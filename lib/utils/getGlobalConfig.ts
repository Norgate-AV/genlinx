import path from "path";
import os from "os";
import { cosmiconfig, CosmiconfigResult } from "cosmiconfig";
import { GlobalConfig } from "../@types/index.js";
import pkg from "../../package.json";
import { get } from "config";

export function getGlobalConfigFilePath(): string {
    const file = "config.json";
    const directory =
        process.env.GENLINX_CONFIG_DIR ||
        path.join(os.homedir(), ".config", path.basename(pkg.name));

    return path.join(directory, file);
}

export async function getGlobalConfig(): Promise<CosmiconfigResult> {
    return await cosmiconfig(path.basename(pkg.name)).load(
        getGlobalConfigFilePath(),
    );
}
