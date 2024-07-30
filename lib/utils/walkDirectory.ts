import fs from "node:fs/promises";
import path from "path";
import chalk from "chalk";

function ignoreEntity(entity: string): boolean {
    return (
        entity.includes("node_modules") ||
        entity.includes(".git") ||
        entity.includes(".history")
    );
}

export async function walkDirectory(directory: string): Promise<Array<string>> {
    const files: Array<string> = [];
    let entities: Array<string>;

    try {
        entities = await fs.readdir(directory);
    } catch (error: any) {
        console.log(chalk.red(error.message));
        return files;
    }

    for (const entity of entities) {
        if (ignoreEntity(entity)) {
            continue;
        }

        const entityPath = path.join(directory, entity);

        try {
            const stats = await fs.lstat(entityPath);

            if (stats.isDirectory()) {
                files.push(...(await walkDirectory(entityPath)));
            } else {
                files.push(entityPath);
            }
        } catch (error: any) {
            console.log(chalk.red(error.message));
        }
    }

    return files;
}

export default walkDirectory;
