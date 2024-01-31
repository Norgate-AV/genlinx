import inquirer from "inquirer";

export async function selectFiles(
    files: Array<string>,
): Promise<Array<string>> {
    const { selectedFiles }: { selectedFiles: Array<string> } =
        await inquirer.prompt([
            {
                type: "checkbox",
                name: "selectedFiles",
                message: "Select file(s):",
                choices: files,
                validate: (answer: Array<string>) => {
                    if (answer.length < 1) {
                        return "You must choose at least one file.";
                    }

                    return true;
                },
            },
        ]);

    return selectedFiles;
}

export default selectFiles;
