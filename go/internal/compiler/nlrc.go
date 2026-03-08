package compiler

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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
	if len(options.CFGFiles) > 0 {
		// Combine CFG files with -C flag
		cfgArg := "-C" + strings.Join(options.CFGFiles, ";")
		args = append(args, cfgArg)
	} else if len(options.SourceFiles) > 0 {
		// Convert source files to absolute paths if they exist
		for i, file := range options.SourceFiles {
			if path, err := filepath.Abs(file); err == nil {
				// Check if the file exists or if it's an absolute path already
				if _, err := os.Stat(file); err == nil || filepath.IsAbs(file) {
					options.SourceFiles[i] = path
				}
				// If file doesn't exist and is relative, keep the original path
			}
			// If we can't get absolute path, keep the original
		}

		args = append(args, options.SourceFiles...)
	}

	// Add include paths (flag with quoted path)
	for _, path := range options.IncludePath {
		// Only add include path if it exists and is not empty
		if path != "" {
			args = append(args, "-I\""+path+"\"")
		}
	}

	// Add module paths as semicolon-joined list with -M flag
	if len(options.ModulePath) > 0 {
		combinedModulePath := strings.Join(options.ModulePath, ";")
		args = append(args, "-M\""+combinedModulePath+"\"")
	}

	// Add library paths (flag with quoted path)
	for _, path := range options.LibraryPath {
		if path != "" {
			args = append(args, "-L\""+path+"\"")
		}
	}

	// Add output path if specified
	if options.OutputPath != "" {
		args = append(args, "-O\""+options.OutputPath+"\"")
	}

	return args, nil
}

// readOutput reads from stdout and stderr pipes
func (c *NLRCCompiler) readOutput(stdout, stderr io.Reader) (string, error) {
	var output strings.Builder

	// Read stdout
	stdoutScanner := bufio.NewScanner(stdout)
	for stdoutScanner.Scan() {
		line := stdoutScanner.Text()
		output.WriteString(line + "\n")
		if c.ExecutablePath != "" { // Assuming verbose flag check
			fmt.Println(line)
		}
	}

	// Read stderr
	stderrScanner := bufio.NewScanner(stderr)
	for stderrScanner.Scan() {
		line := stderrScanner.Text()
		output.WriteString(line + "\n")
		if c.ExecutablePath != "" { // Assuming verbose flag check
			fmt.Println(line)
		}
	}

	if err := stdoutScanner.Err(); err != nil {
		return output.String(), err
	}

	if err := stderrScanner.Err(); err != nil {
		return output.String(), err
	}

	return output.String(), nil
}

// parseOutput parses compiler output for errors and warnings
func (c *NLRCCompiler) parseOutput(output string) ([]string, []string) {
	var errors, warnings []string
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check for error patterns (customize based on actual compiler output)
		if strings.Contains(strings.ToLower(line), "error") {
			errors = append(errors, line)
		} else if strings.Contains(strings.ToLower(line), "warning") {
			warnings = append(warnings, line)
		}
	}

	return errors, warnings
}
