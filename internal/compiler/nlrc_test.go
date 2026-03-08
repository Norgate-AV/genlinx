package compiler

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// TestMain acts as a fake NLRC compiler when GENLINX_FAKE_COMPILER is set.
// This allows TestCompile_AlwaysStreamsOutput to use the test binary itself
// as a controllable subprocess without requiring a real NLRC installation.
func TestMain(m *testing.M) {
	if msg := os.Getenv("GENLINX_FAKE_COMPILER"); msg != "" {
		fmt.Println(msg)
		os.Exit(0)
	}
	os.Exit(m.Run())
}

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

	// Verify include paths (joined into single arg)
	assert.Contains(suite.T(), args, "-I\"C:/Include/Path1;C:/Include/Path2\"")

	// Verify module paths (should be combined into single -M"path1;path2" flag)
	assert.Contains(suite.T(), args, "-M\"C:/Module/Path1;C:/Module/Path2\"")
	for _, arg := range args {
		if strings.HasPrefix(arg, "-M") {
			assert.False(suite.T(), strings.HasPrefix(arg, "-M\"-M"), "module arg must not have double -M prefix: %s", arg)
		}
	}

	// Verify library paths (joined into single arg)
	assert.Contains(suite.T(), args, "-L\"C:/Library/Path1;C:/Library/Path2\"")

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

	// Verify CFG files are combined with correct flag
	assert.Contains(suite.T(), args, "-CFGconfig1.cfg;config2.cfg")

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

// TestParseOutput_Empty verifies that empty output produces no errors or warnings.
func (suite *CompilerTestSuite) TestParseOutput_Empty() {
	compiler := NewNLRCCompiler("test.exe")
	errors, warnings := compiler.parseOutput("")
	assert.Empty(suite.T(), errors)
	assert.Empty(suite.T(), warnings)
}

// TestParseOutput_Errors verifies that lines containing the NLRC "ERROR: "
// prefix are classified as errors.
func (suite *CompilerTestSuite) TestParseOutput_Errors() {
	compiler := NewNLRCCompiler("test.exe")
	output := "ERROR: undefined variable 'foo'\nERROR: missing semicolon"
	errors, warnings := compiler.parseOutput(output)
	assert.Len(suite.T(), errors, 2)
	assert.Empty(suite.T(), warnings)
	assert.Contains(suite.T(), errors, "ERROR: undefined variable 'foo'")
	assert.Contains(suite.T(), errors, "ERROR: missing semicolon")
}

// TestParseOutput_Warnings verifies that lines containing the NLRC "WARNING: "
// prefix are classified as warnings.
func (suite *CompilerTestSuite) TestParseOutput_Warnings() {
	compiler := NewNLRCCompiler("test.exe")
	output := "WARNING: unused variable 'x'\nWARNING: deprecated function"
	errors, warnings := compiler.parseOutput(output)
	assert.Empty(suite.T(), errors)
	assert.Len(suite.T(), warnings, 2)
	assert.Contains(suite.T(), warnings, "WARNING: unused variable 'x'")
	assert.Contains(suite.T(), warnings, "WARNING: deprecated function")
}

// TestParseOutput_Mixed verifies that a realistic NLRC output block is
// classified correctly — plain lines are ignored.
func (suite *CompilerTestSuite) TestParseOutput_Mixed() {
	compiler := NewNLRCCompiler("test.exe")
	output := "Compiling test.axs\nERROR: undefined variable\nWARNING: deprecated usage\nCompilation complete"
	errors, warnings := compiler.parseOutput(output)
	assert.Len(suite.T(), errors, 1)
	assert.Len(suite.T(), warnings, 1)
	assert.Contains(suite.T(), errors, "ERROR: undefined variable")
	assert.Contains(suite.T(), warnings, "WARNING: deprecated usage")
}

// TestParseOutput_NonNLRCFormatIgnored verifies that lines using mixed-case or
// non-NLRC prefixes (e.g. "Error:", "warning:", "0 errors found") are not
// matched, preventing false positives.
func (suite *CompilerTestSuite) TestParseOutput_NonNLRCFormatIgnored() {
	compiler := NewNLRCCompiler("test.exe")
	output := "Error: wrong case\nwarning: also wrong\n0 error(s) found\nBuild succeeded with 0 warnings"
	errors, warnings := compiler.parseOutput(output)
	assert.Empty(suite.T(), errors)
	assert.Empty(suite.T(), warnings)
}

// TestParseOutput_DuplicatesDeduped verifies that duplicate log lines are only
// reported once, matching the TypeScript Set deduplication behaviour.
func (suite *CompilerTestSuite) TestParseOutput_DuplicatesDeduped() {
	compiler := NewNLRCCompiler("test.exe")
	output := "ERROR: undefined variable\nERROR: undefined variable\nWARNING: deprecated usage\nWARNING: deprecated usage"
	errors, warnings := compiler.parseOutput(output)
	assert.Len(suite.T(), errors, 1)
	assert.Len(suite.T(), warnings, 1)
}

// TestParseOutput_ErrorPrecedesWarningKeyword verifies that a line with NLRC
// "ERROR: " prefix is classified as an error even if it also contains the word
// "warning" in its message.
func (suite *CompilerTestSuite) TestParseOutput_ErrorPrecedesWarningKeyword() {
	compiler := NewNLRCCompiler("test.exe")
	output := "ERROR: this warning-like message is still an error"
	errors, warnings := compiler.parseOutput(output)
	assert.Len(suite.T(), errors, 1)
	assert.Empty(suite.T(), warnings)
}

