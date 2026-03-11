package cmd

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// structToMap
// ---------------------------------------------------------------------------

func TestStructToMap_PreservesJSONKeys(t *testing.T) {
	type inner struct {
		Value string `json:"myValue"`
	}
	type outer struct {
		Nested inner `json:"nested"`
	}

	m, err := structToMap(outer{Nested: inner{Value: "hello"}})
	require.NoError(t, err)
	nested, ok := m["nested"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "hello", nested["myValue"])
}

func TestStructToMap_EmptyStruct(t *testing.T) {
	type empty struct{}
	m, err := structToMap(empty{})
	require.NoError(t, err)
	assert.Empty(t, m)
}

// ---------------------------------------------------------------------------
// getByDottedKey
// ---------------------------------------------------------------------------

func TestGetByDottedKey_TopLevel(t *testing.T) {
	m := map[string]any{"key": "value"}
	v, ok := getByDottedKey(m, "key")
	assert.True(t, ok)
	assert.Equal(t, "value", v)
}

func TestGetByDottedKey_Nested(t *testing.T) {
	m := map[string]any{
		"nlrc": map[string]any{
			"path": "/usr/bin/nlrc",
		},
	}
	v, ok := getByDottedKey(m, "nlrc.path")
	assert.True(t, ok)
	assert.Equal(t, "/usr/bin/nlrc", v)
}

func TestGetByDottedKey_MissingKey(t *testing.T) {
	m := map[string]any{"key": "value"}
	_, ok := getByDottedKey(m, "missing")
	assert.False(t, ok)
}

func TestGetByDottedKey_MissingNestedKey(t *testing.T) {
	m := map[string]any{"a": map[string]any{"b": "val"}}
	_, ok := getByDottedKey(m, "a.c")
	assert.False(t, ok)
}

func TestGetByDottedKey_NonMapIntermediate(t *testing.T) {
	m := map[string]any{"a": "string-not-map"}
	_, ok := getByDottedKey(m, "a.b")
	assert.False(t, ok)
}

// ---------------------------------------------------------------------------
// printValue
// ---------------------------------------------------------------------------

func TestPrintValue_String(t *testing.T) {
	out := captureStdout(t, func() { printValue("hello world") })
	assert.Equal(t, "hello world", strings.TrimSpace(out))
}

func TestPrintValue_Bool(t *testing.T) {
	out := captureStdout(t, func() { printValue(true) })
	assert.Equal(t, "true", strings.TrimSpace(out))
}

func TestPrintValue_WholeNumber(t *testing.T) {
	// JSON numbers unmarshal as float64; whole numbers should print as integers.
	out := captureStdout(t, func() { printValue(float64(42)) })
	assert.Equal(t, "42", strings.TrimSpace(out))
}

func TestPrintValue_Float(t *testing.T) {
	out := captureStdout(t, func() { printValue(3.14) })
	assert.Contains(t, strings.TrimSpace(out), "3.14")
}

func TestPrintValue_Slice(t *testing.T) {
	out := captureStdout(t, func() { printValue([]any{"a", "b"}) })
	// Should be JSON array.
	assert.Contains(t, out, `"a"`)
	assert.Contains(t, out, `"b"`)
}

// ---------------------------------------------------------------------------
// printConfig
// ---------------------------------------------------------------------------

func TestPrintConfig_StructProducesJSON(t *testing.T) {
	type cfg struct {
		Name string `json:"name"`
	}

	out := captureStdout(t, func() { printConfig(cfg{Name: "test"}) })
	assert.Contains(t, out, `"name"`)
	assert.Contains(t, out, `"test"`)
}

// ---------------------------------------------------------------------------
// resolveEditor
// ---------------------------------------------------------------------------

func TestResolveEditor_UsesEditorEnvVar(t *testing.T) {
	t.Setenv("VISUAL", "")
	t.Setenv("EDITOR", "vim")
	e, err := resolveEditor()
	require.NoError(t, err)
	assert.Equal(t, "vim", e)
}

func TestResolveEditor_VisualTakesPrecedenceOverEditor(t *testing.T) {
	t.Setenv("VISUAL", "emacs")
	t.Setenv("EDITOR", "vim")
	e, err := resolveEditor()
	require.NoError(t, err)
	assert.Equal(t, "emacs", e)
}

