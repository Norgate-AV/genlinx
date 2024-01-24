import fs from "fs-extra";
import path from "path";
import chalk from "chalk";

export async function getFilesByExtension(extension: string) {
    const files: Array<string> = [];

    try {
        const entities = await fs.readdir(process.cwd());

        files.push(
            ...entities.filter((entity) => path.extname(entity) === extension),
        );

        return files;
    } catch (error) {
        console.log(chalk.red(error.message));
    }

    return files;
}

export default getFilesByExtension;
