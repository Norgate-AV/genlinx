package cmd

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Norgate-AV/genlinx/internal/options"
)

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
