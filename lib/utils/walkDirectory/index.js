import fs from "fs-extra";
import path from "path";
import chalk from "chalk";

function ignoreEntity(entity) {
    return (
        entity.includes("node_modules") ||
        entity.includes(".git") ||
        entity.includes(".history")
    );
}

export async function walkDirectory(directory) {
    const files = [];
    let entities;

    try {
        entities = await fs.readdir(directory);
    } catch (error) {
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
        } catch (error) {
            console.log(chalk.red(error.message));
        }
    }

    return files;
}

export default walkDirectory;
