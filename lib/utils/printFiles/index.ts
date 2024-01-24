import chalk from "chalk";

export function printFiles(files: Array<string>) {
    console.log(chalk.green(`Found ${files.length} file(s)`));

    for (const file of files) {
        console.log(chalk.cyan(`--> ${file}`));
    }
}

export default printFiles;
