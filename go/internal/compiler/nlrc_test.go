package compiler

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// CompilerTestSuite defines the test suite for compiler functions
type CompilerTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest sets up the test environment
func (suite *CompilerTestSuite) SetupTest() {
	tempDir, err := os.MkdirTemp("", "genlinx_compiler_test")
	suite.Require().NoError(err)
	suite.tempDir = tempDir
}

// TearDownTest cleans up the test environment
func (suite *CompilerTestSuite) TearDownTest() {
	os.RemoveAll(suite.tempDir)
}

// TestNewNLRCCompiler tests creating a new NLRC compiler instance
func (suite *CompilerTestSuite) TestNewNLRCCompiler() {
	executablePath := "C:/Program Files/NLRC.exe"
	compiler := NewNLRCCompiler(executablePath)

	assert.NotNil(suite.T(), compiler)
	assert.Equal(suite.T(), executablePath, compiler.ExecutablePath)
}

// TestBuildArgs tests building command line arguments
func (suite *CompilerTestSuite) TestBuildArgs() {
	compiler := NewNLRCCompiler("test.exe")

	options := CompileOptions{
		SourceFiles: []string{"file1.axs", "file2.axi"},
		CFGFiles:    []string{}, // Empty to test source files path
		IncludePath: []string{"C:/Include/Path1", "C:/Include/Path2"},
		ModulePath:  []string{"C:/Module/Path1", "C:/Module/Path2"},
		LibraryPath: []string{"C:/Library/Path1", "C:/Library/Path2"},
		OutputPath:  "C:/Output/Path",
		Verbose:     true,
	}

	args, err := compiler.BuildArgs(options)
	suite.Require().NoError(err)

	// Verify source files come first
	assert.Equal(suite.T(), "file1.axs", args[0])
	assert.Equal(suite.T(), "file2.axi", args[1])

	// Verify include paths
	assert.Contains(suite.T(), args, "-I\"C:/Include/Path1\"")
	assert.Contains(suite.T(), args, "-I\"C:/Include/Path2\"")

	// Verify module paths (should be combined with -M flag)
	moduleArg := ""
	for _, arg := range args {
		if strings.HasPrefix(arg, "-M\"-M") {
			moduleArg = arg
			break
		}
	}
	assert.NotEmpty(suite.T(), moduleArg)
	assert.Contains(suite.T(), moduleArg, "C:/Module/Path1")
	assert.Contains(suite.T(), moduleArg, "C:/Module/Path2")

	// Verify library paths
	assert.Contains(suite.T(), args, "-L\"C:/Library/Path1\"")
	assert.Contains(suite.T(), args, "-L\"C:/Library/Path2\"")

	// Verify output path
	assert.Contains(suite.T(), args, "-O\"C:/Output/Path\"")
}

// TestBuildArgsWithCFG tests building arguments with CFG files
func (suite *CompilerTestSuite) TestBuildArgsWithCFG() {
	compiler := NewNLRCCompiler("test.exe")

	options := CompileOptions{
		SourceFiles: []string{}, // Empty to test CFG path
		CFGFiles:    []string{"config1.cfg", "config2.cfg"},
		IncludePath: []string{"C:/Include"},
		Verbose:     false,
	}

	args, err := compiler.BuildArgs(options)
	suite.Require().NoError(err)

	// Verify CFG files are combined
	assert.Contains(suite.T(), args, "-Cconfig1.cfg;config2.cfg")

	// Verify include paths are still included
	assert.Contains(suite.T(), args, "-I\"C:/Include\"")
}

// TestBuildArgsEmptyPaths tests building arguments with empty paths
func (suite *CompilerTestSuite) TestBuildArgsEmptyPaths() {
	compiler := NewNLRCCompiler("test.exe")

	options := CompileOptions{
		SourceFiles: []string{"test.axs"},
		CFGFiles:    []string{},
		IncludePath: []string{"", "valid/path"}, // Empty string should be ignored
		ModulePath:  []string{},                 // Empty slice
		LibraryPath: []string{"valid/lib"},
		Verbose:     false,
	}

	args, err := compiler.BuildArgs(options)
	suite.Require().NoError(err)

	// Should contain source file
	assert.Contains(suite.T(), args, "test.axs")

	// Should contain valid include path but not empty one
	assert.Contains(suite.T(), args, "-I\"valid/path\"")
	assert.NotContains(suite.T(), args, "-I")

	// Should contain library path
	assert.Contains(suite.T(), args, "-L\"valid/lib\"")

	// Should not contain module arguments (empty slice)
	for _, arg := range args {
		assert.False(suite.T(), strings.HasPrefix(arg, "-M"), "Should not contain module arguments for empty slice")
	}
}

