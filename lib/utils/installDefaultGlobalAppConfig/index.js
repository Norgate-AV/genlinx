import fs from "fs-extra";
import path from "path";
import pkg from "../../../package.json";
import { defaultGlobalAppConfig } from "../../../config";
import { getGlobalAppConfigFilePath } from "../getGlobalAppConfigFilePath";

export function installDefaultGlobalAppConfig() {
    const configFilePath = getGlobalAppConfigFilePath();

    if (!fs.existsSync(configFilePath)) {
        console.log(
            `Installing default global config file for ${path.basename(
                pkg.name,
            )} at: ${configFilePath}`,
        );

        fs.ensureDirSync(path.dirname(configFilePath));
        fs.writeFileSync(
            configFilePath,
            JSON.stringify(defaultGlobalAppConfig, null, 4),
        );
    }
}

export default installDefaultGlobalAppConfig;