// TestBuildArgs_CFGTakesPrecedenceOverSourceFiles verifies that when both
// CFGFiles and SourceFiles are provided, CFG mode is used exclusively and the
// source files are not added to the argument list.
func (suite *CompilerTestSuite) TestBuildArgs_CFGTakesPrecedenceOverSourceFiles() {
	compiler := NewNLRCCompiler("test.exe")
	options := CompileOptions{
		SourceFiles: []string{"should_be_ignored.axs"},
		CFGFiles:    []string{"project.cfg"},
	}

	args, err := compiler.BuildArgs(options)
	suite.Require().NoError(err)

	// The -CFG flag must be present
	assert.Contains(suite.T(), args, "-CFGproject.cfg")

	// Source files must NOT appear
	for _, arg := range args {
		assert.NotEqual(suite.T(), "should_be_ignored.axs", arg,
			"source file should be ignored when CFG files are provided")
	}
}

// TestBuildArgs_BothSourceAndCFGEmpty verifies that when both SourceFiles and
// CFGFiles are empty, no source or cfg arguments are generated.
func (suite *CompilerTestSuite) TestBuildArgs_BothSourceAndCFGEmpty() {
	compiler := NewNLRCCompiler("test.exe")
	options := CompileOptions{
		SourceFiles: []string{},
		CFGFiles:    []string{},
	}

	args, err := compiler.BuildArgs(options)
	suite.Require().NoError(err)

	for _, arg := range args {
		assert.False(suite.T(), strings.HasPrefix(arg, "-CFG"),
			"should not contain a -CFG arg when no CFG files are provided")
	}
	assert.Empty(suite.T(), args)
}

// TestBuildArgs_SourceFileExistsOnDisk verifies that a source file that exists
// on the filesystem is resolved to its absolute path in the argument list.
func (suite *CompilerTestSuite) TestBuildArgs_SourceFileExistsOnDisk() {
	// Create a real source file in the temp directory
	sourceFile := filepath.Join(suite.tempDir, "TestMain.axs")
	suite.Require().NoError(os.WriteFile(sourceFile, []byte("PROGRAM_NAME='TestMain'\n"), 0o644))

	compiler := NewNLRCCompiler("test.exe")
	options := CompileOptions{
		SourceFiles: []string{sourceFile},
		CFGFiles:    []string{},
	}

	args, err := compiler.BuildArgs(options)
	suite.Require().NoError(err)
	suite.Require().NotEmpty(args)

	assert.True(suite.T(), filepath.IsAbs(args[0]),
		"source file arg should be an absolute path, got: %s", args[0])
	assert.Equal(suite.T(), sourceFile, args[0])
}

// TestCompile_BadExecutable verifies that Compile returns an error whose message
// contains "failed to start compiler" when the executable does not exist.
func (suite *CompilerTestSuite) TestCompile_BadExecutable() {
	compiler := NewNLRCCompiler("/no/such/compiler/nlrc.exe")
	options := CompileOptions{
		SourceFiles: []string{"test.axs"},
	}

	result, err := compiler.Compile(options)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "failed to start compiler")
}

// captureStdout runs fn and returns whatever it wrote to os.Stdout.
func captureStdout(fn func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r) //nolint:errcheck
	return buf.String()
}

// TestReadOutput_StreamsAllOutput verifies that readOutput prints every line
// immediately (streams in real time) and returns the full accumulated output.
func (suite *CompilerTestSuite) TestReadOutput_StreamsAllOutput() {
	c := NewNLRCCompiler("test.exe")

	captured := captureStdout(func() {
		out, err := c.readOutput(
			strings.NewReader("line1\nline2\n"),
			strings.NewReader("err1\n"),
		)
		suite.Require().NoError(err)
		suite.Contains(out, "line1")
		suite.Contains(out, "line2")
		suite.Contains(out, "err1")
	})

	assert.Contains(suite.T(), captured, "line1")
	assert.Contains(suite.T(), captured, "line2")
	assert.Contains(suite.T(), captured, "err1")
}

// TestCompile_AlwaysStreamsOutput is a regression test for the bug where
// readOutput was changed to accept a `verbose bool` parameter and only printed
// when verbose=true. Because the verbose flag was never registered in the build
// command's init(), it always defaulted to false, silencing all output.
//
// This test exercises the full Compile path using the test binary itself as a
// fake NLRC compiler (via the GENLINX_FAKE_COMPILER env var handled in
// TestMain) and asserts that every output line is streamed to stdout even when
// CompileOptions.Verbose is false.
func (suite *CompilerTestSuite) TestCompile_AlwaysStreamsOutput() {
	const fakeOutput = "Compiling placeholder.axs..."

	captured := captureStdout(func() {
		c := &NLRCCompiler{
			ExecutablePath: os.Args[0],
			env:            append(os.Environ(), "GENLINX_FAKE_COMPILER="+fakeOutput),
		}
		_, _ = c.Compile(CompileOptions{
			SourceFiles: []string{"placeholder.axs"},
			Verbose:     false, // output must still stream regardless
		})
	})

	assert.Contains(suite.T(), captured, fakeOutput,
		"Compile must always stream output to stdout regardless of the Verbose option")
}

// TestCompilerTestSuite runs the test suite
func TestCompilerTestSuite(t *testing.T) {
	suite.Run(t, new(CompilerTestSuite))
}
