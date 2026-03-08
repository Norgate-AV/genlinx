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

// captureStdout redirects os.Stdout for the duration of fn and returns the
// captured output.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	require.NoError(t, err)

	orig := os.Stdout
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = orig

	out, err := io.ReadAll(r)
	require.NoError(t, err)

	return string(out)
}

// requireBuildNLRC is a helper that walks the parsed JSON map down to the
// build.nlrc object, failing the test if any step is missing.
func requireBuildNLRC(t *testing.T, out string) map[string]any {
	t.Helper()

	var m map[string]any
	require.NoError(t, json.Unmarshal([]byte(out), &m), "output must be valid JSON")

	build, ok := m["build"].(map[string]any)
	require.True(t, ok, "build key must be present, got: %s", out)

	nlrc, ok := build["nlrc"].(map[string]any)
	require.True(t, ok, "nlrc key must be present")

	return nlrc
}

// TestPrintRawFileConfig_JSON verifies JSON pretty-printing with camelCase key preservation.
func TestPrintRawFileConfig_JSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".genlinxrc.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"build":{"nlrc":{"includePath":["./include"]}}}`), 0o644))

	out := captureStdout(t, func() { printRawFileConfig(path) })

	nlrc := requireBuildNLRC(t, out)
	paths, ok := nlrc["includePath"].([]any)
	require.True(t, ok, "includePath key must be present (camelCase preserved)")
	assert.Equal(t, "./include", paths[0])
}

// TestPrintRawFileConfig_YAML verifies YAML → JSON with original key casing.
// Regression test for config --list --local showing "{}" for YAML files.
func TestPrintRawFileConfig_YAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".genlinxrc.yaml")
	content := "build:\n  nlrc:\n    includePath:\n      - ./include\n    modulePath:\n      - ./module\n"
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))

	out := captureStdout(t, func() { printRawFileConfig(path) })

	assert.NotEqual(t, "{}", strings.TrimSpace(out), "YAML file must not produce empty output")

	nlrc := requireBuildNLRC(t, out)
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
	require.NoError(t, os.WriteFile(path, []byte("build:\n  nlrc:\n    includePath:\n      - ./include\n"), 0o644))

	out := captureStdout(t, func() { printRawFileConfig(path) })

	assert.NotEqual(t, "{}", strings.TrimSpace(out))
	nlrc := requireBuildNLRC(t, out)
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
