import inquirer from "inquirer";

export async function selectWorkspaceFiles(files) {
    const { selectedWorkspaceFiles } = await inquirer.prompt([
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
