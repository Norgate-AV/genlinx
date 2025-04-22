import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { getLocalConfig } from "./getLocalConfig.js";
import { getPackageJson } from "./index.js";
import { cosmiconfig } from "cosmiconfig";
import type { CliOptions } from "../@types/index.js";

// Mock dependencies
vi.mock("./index.js");
vi.mock("cosmiconfig");
vi.mock("chalk", () => ({
    default: {
        blue: vi.fn((text) => text),
    },
}));

// Mock node:path to ensure consistent cross-platform behavior
vi.mock("node:path", async (importOriginal) => {
    const actualPath = (await importOriginal()) as typeof import("node:path");

    return {
        ...actualPath,
        resolve: vi.fn((root, p) => `${root}/${p}`),
        basename: vi.fn((p) => p.split(/[\\/]/).pop() || ""),
        dirname: vi.fn((p) => {
            const parts = p.split(/[\\/]/);
            parts.pop();
            return parts.join("/") || ".";
        }),
    };
});

// Helper function to normalize paths in test assertions
function normalizePaths(
    paths: string | string[] | undefined,
): string | string[] | undefined {
    if (paths === undefined) {
        return undefined;
    }

    // Normalize function that handles both backslashes and drive letters
    const normalize = (p: string) =>
        p.replace(/\\/g, "/").replace(/^[A-Z]:/, "");

    return Array.isArray(paths) ? paths.map(normalize) : normalize(paths);
}

describe("getLocalConfig", () => {
    const consoleSpy = vi.spyOn(console, "log").mockImplementation(() => {});

    beforeEach(() => {
        vi.resetAllMocks();
    });

    afterEach(() => {
        consoleSpy.mockClear();
    });

    it("should throw if package.json is not found", async () => {
        // Setup: getPackageJson returns null (no package.json)
        vi.mocked(getPackageJson).mockResolvedValue(null);

        // Test: getLocalConfig should throw
        await expect(getLocalConfig({})).rejects.toThrow(
            "package.json not found",
        );
    });

    it("should throw if package.json has no name field", async () => {
        // Setup: getPackageJson returns a package.json without name
        vi.mocked(getPackageJson).mockResolvedValue({});

        // Test: getLocalConfig should throw
        await expect(getLocalConfig({})).rejects.toThrow(
            "package.json missing name field",
        );
    });

    it("should return null if no config file is found", async () => {
        // Setup
        vi.mocked(getPackageJson).mockResolvedValue({ name: "genlinx" });
        const mockSearch = vi.fn().mockResolvedValue(null);
        vi.mocked(cosmiconfig).mockReturnValue({ search: mockSearch } as any);

        // Test
        const result = await getLocalConfig({});

        // Assert
        expect(result).toBeNull();
        expect(cosmiconfig).toHaveBeenCalledWith("genlinx");
        expect(mockSearch).toHaveBeenCalled();
    });

    it("should resolve paths when config file is found", async () => {
        // Setup
        vi.mocked(getPackageJson).mockResolvedValue({ name: "genlinx" });

        const mockConfig = {
            cfg: {
                includePath: ["include", "another/include"],
                modulePath: ["module"],
                libraryPath: ["library"],
            },
            archive: {
                extraFileSearchLocations: ["archive"],
            },
            build: {
                nlrc: {
                    includePath: ["nlrc/include"],
                    modulePath: ["nlrc/module"],
                    libraryPath: ["nlrc/library"],
                },
            },
        };

        const mockResult = {
            filepath: "/path/to/config",
            config: { ...mockConfig },
            isEmpty: false,
        };

        const mockSearch = vi.fn().mockResolvedValue(mockResult);
        vi.mocked(cosmiconfig).mockReturnValue({ search: mockSearch } as any);

        // Test
        const result = await getLocalConfig({});

        // Assert
        expect(result).not.toBeNull();
        expect(result?.filepath).toBe("/path/to/config");

        // Check path resolution
        expect(normalizePaths(result?.config.cfg.includePath)).toEqual([
            "/path/to/include",
            "/path/to/another/include",
        ]);
        expect(normalizePaths(result?.config.cfg.modulePath)).toEqual([
            "/path/to/module",
        ]);
        expect(normalizePaths(result?.config.cfg.libraryPath)).toEqual([
            "/path/to/library",
        ]);
        expect(
            normalizePaths(result?.config.archive.extraFileSearchLocations),
        ).toEqual(["/path/to/archive"]);
        expect(normalizePaths(result?.config.build.nlrc.includePath)).toEqual([
            "/path/to/nlrc/include",
        ]);
        expect(normalizePaths(result?.config.build.nlrc.modulePath)).toEqual([
            "/path/to/nlrc/module",
        ]);
        expect(normalizePaths(result?.config.build.nlrc.libraryPath)).toEqual([
            "/path/to/nlrc/library",
        ]);
    });

    it("should handle missing path arrays in config gracefully", async () => {
        // Setup with minimal config
        vi.mocked(getPackageJson).mockResolvedValue({ name: "genlinx" });

        const mockResult = {
            filepath: "/path/to/config",
            config: {
                cfg: {},
            },
            isEmpty: false,
        };

        const mockSearch = vi.fn().mockResolvedValue(mockResult);
        vi.mocked(cosmiconfig).mockReturnValue({ search: mockSearch } as any);

        // Test
        const result = await getLocalConfig({});

        // Assert - should not throw and return config as is
        expect(result).not.toBeNull();
        expect(result?.config.cfg).toBeDefined();
    });

    it("should log config filepath when verbose is enabled", async () => {
        // Setup
        vi.mocked(getPackageJson).mockResolvedValue({ name: "genlinx" });

        const mockResult = {
            filepath: "/path/to/config",
            config: { cfg: {} },
            isEmpty: false,
        };

        const mockSearch = vi.fn().mockResolvedValue(mockResult);
        vi.mocked(cosmiconfig).mockReturnValue({ search: mockSearch } as any);

        const options: CliOptions = {
            cfg: {
                verbose: true,
                includePath: [],
                workspaceFiles: [],
                rootDirectory: "",
            },
        };

        // Test
        await getLocalConfig(options);

        // Assert
        expect(consoleSpy).toHaveBeenCalledWith(
            expect.stringContaining(
                "Using local config file at /path/to/config",
            ),
        );
    });

    it("should not log when verbose is disabled", async () => {
        // Setup
        vi.mocked(getPackageJson).mockResolvedValue({ name: "genlinx" });

        const mockResult = {
            filepath: "/path/to/config",
            config: { cfg: {} },
            isEmpty: false,
        };

        const mockSearch = vi.fn().mockResolvedValue(mockResult);
        vi.mocked(cosmiconfig).mockReturnValue({ search: mockSearch } as any);

        const options: CliOptions = {
            cfg: {
                verbose: false,
                includePath: [],
                workspaceFiles: [],
                rootDirectory: "",
            },
        };

        // Test
        await getLocalConfig(options);

        // Assert
        expect(consoleSpy).not.toHaveBeenCalled();
    });
});
