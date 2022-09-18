import fs from "fs-extra";
import path from "path";
import { defaultGlobalAppConfig } from "../../../config";
import { getGlobalAppConfigFilePath } from "../getGlobalAppConfigFilePath";
import { installDefaultGlobalAppConfig } from "../installDefaultGlobalAppConfig";
import pkg from "../../../package.json";

export async function getGlobalAppConfig() {
    const configFilePath = getGlobalAppConfigFilePath();

    const exists = await fs.pathExists(configFilePath);

    if (!exists) {
        console.log(
            `Could not locate config file for ${path.basename(
                pkg.name,
            )} at: ${configFilePath}`,
        );

        await installDefaultGlobalAppConfig();
        return defaultGlobalAppConfig;
    }

    return JSON.parse(await fs.readFile(configFilePath, "utf8"));
}

export default getGlobalAppConfig;
