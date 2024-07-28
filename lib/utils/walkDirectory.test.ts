import { describe, expect, it } from "vitest";
import { walkDirectory } from "./walkDirectory.js";

describe("walkDirectory", () => {
    it("should return an array of files", async () => {
        const files = await walkDirectory("./lib");

        expect(files).toBeInstanceOf(Array);
        expect(files.length).toBeGreaterThan(0);
    });

    it("should return an empty array if the directory does not exist", async () => {
        const files = await walkDirectory("./does-not-exist");

        expect(files).toBeInstanceOf(Array);
        expect(files.length).toBe(0);
    });
});
