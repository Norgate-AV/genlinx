import { describe, expect, it, vi } from "vitest";
import * as inquirer from "@inquirer/prompts";
import { selectWorkspaceFiles } from "./selectWorkspaceFiles.js";

vi.mock("inquirer");

describe("selectWorkspaceFiles", () => {
    it("should call inquirer.checkbox with the correct arguments", async () => {
        vi.mock("inquirer");

        describe("selectWorkspaceFiles", () => {
            it("should call inquirer.checkbox with the correct arguments", async () => {
                const files = ["file1", "file2"];
                const checkboxSpy = vi.spyOn(inquirer, "checkbox");

                await selectWorkspaceFiles(files);

                expect(checkboxSpy).toHaveBeenCalledWith({
                    message: "Select file(s):",
                    choices: [{ value: "file1" }, { value: "file2" }],
                    validate: expect.any(Function),
                });
            });

            it("should return selected files when inquirer.checkbox resolves with selected choices", async () => {
                const files = ["file1", "file2"];
                const selectedChoices = ["file1"];
                vi.spyOn(inquirer, "checkbox").mockResolvedValue(
                    selectedChoices,
                );

                const result = await selectWorkspaceFiles(files);

                expect(result).toEqual(selectedChoices);
            });

            it("should return an empty array when inquirer.checkbox resolves with no selected choices", async () => {
                const files = ["file1", "file2"];
                const selectedChoices: Array<string> = [];
                vi.spyOn(inquirer, "checkbox").mockResolvedValue(
                    selectedChoices,
                );

                const result = await selectWorkspaceFiles(files);

                expect(result).toEqual(selectedChoices);
            });

            it("should throw an error when inquirer.checkbox rejects with an error", async () => {
                const files = ["file1", "file2"];
                const error = new Error("Checkbox error");
                vi.spyOn(inquirer, "checkbox").mockRejectedValue(error);

                await expect(selectWorkspaceFiles(files)).rejects.toThrow(
                    error,
                );
            });
        });
    });
});
