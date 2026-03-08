import fs from "node:fs/promises";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { APW } from "./APW.js";
import * as pathUtils from "../../utils/pathExists.js";

// Mock the fs module and pathExists utility
vi.mock("node:fs/promises");
vi.mock("../../utils/pathExists", () => ({
    pathExists: vi.fn(),
}));

describe("APW", () => {
    const mockWorkspaceData = `<!DOCTYPE Workspace [
      <Workspace Version="4">
        <Identifier>WORKSPACE-ID-12345</Identifier>
        <FileManager Default="">
          <File GUID="{12345-ABCDE}" Type="MasterSrc" Path="test.axs">
            <Identifier>FILE-ID-12345</Identifier>
          </File>
        </FileManager>
      </Workspace>
    `;

    beforeEach(() => {
        vi.mocked(pathUtils.pathExists).mockResolvedValue(true);
        vi.mocked(fs.readFile).mockResolvedValue(mockWorkspaceData as any);
    });

    it("should initialize with correct file path", () => {
        const apw = new APW("test.apw");
        expect(apw.filePath).toContain("test.apw");
    });

    it("should load APW file successfully", async () => {
        const apw = new APW("test.apw");
        await apw.load();

        expect(fs.readFile).toHaveBeenCalled();
        // We would test for more properties but they're private
        // This is where you might consider making some properties or methods public for testing
    });

    it("should throw error if file does not exist", async () => {
        vi.mocked(pathUtils.pathExists).mockResolvedValue(false);

        const apw = new APW("nonexistent.apw");
        await expect(apw.load()).rejects.toThrow("Failed to load APW file");
    });
});
