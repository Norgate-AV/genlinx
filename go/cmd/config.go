package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Norgate-AV/genlinx-go/internal/options"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config [key]",
	Short: "view/edit configuration properties for genlinx",
	Long: `View/edit configuration properties for genlinx.

If an option is omitted, genlinx will:
  1. look for a .genlinxrc.json file in the root directory of the project;
     if one is found, any options defined will be set to the value in the file
  2. use the default values from the global config file for any remaining options`,
	Args: cobra.MaximumNArgs(1),
	RunE: runConfig,
}

func runConfig(cmd *cobra.Command, args []string) error {
	list, _ := cmd.Flags().GetBool("list")
	edit, _ := cmd.Flags().GetBool("edit")
	global, _ := cmd.Flags().GetBool("global")
	local, _ := cmd.Flags().GetBool("local")

	if list {
		return configList(global, local)
	}

	if edit {
		return configEdit(global, local)
	}

	if len(args) > 0 {
		return configGet(args[0], global, local)
	}

	return cmd.Help()
}

// ---------------------------------------------------------------------------
// list
// ---------------------------------------------------------------------------

func configList(global, local bool) error {
	if !global && !local {
		// Combined config (global + local merged).
		globalResult, _ := options.LoadGlobalConfig()
		localResult, _ := options.LoadLocalConfig()
		printMergedConfig(globalResult, localResult)
		return nil
	}

	r, label, err := loadConfigResult(global, local)
	if err != nil {
		return err
	}

	if !r.Found {
		fmt.Printf("No %s configuration found.\n", label)
		return nil
	}

	printConfig(r.Config)
	return nil
}

// ---------------------------------------------------------------------------
// edit
// ---------------------------------------------------------------------------

func configEdit(global, local bool) error {
	if !global && !local {
		fmt.Fprintln(os.Stderr, "Please specify --global or --local")
		return nil
	}

	r, label, err := loadConfigResult(global, local)
	if err != nil {
		return err
	}

	var filePath string

	if !r.Found {
		if global {
			// Create the file so it can be edited.
			filePath, err = options.GlobalConfigPath()
			if err != nil {
				return err
			}

			if err := ensureConfigFile(filePath); err != nil {
				return err
			}
		} else {
			fmt.Printf("No %s configuration found.\n", label)
			return nil
		}
	} else {
		filePath = r.Path
	}

	editor := resolveEditor()

	fmt.Printf("Opening %s in %s...\n", filePath, editor)

	c := exec.Command(editor, filePath) //nolint:gosec // editor comes from a trusted env var / constant
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return c.Run()
}

// ---------------------------------------------------------------------------
// get
// ---------------------------------------------------------------------------

func configGet(key string, global, local bool) error {
	if !global && !local {
		// Query the merged config.
		globalResult, _ := options.LoadGlobalConfig()
		localResult, _ := options.LoadLocalConfig()

		merged := mergePrintConfigs(globalResult, localResult)

		value, found := getByDottedKey(merged, key)
		if !found {
			fmt.Printf("No configuration found for key: %s\n", key)
			return nil
		}

		printValue(value)
		return nil
	}

	r, label, err := loadConfigResult(global, local)
	if err != nil {
		return err
	}

	if !r.Found {
		fmt.Printf("No %s configuration found.\n", label)
		return nil
	}

	m, err := structToMap(r.Config)
	if err != nil {
		return err
	}

	value, found := getByDottedKey(m, key)
	if !found {
		fmt.Printf("No configuration found for key: %s\n", key)
		return nil
	}

	printValue(value)
	return nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func loadConfigResult(global, local bool) (options.ConfigLoadResult, string, error) {
	if global {
		r, err := options.LoadGlobalConfig()
		return r, "global", err
	}

	r, err := options.LoadLocalConfig()
	return r, "local", err
}

// resolveEditor returns the editor to use, mirroring the TS getEditor logic:
// $EDITOR env var → "code" (VS Code) as fallback.
func resolveEditor() string {
	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}

	return "code"
}

// ensureConfigFile creates an empty JSON config file at path (and its parent
// directories) if it doesn't already exist.
func ensureConfigFile(path string) error {
	// Create parent directory if needed.
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, []byte("{}\n"), 0o644); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
	}

	return nil
}

// structToMap converts a struct to a map[string]any using JSON round-trip so
// that dotted-key lookups work against the JSON field names.
func structToMap(v any) (map[string]any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return m, nil
}

// getByDottedKey resolves a dot-separated key (e.g. "archive.outputFileSuffix")
// against a nested map, mirroring lodash's _.get behaviour.
func getByDottedKey(m map[string]any, key string) (any, bool) {
	parts := strings.Split(key, ".")
	var current any = m

	for _, part := range parts {
		v, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}

		val, ok := v[part]
		if !ok {
			return nil, false
		}

		current = val
	}

	return current, true
}

// mergePrintConfigs returns a combined map for key lookup (default + global + local).
func mergePrintConfigs(global, local options.ConfigLoadResult) map[string]any {
	// Load default as base, then overlay global and local.
	m := map[string]any{}

	for _, r := range []options.ConfigLoadResult{global, local} {
		if !r.Found || r.Config == nil {
			continue
		}

		overlay, err := structToMap(r.Config)
		if err != nil {
			continue
		}

		deepMerge(m, overlay)
	}

	return m
}

// deepMerge merges src into dst, recursing into nested maps.
func deepMerge(dst, src map[string]any) {
	for k, v := range src {
		if srcMap, ok := v.(map[string]any); ok {
			if dstMap, ok := dst[k].(map[string]any); ok {
				deepMerge(dstMap, srcMap)
				continue
			}
		}

		dst[k] = v
	}
}

func printConfig(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("%+v\n", v)
		return
	}

	fmt.Println(string(b))
}

// printMergedConfig pretty-prints the merged (global + local) configuration.
func printMergedConfig(global, local options.ConfigLoadResult) {
	m := mergePrintConfigs(global, local)
	printConfig(m)
}

func printValue(v any) {
	switch val := v.(type) {
	case string:
		fmt.Println(val)
	case bool:
		fmt.Println(val)
	case float64:
		// JSON numbers unmarshal as float64.
		if val == float64(int64(val)) {
			fmt.Println(int64(val))
		} else {
			fmt.Println(val)
		}
	default:
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			fmt.Printf("%+v\n", v)
			return
		}

		fmt.Println(string(b))
	}
}

func init() {
	configCmd.Flags().Bool("global", false, "use global configuration")
	configCmd.Flags().Bool("local", false, "use local configuration")
	configCmd.Flags().BoolP("list", "l", false, "display the configuration in stdout")
	configCmd.Flags().BoolP("edit", "e", false, "edit the configuration with default text editor")

	configCmd.MarkFlagsMutuallyExclusive("global", "local")
	configCmd.MarkFlagsMutuallyExclusive("list", "edit")
}
