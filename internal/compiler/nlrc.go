package compiler

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	cfgFlag       = "-CFG"
	includeFlag   = "-I"
	moduleFlag    = "-M"
	libraryFlag   = "-L"
	outputFlag    = "-O"
	pathDelimiter = ";"
)

// logPattern matches the NLRC compiler log prefix: "ERROR: ..." or "WARNING: ..."
// This mirrors the TypeScript regex (?<level>ERROR|WARNING): .+
var logPattern = regexp.MustCompile(`(ERROR|WARNING): .+`)

// Compiler represents a NetLinx compiler
type Compiler interface {
	Compile(options CompileOptions) (*CompileResult, error)
}

// CompileOptions represents compilation options
type CompileOptions struct {
	SourceFiles []string
	CFGFiles    []string
	IncludePath []string
	ModulePath  []string
	LibraryPath []string
	OutputPath  string
	Verbose     bool
}

// CompileResult represents the result of a compilation
type CompileResult struct {
	Success  bool
	Output   string
	Errors   []string
	Warnings []string
	ExitCode int
}

// NLRCCompiler implements the Compiler interface for NetLinx compiler
type NLRCCompiler struct {
	ExecutablePath string
	env            []string // override subprocess environment (tests only; nil = inherit)
}

// NewNLRCCompiler creates a new NLRC compiler instance
func NewNLRCCompiler(executablePath string) *NLRCCompiler {
	return &NLRCCompiler{
		ExecutablePath: executablePath,
	}
}

// Compile compiles NetLinx source files
func (c *NLRCCompiler) Compile(options CompileOptions) (*CompileResult, error) {
	args, err := c.BuildArgs(options)
	if err != nil {
		return nil, err
	}

	// Create the command with the executable path and arguments
	// Don't wrap the executable path in quotes - let Go handle it
	cmd := exec.Command(c.ExecutablePath, args...)
	if c.env != nil {
		cmd.Env = c.env
	}

	// Capture stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start compiler: %w", err)
	}

	// Read output
	output, err := c.readOutput(stdout, stderr)
	if err != nil {
		return nil, fmt.Errorf("failed to read output: %w", err)
	}

	// Wait for command to complete
	if err := cmd.Wait(); err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			result := &CompileResult{
				Success:  false,
				Output:   output,
				ExitCode: exitErr.ExitCode(),
			}

			result.Errors, result.Warnings = c.parseOutput(output)
			return result, nil
		}

		return nil, fmt.Errorf("compiler execution failed: %w", err)
	}

	result := &CompileResult{
		Success:  true,
		Output:   output,
		ExitCode: 0,
	}

	result.Errors, result.Warnings = c.parseOutput(output)

	return result, nil
}

// BuildArgs builds the command line arguments for the compiler
func (c *NLRCCompiler) BuildArgs(options CompileOptions) ([]string, error) {
	var args []string

	// Add source files first (they must come before options and be absolute paths)
	switch {
	case len(options.CFGFiles) > 0:
		for i, file := range options.CFGFiles {
			abs, err := filepath.Abs(file)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve path for %q: %w", file, err)
			}

			if _, err := os.Stat(abs); err != nil {
				return nil, fmt.Errorf("CFG file not found: %s", abs)
			}

			options.CFGFiles[i] = abs
		}

		cfgArg := cfgFlag + strings.Join(options.CFGFiles, pathDelimiter)
		args = append(args, cfgArg)
	case len(options.SourceFiles) > 0:
		for i, file := range options.SourceFiles {
			ext := strings.ToLower(filepath.Ext(file))
			if ext != ".axs" && ext != ".axi" {
				return nil, fmt.Errorf("invalid source file %q: must have .axs or .axi extension", file)
			}

			abs, err := filepath.Abs(file)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve path for %q: %w", file, err)
			}

			if _, err := os.Stat(abs); err != nil {
				return nil, fmt.Errorf("source file not found: %s", abs)
			}

			options.SourceFiles[i] = abs
		}

		args = append(args, options.SourceFiles...)
	}

	if joined := joinNonEmpty(options.IncludePath); joined != "" {
		args = append(args, includeFlag+"\""+joined+"\"")
	}

	if joined := joinNonEmpty(options.ModulePath); joined != "" {
		args = append(args, moduleFlag+"\""+joined+"\"")
	}

	if joined := joinNonEmpty(options.LibraryPath); joined != "" {
		args = append(args, libraryFlag+"\""+joined+"\"")
	}

	if options.OutputPath != "" {
		args = append(args, outputFlag+"\""+options.OutputPath+"\"")
	}

	return args, nil
}

// joinNonEmpty filters out empty strings and joins the remainder with pathDelimiter.
func joinNonEmpty(paths []string) string {
	var filtered []string

	for _, p := range paths {
		if p != "" {
			filtered = append(filtered, p)
		}
	}

	return strings.Join(filtered, pathDelimiter)
}

// readOutput streams stdout and stderr line-by-line as the compiler runs,
// printing every line immediately and also accumulating them for parseOutput.
func (c *NLRCCompiler) readOutput(stdout, stderr io.Reader) (string, error) {
	var output strings.Builder

	// Read stdout
	stdoutScanner := bufio.NewScanner(stdout)
	for stdoutScanner.Scan() {
		line := stdoutScanner.Text()
		output.WriteString(line + "\n")
		fmt.Println(line)
	}

	// Read stderr
	stderrScanner := bufio.NewScanner(stderr)
	for stderrScanner.Scan() {
		line := stderrScanner.Text()
		output.WriteString(line + "\n")
		fmt.Println(line)
	}

	if err := stdoutScanner.Err(); err != nil {
		return output.String(), err
	}

	if err := stderrScanner.Err(); err != nil {
		return output.String(), err
	}

	return output.String(), nil
}

// parseOutput parses compiler output for errors and warnings.
// It mirrors the TypeScript regex (?<level>ERROR|WARNING): .+ and deduplicates
// results, matching NLRC's output format exactly.
func (c *NLRCCompiler) parseOutput(output string) ([]string, []string) {
	seenErrors := make(map[string]struct{})
	seenWarnings := make(map[string]struct{})

	var errors, warnings []string

	for line := range strings.SplitSeq(output, "\n") {
		match := logPattern.FindString(line)
		if match == "" {
			continue
		}

		if strings.HasPrefix(match, "ERROR:") {
			if _, seen := seenErrors[match]; !seen {
				seenErrors[match] = struct{}{}
				errors = append(errors, match)
			}
		} else {
			if _, seen := seenWarnings[match]; !seen {
				seenWarnings[match] = struct{}{}
				warnings = append(warnings, match)
			}
		}
	}

	return errors, warnings
}
