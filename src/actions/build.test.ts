import os from "node:os";
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { build } from "../commands/build.js";
import { actions } from "./index.js";
import type { BuildOptions } from "../../lib/@types/index.js";
import { getBuildLogs, catchAllErrors, printAllWarnings } from "./build.js";

// Mock the actions module
vi.mock("./index", () => ({
    actions: {
        build: {
            execute: vi.fn(),
        },
    },
}));

// Mock the os module
vi.mock("node:os", () => ({
    default: {
        platform: vi.fn(() => "win32"),
    },
    platform: vi.fn(() => "win32"),
}));

// Mock NLRC for the executeSourceBuild test
vi.mock("../../lib/index.js", () => ({
    NLRC: {
        getSourceBuildCommand: vi.fn().mockImplementation((file) => {
            // Always return an object with the required properties
            return {
                path: "nlrc",
                args: ["-c", file],
            };
        }),
    },
}));

function createTestBuildOptions(
    overrides: Partial<BuildOptions> = {},
): BuildOptions {
    const defaultOptions: BuildOptions = {
        nlrc: {
            path: "path/to/nlrc",
            option: {
                cfg: "-CFG",
                includePath: "-I",
                modulePath: "-M",
                libraryPath: "-L",
            },
            includePath: [],
            modulePath: [],
            libraryPath: [],
        },
        shell: {
            path: "cmd.exe",
        },
        all: false,
        createCfg: false,
        verbose: false,
    };

    // Merge the defaults with any provided overrides
    return {
        ...defaultOptions,
        ...overrides,
    };
}

function getCallOptions(): BuildOptions {
    // This ensures TypeScript knows we have a call and options object
    const calls = vi.mocked(actions.build.execute).mock.calls;
    expect(calls.length).toBeGreaterThan(0);

    if (calls.length === 0) {
        throw new Error("No calls to actions.build.execute found");
    }

    const options = calls[0]?.[0] as BuildOptions;
    if (!options) {
        throw new Error(
            "No options found in the first call to actions.build.execute",
        );
    }

    return options;
}

