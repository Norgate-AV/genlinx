import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { getGlobalConfig, getGlobalConfigFilePath } from "./getGlobalConfig.js";
import { getPackageJson } from "./index.js";
import { cosmiconfig } from "cosmiconfig";

// Mock dependencies
vi.mock("./index.js");
vi.mock("cosmiconfig");

function normalizePaths(
    paths: string | string[] | undefined,
): string | string[] | undefined {
    if (paths === undefined) {
        return undefined;
    }

    const normalize = (p: string) => {
        // First replace backslashes with forward slashes
        let normalized = p.replace(/\\/g, "/");

        // Replace common home directory patterns with /home/user/
        // 1. macOS: /Users/username/
        normalized = normalized.replace(/\/Users\/[^\/]+\//, "/home/user/");
        // 2. Windows: C:/Users/username/
        normalized = normalized.replace(/C:\/Users\/[^\/]+\//, "/home/user/");
        // 3. Linux: /home/username/ (but not if already /home/user/)
        normalized = normalized.replace(
            /\/home\/(?!user\/)[^\/]+\//,
            "/home/user/",
        );

        // Additionally strip drive letter if still present
        normalized = normalized.replace(/^[A-Z]:/, "");
        return normalized;
    };

    return Array.isArray(paths) ? paths.map(normalize) : normalize(paths);
}

describe("getGlobalConfigFilePath", () => {
    const originalEnv = process.env;

    beforeEach(() => {
        process.env = { ...originalEnv };
    });

    afterEach(() => {
        process.env = originalEnv;
    });

    it("should return default path when no environment variable is set", () => {
        const name = "genlinx";
        const result = getGlobalConfigFilePath(name);

        expect(normalizePaths(result)).toBe(
            "/home/user/.config/genlinx/config.json",
        );
    });

    it("should use GENLINX_CONFIG_DIR when set", () => {
        process.env.GENLINX_CONFIG_DIR = "/custom/config/dir";
        const name = "genlinx";
        const result = getGlobalConfigFilePath(name);

        expect(normalizePaths(result)).toBe("/custom/config/dir/config.json");
    });

    it("should use basename of name parameter", () => {
        const name = "/path/to/genlinx";
        const result = getGlobalConfigFilePath(name);

        expect(normalizePaths(result)).toBe(
            "/home/user/.config/genlinx/config.json",
        );
    });
});

describe("getGlobalConfig", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("should throw if package.json is not found", async () => {
        vi.mocked(getPackageJson).mockResolvedValue(null);

        await expect(getGlobalConfig()).rejects.toThrow(
            "package.json not found",
        );
    });

    it("should throw if package.json has no name field", async () => {
        vi.mocked(getPackageJson).mockResolvedValue({});

        await expect(getGlobalConfig()).rejects.toThrow(
            "package.json missing name field",
        );
    });

    it("should load the global config file when it exists", async () => {
        vi.mocked(getPackageJson).mockResolvedValue({ name: "genlinx" });

        const mockResult = {
            filepath: "/home/user/.config/genlinx/config.json",
            config: {
                cfg: {
                    includePath: ["global/include"],
                },
            },
            isEmpty: false,
        };

        // Make sure the mock structure matches exactly what getGlobalConfig expects
        const mockLoad = vi.fn().mockResolvedValue(mockResult);
        vi.mocked(cosmiconfig).mockReturnValue({
            load: mockLoad,
            search: vi.fn(),
        } as any);

        const result = await getGlobalConfig();

        expect(mockLoad).toHaveBeenCalled();
        const calledPath = mockLoad.mock.calls?.[0]?.[0] ?? "";
        expect(normalizePaths(calledPath)).toBe(
            "/home/user/.config/genlinx/config.json",
        );

        expect(result).toEqual(mockResult);
    });

    it("should return null when no config file exists", async () => {
        vi.mocked(getPackageJson).mockResolvedValue({ name: "genlinx" });

        const mockLoad = vi.fn().mockRejectedValue(new Error("File not found"));
        vi.mocked(cosmiconfig).mockReturnValue({ load: mockLoad } as any);

        const result = await getGlobalConfig();

        expect(mockLoad).toHaveBeenCalled();
        expect(result).toBeNull();
    });

    it("should return null on any error during loading", async () => {
        vi.mocked(getPackageJson).mockResolvedValue({ name: "genlinx" });

        const mockLoad = vi
            .fn()
            .mockRejectedValue(new Error("Some unexpected error"));
        vi.mocked(cosmiconfig).mockReturnValue({ load: mockLoad } as any);

        const result = await getGlobalConfig();

        expect(mockLoad).toHaveBeenCalled();
        expect(result).toBeNull();
    });
});