func TestResolveEditor_NoEditorNoFallback_ReturnsError(t *testing.T) {
	t.Setenv("VISUAL", "")
	t.Setenv("EDITOR", "")
	t.Setenv("PATH", t.TempDir()) // empty dir — no binaries available
	_, err := resolveEditor()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "EDITOR")
}

// ---------------------------------------------------------------------------
// ensureConfigFile
// ---------------------------------------------------------------------------

func TestEnsureConfigFile_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "config.json")

	require.NoError(t, ensureConfigFile(path))

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "{}\n", string(data))
}

func TestEnsureConfigFile_DoesNotOverwriteExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"existing":true}`), 0o644))

	require.NoError(t, ensureConfigFile(path))

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, `{"existing":true}`, string(data))
}

// ---------------------------------------------------------------------------
// loadConfigResult
// ---------------------------------------------------------------------------

func TestLoadConfigResult_GlobalFlag(t *testing.T) {
	// Point global config dir at an empty temp dir — should return not-found.
	t.Setenv("GENLINX_CONFIG_DIR", t.TempDir())

	r, label, err := loadConfigResult(true, false)
	require.NoError(t, err)
	assert.Equal(t, "global", label)
	assert.False(t, r.Found)
}

func TestLoadConfigResult_LocalFlag(t *testing.T) {
	// Run from a temp dir with no config file — should return not-found.
	dir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd) //nolint:errcheck
	require.NoError(t, os.Chdir(dir))

	r, label, err := loadConfigResult(false, true)
	require.NoError(t, err)
	assert.Equal(t, "local", label)
	assert.False(t, r.Found)
}

// ---------------------------------------------------------------------------
// mergePrintConfigs
// ---------------------------------------------------------------------------

func TestMergePrintConfigs_ReturnsBuildKey(t *testing.T) {
	// Point global config to empty dir so test is isolated.
	t.Setenv("GENLINX_CONFIG_DIR", t.TempDir())

	dir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd) //nolint:errcheck
	require.NoError(t, os.Chdir(dir))

	m, err := mergePrintConfigs()
	require.NoError(t, err)
	_, hasBuild := m["build"]
	assert.True(t, hasBuild, "merged config must have a 'build' key")
}

// ---------------------------------------------------------------------------
// configList
// ---------------------------------------------------------------------------

func TestConfigList_CombinedPrintsConfig(t *testing.T) {
	t.Setenv("GENLINX_CONFIG_DIR", t.TempDir())
	dir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd) //nolint:errcheck
	require.NoError(t, os.Chdir(dir))

	out := captureStdout(t, func() {
		err := configList(false, false)
		require.NoError(t, err)
	})

	// Combined list always prints the merged config as JSON.
	assert.Contains(t, out, `"build"`)
}

func TestConfigList_GlobalNotFound_PrintsMessage(t *testing.T) {
	t.Setenv("GENLINX_CONFIG_DIR", t.TempDir())

	out := captureStdout(t, func() {
		err := configList(true, false)
		require.NoError(t, err)
	})

	assert.Contains(t, out, "No global configuration found.")
}

func TestConfigList_LocalNotFound_PrintsMessage(t *testing.T) {
	// Run from a temp dir with no config file.
	dir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd) //nolint:errcheck
	require.NoError(t, os.Chdir(dir))

	out := captureStdout(t, func() {
		err := configList(false, true)
		require.NoError(t, err)
	})

	assert.Contains(t, out, "No local configuration found.")
}

func TestConfigList_LocalFound_PrintsFileContents(t *testing.T) {
	dir := t.TempDir()
	content := `{"build":{"nlrc":{"path":"custom.exe"}}}`
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".genlinxrc.json"), []byte(content), 0o644))

	origWd, _ := os.Getwd()
	defer os.Chdir(origWd) //nolint:errcheck
	require.NoError(t, os.Chdir(dir))

	out := captureStdout(t, func() {
		err := configList(false, true)
		require.NoError(t, err)
	})

	assert.Contains(t, out, "custom.exe")
}

// ---------------------------------------------------------------------------
// configGet
// ---------------------------------------------------------------------------

func TestConfigGet_ExistingKey_PrintsValue(t *testing.T) {
	t.Setenv("GENLINX_CONFIG_DIR", t.TempDir())
	dir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd) //nolint:errcheck
	require.NoError(t, os.Chdir(dir))

	out := captureStdout(t, func() {
		require.NoError(t, configGet("nlrc.path", false, false))
	})

	// Default NLRC path is non-empty; output should be the path value.
	assert.NotEmpty(t, strings.TrimSpace(out))
}

func TestConfigGet_MissingKey_PrintsNoConfigFound(t *testing.T) {
	t.Setenv("GENLINX_CONFIG_DIR", t.TempDir())
	dir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd) //nolint:errcheck
	require.NoError(t, os.Chdir(dir))

	out := captureStdout(t, func() {
		err := configGet("nonexistent.key.path", false, false)
		require.NoError(t, err)
	})

	assert.Contains(t, out, "No configuration found for key")
}

func TestConfigGet_GlobalNotFound(t *testing.T) {
	t.Setenv("GENLINX_CONFIG_DIR", t.TempDir())

	out := captureStdout(t, func() {
		err := configGet("nlrc.path", true, false)
		require.NoError(t, err)
	})

	assert.Contains(t, out, "No global configuration found.")
}

// captureStdout redirects os.Stdout for the duration of fn and returns the
// captured output.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	require.NoError(t, err)

	orig := os.Stdout
	os.Stdout = w

	fn()

	require.NoError(t, w.Close())
	os.Stdout = orig

	out, err := io.ReadAll(r)
	require.NoError(t, err)

	return string(out)
}

// requireNLRC is a helper that walks the parsed JSON map down to the
// nlrc object, failing the test if it is missing.
func requireNLRC(t *testing.T, out string) map[string]any {
	t.Helper()

	var m map[string]any
	require.NoError(t, json.Unmarshal([]byte(out), &m), "output must be valid JSON")

	nlrc, ok := m["nlrc"].(map[string]any)
	require.True(t, ok, "nlrc key must be present, got: %s", out)

	return nlrc
}

// TestPrintRawFileConfig_JSON verifies JSON pretty-printing with camelCase key preservation.
func TestPrintRawFileConfig_JSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".genlinxrc.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"nlrc":{"includePath":["./include"]}}`), 0o644))

	out := captureStdout(t, func() { printRawFileConfig(path) })

	nlrc := requireNLRC(t, out)
	paths, ok := nlrc["includePath"].([]any)
	require.True(t, ok, "includePath key must be present (camelCase preserved)")
	assert.Equal(t, "./include", paths[0])
}

