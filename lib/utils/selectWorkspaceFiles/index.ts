import inquirer from "inquirer";

export async function selectWorkspaceFiles(
    files: Array<string>,
): Promise<Array<string>> {
    const {
        selectedWorkspaceFiles,
    }: { selectedWorkspaceFiles: Array<string> } = await inquirer.prompt([
        {
            type: files.length > 1 ? "checkbox" : "confirm",
            name: "selectedWorkspaceFiles",
            message: "Select file(s):",
            choices: files,
        },
    ]);

    return selectedWorkspaceFiles;
}

export default selectWorkspaceFiles;