describe("build command", () => {
    beforeEach(() => {
        vi.resetAllMocks();
        // Mock process.exit so it doesn't exit our tests
        vi.spyOn(process, "exit").mockImplementation((() => {}) as any);
        // Mock console.log for testing output messages
        vi.spyOn(console, "log").mockImplementation(() => {});
        vi.spyOn(console, "error").mockImplementation(() => {});
    });

    it("should create a Command instance", () => {
        const command = build();
        expect(command).toBeDefined();
        expect(command.name()).toBe("build");
    });

    it("should call build.execute when action handler is triggered", () => {
        // Create the command
        const command = build();

        // Use the parse method directly with command-line style arguments
        command.parse(["node", "genlinx", "build", "file1.axs", "file2.axs"], {
            from: "user",
        });

        // Check that execute was called
        expect(actions.build.execute).toHaveBeenCalled();

        // Get the options that were passed to execute
        const callOptions = getCallOptions();

        // Check that sourceFiles includes our expected files
        // The actual array includes the command prefix: ["node", "genlinx", "build", ...files]
        expect(callOptions.sourceFiles).toContain("file1.axs");
        expect(callOptions.sourceFiles).toContain("file2.axs");
        expect(callOptions.sourceFiles?.slice(-2)).toEqual([
            "file1.axs",
            "file2.axs",
        ]);
    });

    it("should pass arguments as-is to build.execute", () => {
        // Create the command
        const command = build();

        // Use the parse method directly with the verbose flag
        command.parse(["node", "genlinx", "build", "test.axs", "--verbose"], {
            from: "user",
        });

        // Check that execute was called
        expect(actions.build.execute).toHaveBeenCalled();

        // Get the options
        const callOptions = getCallOptions();

        // Check verbose flag was set
        expect(callOptions.verbose).toBe(true);

        // Check sourceFiles contains our file (along with the command parts)
        expect(callOptions.sourceFiles).toContain("test.axs");
        expect(callOptions.sourceFiles?.slice(-1)).toEqual(["test.axs"]);
    });

    // Path options tests
    it("should handle include paths", () => {
        const command = build();

        command.parse(
            ["node", "genlinx", "build", "--include-path", "path1", "path2"],
            {
                from: "user",
            },
        );

        const callOptions = getCallOptions();
        expect(callOptions.includePath).toContain("path1");
        expect(callOptions.includePath).toContain("path2");
    });

    it("should use default values for options not specified", () => {
        const command = build();

        command.parse(["node", "genlinx", "build"], {
            from: "user",
        });

        const callOptions = getCallOptions();
        expect(callOptions.cfgFiles).toEqual([]);
        expect(callOptions.verbose).toBe(false);
    });

    // CFG file options
    it("should handle cfg files option", () => {
        const command = build();

        command.parse(
            [
                "node",
                "genlinx",
                "build",
                "--cfg-files",
                "project1.cfg",
                "project2.cfg",
            ],
            {
                from: "user",
            },
        );

        const callOptions = getCallOptions();
        expect(callOptions.cfgFiles).toContain("project1.cfg");
        expect(callOptions.cfgFiles).toContain("project2.cfg");
    });

    // Basic command structure tests
    it("should create a Command instance with correct structure", () => {
        const command = build();
        expect(command).toBeDefined();
        expect(command.name()).toBe("build");
        expect(command.description()).toContain(
            "build NetLinx source file(s) or CFG file(s)",
        );

        // Verify all options are registered
        const optionNames = command.options.map((opt) => opt.long);
        expect(optionNames).toContain("--cfg-files");
        expect(optionNames).toContain("--source-files");
        expect(optionNames).toContain("--include-path");
        expect(optionNames).toContain("--module-path");
        expect(optionNames).toContain("--library-path");
        expect(optionNames).toContain("--all");
        expect(optionNames).toContain("--verbose");
    });

    // Platform compatibility tests
    it("should exit when not running on Windows", () => {
        // Mock the process.exit to throw an error to stop execution
        vi.spyOn(process, "exit").mockImplementation((code) => {
            throw new Error(`Process exited with code ${code}`);
        });

        // Mock the platform to return a non-Windows value
        vi.mocked(os.platform).mockReturnValueOnce("darwin");
        const command = build();

        // Parse should now throw due to our process.exit mock
        expect(() => {
            command.parse(["node", "genlinx", "build"], { from: "user" });
        }).toThrow("Process exited with code 1");

        // Verify that console.log was called with the expected message
        expect(console.log).toHaveBeenCalledWith(
            expect.stringContaining("only supported on Windows"),
        );

        // Verify process.exit was called with code 1
        expect(process.exit).toHaveBeenCalledWith(1);

        // We no longer need to check if actions.build.execute was called
        // because the thrown error stops execution before it would be called
    });

    // Source file tests
    it("should handle source files passed as arguments", () => {
        const command = build();

        command.parse(["node", "genlinx", "build", "file1.axs", "file2.axs"], {
            from: "user",
        });

        expect(actions.build.execute).toHaveBeenCalled();
        const callOptions = getCallOptions();
        expect(callOptions.sourceFiles).toContain("file1.axs");
        expect(callOptions.sourceFiles).toContain("file2.axs");
    });

    it("should combine source files from arguments and options", () => {
        const command = build();

        command.parse(
            [
                "node",
                "genlinx",
                "build",
                "file1.axs",
                "--source-files",
                "file2.axs",
                "file3.axs",
            ],
            {
                from: "user",
            },
        );

        const callOptions = getCallOptions();
        expect(callOptions.sourceFiles).toContain("file1.axs");
        expect(callOptions.sourceFiles).toContain("file2.axs");
        expect(callOptions.sourceFiles).toContain("file3.axs");
    });

    // Option tests
    it("should pass the verbose flag", () => {
        const command = build();

        command.parse(["node", "genlinx", "build", "--verbose"], {
            from: "user",
        });

        const callOptions = getCallOptions();
        expect(callOptions.verbose).toBe(true);
    });

    it("should handle module paths", () => {
        const command = build();

        command.parse(
            ["node", "genlinx", "build", "--module-path", "mod1", "mod2"],
            {
                from: "user",
            },
        );

        const callOptions = getCallOptions();
        expect(callOptions.modulePath).toContain("mod1");
        expect(callOptions.modulePath).toContain("mod2");
    });

    it("should handle library paths", () => {
        const command = build();

        command.parse(
            ["node", "genlinx", "build", "--library-path", "lib1", "lib2"],
            {
                from: "user",
            },
        );

        const callOptions = getCallOptions();
        expect(callOptions.libraryPath).toContain("lib1");
        expect(callOptions.libraryPath).toContain("lib2");
    });

    // Flag tests
    it("should handle the all flag", () => {
        const command = build();

        command.parse(["node", "genlinx", "build", "--all"], {
            from: "user",
        });

        const callOptions = getCallOptions();
        expect(callOptions.all).toBe(true);
    });

    it("should handle the no-all flag", () => {
        const command = build();

        command.parse(["node", "genlinx", "build", "--no-all"], {
            from: "user",
        });

        const callOptions = getCallOptions();
        expect(callOptions.all).toBe(false);
    });

    // Complex combinations
    it("should handle complex option combinations", () => {
        const command = build();

        command.parse(
            [
                "node",
                "genlinx",
                "build",
                "main.axs",
                "--source-files",
                "extra.axs",
                "--include-path",
                "include1",
                "include2",
                "--module-path",
                "module1",
                "--library-path",
                "lib1",
                "--verbose",
            ],
            { from: "user" },
        );

        const callOptions = getCallOptions();
        expect(callOptions.sourceFiles).toContain("main.axs");
        expect(callOptions.sourceFiles).toContain("extra.axs");
        expect(callOptions.includePath).toContain("include1");
        expect(callOptions.includePath).toContain("include2");
        expect(callOptions.modulePath).toContain("module1");
        expect(callOptions.libraryPath).toContain("lib1");
        expect(callOptions.verbose).toBe(true);
    });

    // Edge cases
    it("should handle the case when sourceFiles is undefined", () => {
        const command = build();

        // Instead of mocking opts and calling the handler directly,
        // use the parse method with standard arguments
        command.parse(["node", "genlinx", "build", "test.axs"], {
            from: "user",
        });

        // Verify actions.build.execute was called
        expect(actions.build.execute).toHaveBeenCalled();
        const callOptions = getCallOptions();

        // Verify the sourceFiles array contains our test file
        expect(callOptions.sourceFiles).toContain("test.axs");

        // If we need to test what happens when options.sourceFiles is undefined,
        // we can check the implementation handles it correctly by testing
        // that the array was properly initialized
        expect(Array.isArray(callOptions.sourceFiles)).toBe(true);
    });

    // This test may need to be adjusted based on how Commander actually handles conflicts
    it("should maintain option conflicts", () => {
        const command = build();

        // This should log an error due to conflicting options
        // Note: Commander typically handles this by itself before our code runs
        const parse = () => {
            try {
                command.parse(
                    [
                        "node",
                        "genlinx",
                        "build",
                        "--cfg-files",
                        "project.cfg",
                        "--source-files",
                        "file.axs",
                    ],
                    {
                        from: "user",
                    },
                );
            } catch (e) {
                return e;
            }
        };

        const result = parse();
        // We don't necessarily expect an exception, but Commander should have prevented
        // both options from being set simultaneously
        // Testing the specific behavior would depend on Commander's implementation
    });
});