// TestBuildArgsNoOptions tests building arguments with minimal options
func (suite *CompilerTestSuite) TestBuildArgsNoOptions() {
	compiler := NewNLRCCompiler("test.exe")

	options := CompileOptions{
		SourceFiles: []string{"minimal.axs"},
		CFGFiles:    []string{},
		Verbose:     false,
	}

	args, err := compiler.BuildArgs(options)
	suite.Require().NoError(err)

	// Should only contain the source file
	assert.Len(suite.T(), args, 1)
	assert.Equal(suite.T(), "minimal.axs", args[0])
}

// TestCompileOptionsStruct tests the CompileOptions struct
func (suite *CompilerTestSuite) TestCompileOptionsStruct() {
	options := CompileOptions{
		SourceFiles: []string{"file1.axs", "file2.axi"},
		CFGFiles:    []string{"config.cfg"},
		IncludePath: []string{"include/path"},
		ModulePath:  []string{"module/path"},
		LibraryPath: []string{"library/path"},
		OutputPath:  "output/path",
		Verbose:     true,
	}

	assert.Len(suite.T(), options.SourceFiles, 2)
	assert.Len(suite.T(), options.CFGFiles, 1)
	assert.Len(suite.T(), options.IncludePath, 1)
	assert.Len(suite.T(), options.ModulePath, 1)
	assert.Len(suite.T(), options.LibraryPath, 1)
	assert.Equal(suite.T(), "output/path", options.OutputPath)
	assert.True(suite.T(), options.Verbose)
}

// TestCompileResultStruct tests the CompileResult struct
func (suite *CompilerTestSuite) TestCompileResultStruct() {
	result := &CompileResult{
		Success:  true,
		Output:   "Compilation successful",
		Errors:   []string{"No errors"},
		Warnings: []string{"No warnings"},
		ExitCode: 0,
	}

	assert.True(suite.T(), result.Success)
	assert.Equal(suite.T(), "Compilation successful", result.Output)
	assert.Len(suite.T(), result.Errors, 1)
	assert.Len(suite.T(), result.Warnings, 1)
	assert.Equal(suite.T(), 0, result.ExitCode)
}

// TestNLRCCompilerStruct tests the NLRCCompiler struct
func (suite *CompilerTestSuite) TestNLRCCompilerStruct() {
	compiler := &NLRCCompiler{
		ExecutablePath: "C:/NLRC.exe",
	}

	assert.Equal(suite.T(), "C:/NLRC.exe", compiler.ExecutablePath)
}

// TestCompilerInterface tests that NLRCCompiler implements the Compiler interface
func (suite *CompilerTestSuite) TestCompilerInterface() {
	var compiler Compiler = &NLRCCompiler{
		ExecutablePath: "test.exe",
	}

	assert.NotNil(suite.T(), compiler)
	assert.Implements(suite.T(), (*Compiler)(nil), compiler)
}

// TestBuildArgsOrder tests that arguments are built in the correct order
func (suite *CompilerTestSuite) TestBuildArgsOrder() {
	compiler := NewNLRCCompiler("test.exe")

	options := CompileOptions{
		SourceFiles: []string{"source.axs"},
		CFGFiles:    []string{}, // Empty to use source files
		IncludePath: []string{"include"},
		ModulePath:  []string{"module"},
		LibraryPath: []string{"library"},
		OutputPath:  "output",
		Verbose:     false,
	}

	args, err := compiler.BuildArgs(options)
	suite.Require().NoError(err)

	// Source files should be first
	assert.Equal(suite.T(), "source.axs", args[0])

	// Find positions of different argument types
	sourcePos := -1
	includePos := -1
	modulePos := -1
	libraryPos := -1
	outputPos := -1

	for i, arg := range args {
		switch {
		case arg == "source.axs":
			sourcePos = i
		case strings.HasPrefix(arg, "-I"):
			includePos = i
		case strings.HasPrefix(arg, "-M"):
			modulePos = i
		case strings.HasPrefix(arg, "-L"):
			libraryPos = i
		case strings.HasPrefix(arg, "-O"):
			outputPos = i
		}
	}

	// Verify order: source first, then options
	assert.Equal(suite.T(), 0, sourcePos)
	if includePos >= 0 {
		assert.Greater(suite.T(), includePos, sourcePos)
	}
	if modulePos >= 0 {
		assert.Greater(suite.T(), modulePos, sourcePos)
	}
	if libraryPos >= 0 {
		assert.Greater(suite.T(), libraryPos, sourcePos)
	}
	if outputPos >= 0 {
		assert.Greater(suite.T(), outputPos, sourcePos)
	}
}

// TestCompilerTestSuite runs the test suite
func TestCompilerTestSuite(t *testing.T) {
	suite.Run(t, new(CompilerTestSuite))
}
