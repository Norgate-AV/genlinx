import { checkbox } from "@inquirer/prompts";

export async function selectWorkspaceFiles(
    files: Array<string>,
): Promise<Array<string>> {
    return await checkbox<string>({
        message: "Select file(s):",
        choices: files.map((file) => ({ value: file })),
        validate: (answer) => {
            if (answer.length < 1) {
                return "You must choose at least one file.";
            }

            return true;
        },
    });
}

export default selectWorkspaceFiles;
