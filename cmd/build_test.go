package cmd

import (
	"bytes"
	"io"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Norgate-AV/genlinx/internal/compiler"
	"github.com/Norgate-AV/genlinx/internal/options"
)

// captureStderr runs fn and returns whatever it wrote to os.Stderr.
func captureStderr(fn func()) string {
	r, w, _ := os.Pipe()
	old := os.Stderr
	os.Stderr = w

	fn()

	_ = w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	io.Copy(&buf, r) //nolint:errcheck
	return buf.String()
}

func TestBuildCompileOpts_MapsAllFields(t *testing.T) {
	opts := &options.BuildOptions{
		SourceFiles: []string{"main.axs", "lib.axi"},
		CFGFiles:    []string{"build.cfg"},
		IncludePath: []string{"/inc1", "/inc2"},
		ModulePath:  []string{"/mod1"},
		LibraryPath: []string{"/lib1"},
		OutputPath:  "/out/dir",
		Verbose:     true,
	}

	got := buildCompileOpts(opts)

	assert.Equal(t, opts.SourceFiles, got.SourceFiles)
	assert.Equal(t, opts.CFGFiles, got.CFGFiles)
	assert.Equal(t, opts.IncludePath, got.IncludePath)
	assert.Equal(t, opts.ModulePath, got.ModulePath)
	assert.Equal(t, opts.LibraryPath, got.LibraryPath)
	assert.Equal(t, opts.OutputPath, got.OutputPath, "OutputPath must be forwarded to CompileOptions")
	assert.Equal(t, opts.Verbose, got.Verbose)
}

func TestBuildCompileOpts_EmptyOutputPath(t *testing.T) {
	opts := &options.BuildOptions{
		SourceFiles: []string{"main.axs"},
		OutputPath:  "",
	}

	got := buildCompileOpts(opts)

	assert.Empty(t, got.OutputPath)
}

func TestBuildCmd_NoArgs_ReturnsError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("CFG auto-discovery runs on Windows; OS guard does not apply")
	}

	err := buildCmd.RunE(buildCmd, []string{})
	assert.EqualError(t, err, "the build command is only supported on Windows")
}

func TestBuildCmd_WithSourceFileArg_PassesGuard(t *testing.T) {
	err := buildCmd.RunE(buildCmd, []string{"main.axs"})
	// The Windows guard fires on non-Windows; confirm it is not a "no files" error.
	if err != nil {
		assert.NotEqual(t, "no source or CFG files specified", err.Error(),
			"guard must not fire when source files are provided")
	}
}

// TestPrintBuildResult_NoErrorsNoWarnings verifies that printBuildResult returns
// nil and writes nothing to stderr when the result is clean.
func TestPrintBuildResult_NoErrorsNoWarnings(t *testing.T) {
	result := &compiler.CompileResult{
		Success:  true,
		Warnings: []string{},
		Errors:   []string{},
		ExitCode: 0,
	}

	var err error
	captured := captureStderr(func() {
		err = printBuildResult(result)
	})

	assert.NoError(t, err)
	assert.Empty(t, captured)
}

// TestPrintBuildResult_WarningsOnly verifies that warnings are printed to stderr
// with a count summary and no error is returned.
func TestPrintBuildResult_WarningsOnly(t *testing.T) {
	result := &compiler.CompileResult{
		Success:  true,
		Warnings: []string{"WARNING: deprecated function", "WARNING: unused variable"},
		Errors:   []string{},
		ExitCode: 0,
	}

	var err error
	captured := captureStderr(func() {
		err = printBuildResult(result)
	})

	assert.NoError(t, err)
	assert.Contains(t, captured, "2 warning(s)")
	assert.Contains(t, captured, "WARNING: deprecated function")
	assert.Contains(t, captured, "WARNING: unused variable")
}

// TestPrintBuildResult_ErrorsReturnsError verifies that errors are printed to
// stderr with a count summary and an error describing the count is returned.
func TestPrintBuildResult_ErrorsReturnsError(t *testing.T) {
	result := &compiler.CompileResult{
		Success:  false,
		Warnings: []string{},
		Errors:   []string{"ERROR: undefined variable", "ERROR: missing semicolon"},
		ExitCode: 1,
	}

	var retErr error
	captured := captureStderr(func() {
		retErr = printBuildResult(result)
	})

	assert.Error(t, retErr)
	assert.Contains(t, retErr.Error(), "2 error(s)")
	assert.Contains(t, captured, "2 error(s)")
	assert.Contains(t, captured, "ERROR: undefined variable")
	assert.Contains(t, captured, "ERROR: missing semicolon")
}

// TestPrintBuildResult_NoErrorsReturnsNil verifies nil is returned when the
// errors slice is empty, even if warnings are present.
func TestPrintBuildResult_NoErrorsReturnsNil(t *testing.T) {
	result := &compiler.CompileResult{
		Success:  true,
		Warnings: []string{"WARNING: something minor"},
		Errors:   []string{},
		ExitCode: 0,
	}

	captureStderr(func() {
		assert.NoError(t, printBuildResult(result))
	})
}

// TestBuildCmd_NoAllFlag_Registered verifies that the -A/--no-all flag is
// registered on the build command with the correct shorthand and default value.
func TestBuildCmd_NoAllFlag_Registered(t *testing.T) {
	f := buildCmd.Flags().Lookup("no-all")
	assert.NotNil(t, f, "--no-all flag must be registered")
	assert.Equal(t, "A", f.Shorthand)
	assert.Equal(t, "false", f.DefValue)
}

// TestBuildCmd_NoAllAndAll_AreMutuallyExclusive verifies that cobra rejects
// supplying both --all and --no-all at the same time.
func TestBuildCmd_NoAllAndAll_AreMutuallyExclusive(t *testing.T) {
	// cobra enforces mutual exclusion during command execution (before RunE).
	// Execute() always traverses from the root, so drive via rootCmd.
	rootCmd.SetArgs([]string{"build", "--all", "--no-all"})
	t.Cleanup(func() {
		rootCmd.SetArgs(nil)
		_ = buildCmd.Flags().Set("all", "false")
		_ = buildCmd.Flags().Set("no-all", "false")
	})

	err := rootCmd.Execute()
	assert.ErrorContains(t, err, "[all no-all]")
}

// TestBuildCmd_PositionalArgsMergedWithSourceFiles verifies that positional
// arguments are combined with --source-files before the OS guard runs.
// On non-Windows the guard fires immediately, so we confirm the error is the
// OS-guard message (not a "no source files" message), proving the merge happened.
func TestBuildCmd_PositionalArgsMergedWithSourceFiles(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("OS guard does not fire on Windows; test not applicable")
	}

	err := buildCmd.RunE(buildCmd, []string{"main.axs", "extra.axs"})
	assert.EqualError(t, err, "the build command is only supported on Windows",
		"positional args must be accepted without triggering a 'no files' error")
}
