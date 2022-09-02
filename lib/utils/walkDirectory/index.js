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

export function walkDirectory(directory) {
    const files = [];

    const entities = fs.readdirSync(directory);

    for (const entity of entities) {
        if (ignoreEntity(entity)) {
            continue;
        }

        const entityPath = path.join(directory, entity);

        try {
            const stats = fs.lstatSync(entityPath);

            if (stats.isDirectory()) {
                files.push(...walkDirectory(entityPath));
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