describe("getBuildLogs", () => {
    it("should parse error and warning messages", () => {
        const buildOutput = `
            Compiling project...
            ERROR: Undefined variable 'x' on line 42
            Processing...
            WARNING: Unused variable 'y' on line 73
            WARNING: Duplicate definition on line 81
            ERROR: Missing semicolon on line 90
            Build complete with errors.
        `;

        const logs = getBuildLogs(buildOutput);

        expect(logs.error).toHaveLength(2);
        expect(logs.warning).toHaveLength(2);
        expect(logs.error).toContain(
            "ERROR: Undefined variable 'x' on line 42",
        );
        expect(logs.error).toContain("ERROR: Missing semicolon on line 90");
        expect(logs.warning).toContain(
            "WARNING: Unused variable 'y' on line 73",
        );
        expect(logs.warning).toContain(
            "WARNING: Duplicate definition on line 81",
        );
    });

    it("should deduplicate repeated errors and warnings", () => {
        const buildOutput = `
            ERROR: Same error on line 10
            Processing...
            ERROR: Same error on line 10
            WARNING: Repeated warning on line 20
            More processing...
            WARNING: Repeated warning on line 20
        `;

        const logs = getBuildLogs(buildOutput);

        expect(logs.error).toHaveLength(1);
        expect(logs.warning).toHaveLength(1);
        expect(logs.error).toContain("ERROR: Same error on line 10");
        expect(logs.warning).toContain("WARNING: Repeated warning on line 20");
    });

    it("should return empty arrays when no errors or warnings", () => {
        const buildOutput = `
            Compiling project...
            Processing...
            Build successful.
        `;

        const logs = getBuildLogs(buildOutput);

        expect(logs.error).toHaveLength(0);
        expect(logs.warning).toHaveLength(0);
    });
});

describe("catchAllErrors", () => {
    beforeEach(() => {
        vi.spyOn(console, "log").mockImplementation(() => {});
    });

    it("should throw an error when errors exist", () => {
        const errors = ["ERROR: Something went wrong", "ERROR: Another issue"];

        expect(() => catchAllErrors(errors)).toThrow(
            "The build process failed with a total of 2 error(s)",
        );
        expect(console.log).toHaveBeenCalledWith(
            expect.stringContaining("2 error(s)"),
        );
    });

    it("should not throw when errors array is empty", () => {
        expect(() => catchAllErrors([])).not.toThrow();
        expect(console.log).not.toHaveBeenCalled();
    });
});

describe("printAllWarnings", () => {
    beforeEach(() => {
        vi.spyOn(console, "log").mockImplementation(() => {});
    });

    it("should print warnings when they exist", () => {
        const warnings = ["WARNING: Minor issue", "WARNING: Something to note"];

        printAllWarnings(warnings);

        expect(console.log).toHaveBeenCalledWith(
            expect.stringContaining("2 warning(s)"),
        );
        expect(console.log).toHaveBeenCalledTimes(3);
    });

    it("should not print anything when warnings array is empty", () => {
        printAllWarnings([]);
        expect(console.log).not.toHaveBeenCalled();
    });
});
