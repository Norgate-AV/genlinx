import fs from "fs-extra";
import path from "path";
import { APW } from "../..";

export async function getWorkspaceFiles() {
    const entities = await fs.readdir(process.cwd());
    const workspaceFiles = entities.filter(
        (entity) =>
            path.extname(entity) === APW.fileExtensions[APW.fileType.workspace],
    );

    return workspaceFiles;
}

export default getWorkspaceFiles;
