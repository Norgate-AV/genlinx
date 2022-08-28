import fs from "fs-extra";
import path from "path";
import pkg from "../../../package.json";
import { defaultAppConfig } from "../../../config";
import { getAppConfigDirectory } from "../getAppConfigDirectory";

export function installDefaultAppConfig() {
    const configDirectory = getAppConfigDirectory();
    const configFileName = "config.json";
    const configFilePath = path.join(`${configDirectory}`, `${configFileName}`);

    if (!fs.existsSync(configFilePath)) {
        console.log(
            `Installing default config file for ${path.basename(
                pkg.name,
            )} at: ${configFilePath}`,
        );

        fs.ensureDirSync(configDirectory);
        fs.writeFileSync(
            configFilePath,
            JSON.stringify(defaultAppConfig, null, 4),
        );
    }
}

export default installDefaultAppConfig;
