import fs from "fs-extra";
import path from "path";
import pkg from "../../../package.json";
import { defaultGlobalAppConfig } from "../../../config";
import { getGlobalAppConfigFilePath } from "../getGlobalAppConfigFilePath";

export async function installDefaultGlobalAppConfig() {
    const configFilePath = getGlobalAppConfigFilePath();

    const exists = await fs.pathExists(configFilePath);

    if (!exists) {
        console.log(
            `Installing default global config file for ${path.basename(
                pkg.name,
            )} at: ${configFilePath}`,
        );

        await fs.ensureDir(path.dirname(configFilePath));
        fs.writeFile(
            configFilePath,
            JSON.stringify(defaultGlobalAppConfig, null, 4),
        );
    }
}

export default installDefaultGlobalAppConfig;
