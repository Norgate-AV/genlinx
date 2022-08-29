import fs from "fs-extra";
import path from "path";
import pkg from "../../../package.json";
import { defaultAppConfig } from "../../../config";
import { getAppConfigFilePath } from "../getAppConfigFilePath";

export function installDefaultAppConfig() {
    const configFilePath = getAppConfigFilePath();

    if (!fs.existsSync(configFilePath)) {
        console.log(
            `Installing default config file for ${path.basename(
                pkg.name,
            )} at: ${configFilePath}`,
        );

        fs.ensureDirSync(path.dirname(configFilePath));
        fs.writeFileSync(
            configFilePath,
            JSON.stringify(defaultAppConfig, null, 4),
        );
    }
}

export default installDefaultAppConfig;
