import fs from "fs-extra";
import path from "path";
import chalk from "chalk";

export function walkDirectory(directory) {
    const files = [];

    const entities = fs.readdirSync(directory);

    for (const entity of entities) {
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
