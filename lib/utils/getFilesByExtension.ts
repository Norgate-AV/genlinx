import fs from "node:fs/promises";
import path from "path";
import chalk from "chalk";

export async function getFilesByExtension(
    extension: string,
): Promise<Array<string>> {
    const files: Array<string> = [];

    try {
        const entities = await fs.readdir(process.cwd());

        files.push(
            ...entities.filter((entity) => path.extname(entity) === extension),
        );

        return files;
    } catch (error: any) {
        console.log(chalk.red(error.message));
    }

    return files;
}

export default getFilesByExtension;
