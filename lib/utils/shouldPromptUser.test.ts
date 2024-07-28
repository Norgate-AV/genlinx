import { describe, expect, it } from "vitest";
import { shouldPromptUser } from "./shouldPromptUser.js";

describe("shouldPromptUser", () => {
    it("should return false when options.all is true", () => {
        const options = { all: true };
        const files = ["file1", "file2"];

        expect(shouldPromptUser(options, files)).toBe(false);
    });

    it("should return false when there is only one file", () => {
        const options = { all: false };
        const files = ["file1"];

        expect(shouldPromptUser(options, files)).toBe(false);
    });

    it("should return true when options.all is false and there are multiple files", () => {
        const options = { all: false };
        const files = ["file1", "file2"];

        expect(shouldPromptUser(options, files)).toBe(true);
    });
});