// TestPrintRawFileConfig_YAML verifies YAML → JSON with original key casing.
// Regression test for config --list --local showing "{}" for YAML files.
func TestPrintRawFileConfig_YAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".genlinxrc.yaml")
	content := "nlrc:\n  includePath:\n    - ./include\n  modulePath:\n    - ./module\n"
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))

	out := captureStdout(t, func() { printRawFileConfig(path) })

	assert.NotEqual(t, "{}", strings.TrimSpace(out), "YAML file must not produce empty output")

	nlrc := requireNLRC(t, out)
	paths, ok := nlrc["includePath"].([]any)
	require.True(t, ok, "includePath key must be present (camelCase preserved)")
	assert.Equal(t, "./include", paths[0])

	modPaths, ok := nlrc["modulePath"].([]any)
	require.True(t, ok, "modulePath key must be present")
	assert.Equal(t, "./module", modPaths[0])
}

// TestPrintRawFileConfig_YML verifies .yml extension is handled identically to .yaml.
func TestPrintRawFileConfig_YML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".genlinxrc.yml")
	require.NoError(t, os.WriteFile(path, []byte("nlrc:\n  includePath:\n    - ./include\n"), 0o644))

	out := captureStdout(t, func() { printRawFileConfig(path) })

	assert.NotEqual(t, "{}", strings.TrimSpace(out))
	nlrc := requireNLRC(t, out)
	paths, ok := nlrc["includePath"].([]any)
	require.True(t, ok, "includePath key must be present (camelCase preserved)")
	assert.Equal(t, "./include", paths[0])
}

// TestPrintRawFileConfig_MissingFile verifies that a missing file prints "{}".
func TestPrintRawFileConfig_MissingFile(t *testing.T) {
	out := captureStdout(t, func() {
		printRawFileConfig(filepath.Join(t.TempDir(), "nonexistent.json"))
	})

	assert.Equal(t, "{}", strings.TrimSpace(out))
}

// TestPrintRawFileConfig_EmptyFile verifies that an empty file prints "{}".
func TestPrintRawFileConfig_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".genlinxrc.json")
	require.NoError(t, os.WriteFile(path, []byte(""), 0o644))

	out := captureStdout(t, func() { printRawFileConfig(path) })
	assert.Equal(t, "{}", strings.TrimSpace(out))
}
