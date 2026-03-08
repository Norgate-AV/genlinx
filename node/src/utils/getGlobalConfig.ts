import path from "node:path";
import os from "node:os";
import { cosmiconfig, CosmiconfigResult } from "cosmiconfig";
import { getPackageJson } from "./index.js";

export function getGlobalConfigFilePath(name: string): string {
    const file = "config.json";
    const directory =
        process.env.GENLINX_CONFIG_DIR ||
        path.join(os.homedir(), ".config", path.basename(name));

    return path.join(directory, file);
}

export async function getGlobalConfig(): Promise<CosmiconfigResult | null> {
    const packageJson = await getPackageJson();

    if (!packageJson) {
        throw new Error("package.json not found");
    }

    if (!packageJson.name) {
        throw new Error("package.json missing name field");
    }

    try {
        const result = await cosmiconfig(path.basename(packageJson.name)).load(
            getGlobalConfigFilePath(packageJson.name),
        );

        return result;
    } catch (error: any) {
        return null;
    }
}
