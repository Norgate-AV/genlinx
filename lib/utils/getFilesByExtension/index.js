import fs from "fs-extra";
import path from "path";

export async function getFilesByExtension(extension) {
    const entities = await fs.readdir(process.cwd());
    const files = entities.filter(
        (entity) => path.extname(entity) === extension,
    );

    return files;
}

export default getFilesByExtension;
