package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/neilotoole/jsoncolor"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/Norgate-AV/genlinx/internal/options"
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
		// Combined config: default < global < local, fully normalized and deduplicated.
		// Use the same merge path as the actual build/archive/cfg commands so the
		// display always matches what those commands will use.
		mergedCfg, _, err := options.LoadMergedConfig()
		if err != nil {
			return err
		}

		printConfig(mergedCfg)
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

	printRawFileConfig(r.Path)
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

	editor, err := resolveEditor()
	if err != nil {
		return err
	}

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
		merged, err := mergePrintConfigs()
		if err != nil {
			return err
		}

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

// resolveEditor returns the editor to use.
// Resolution order: $VISUAL → $EDITOR (honoured as-is, user's explicit choice)
// → core.editor from the genlinx config → first available fallback found on PATH.
// Returns an error when no fallback is found, directing the user to set $EDITOR.
func resolveEditor() (string, error) {
	for _, env := range []string{"VISUAL", "EDITOR"} {
		if e := os.Getenv(env); e != "" {
			return e, nil
		}
	}

	// Honour core.editor from the merged genlinx config file.
	if mergedCfg, _, err := options.LoadMergedConfig(); err == nil {
		if e := mergedCfg.Core.Editor; e != "" {
			return e, nil
		}
	}

	for _, candidate := range []string{"nvim", "vim", "nano", "edit", "code", "notepad"} {
		if _, err := exec.LookPath(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("no editor found on PATH; set the EDITOR environment variable to your preferred editor")
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
func mergePrintConfigs() (map[string]any, error) {
	mergedCfg, _, err := options.LoadMergedConfig()
	if err != nil {
		return nil, err
	}

	return structToMap(mergedCfg)
}

func printConfig(v any) {
	enc := jsoncolor.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	if jsoncolor.IsColorTerminal(os.Stdout) {
		enc.SetColors(jsoncolor.DefaultColors())
	}

	if err := enc.Encode(v); err != nil {
		fmt.Printf("%+v\n", v)
	}
}

// printRawFileConfig reads the config file at path and pretty-prints its
// contents as JSON. Supports both JSON and YAML file formats, preserving the
// original key casing from the file.
// If the file is absent or empty it prints "{}".
func printRawFileConfig(path string) {
	data, err := os.ReadFile(path)
	if err != nil || len(bytes.TrimSpace(data)) == 0 {
		fmt.Println("{}")
		return
	}

	ext := strings.ToLower(filepath.Ext(path))

	var m any
	if ext == ".yaml" || ext == ".yml" {
		if err := yaml.Unmarshal(data, &m); err != nil {
			fmt.Println("{}")
			return
		}
	} else {
		if err := json.Unmarshal(data, &m); err != nil {
			fmt.Println("{}")
			return
		}
	}

	enc := jsoncolor.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	if jsoncolor.IsColorTerminal(os.Stdout) {
		enc.SetColors(jsoncolor.DefaultColors())
	}

	if err := enc.Encode(m); err != nil {
		fmt.Println("{}")
	}
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
