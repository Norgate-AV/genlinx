import inquirer from "inquirer";

export async function selectFiles(files) {
    const { selectedFiles } = await inquirer.prompt([
        {
            type: "checkbox",
            name: "selectedFiles",
            message: "Select file(s):",
            choices: files,
            validate: (answer) => {
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
