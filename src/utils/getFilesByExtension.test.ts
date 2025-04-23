import fs from "node:fs/promises";
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { getFilesByExtension } from "./getFilesByExtension.js";

// Mock the fs module
vi.mock("node:fs/promises");

describe("getFilesByExtension", () => {
    const mockFiles = [
        "file1.axs",
        "file2.axs",
        "file3.axi",
        "file4.cfg",
        "file5.txt",
    ];

    beforeEach(() => {
        // Mock the readdir function to return our mock files
        vi.mocked(fs.readdir).mockResolvedValue(mockFiles as any);
    });

    afterEach(() => {
        vi.resetAllMocks();
    });

    it("should return files with the specified extension", async () => {
        const result = await getFilesByExtension(".axs");
        expect(result).toHaveLength(2);
        expect(result).toContain("file1.axs");
        expect(result).toContain("file2.axs");
        expect(result).not.toContain("file3.axi");
    });

    it("should return empty array when no files match", async () => {
        const result = await getFilesByExtension(".xyz");
        expect(result).toHaveLength(0);
    });

    it("should handle errors gracefully", async () => {
        vi.mocked(fs.readdir).mockRejectedValue(new Error("Access denied"));
        const result = await getFilesByExtension(".axs");
        expect(result).toHaveLength(0);
    });
});
