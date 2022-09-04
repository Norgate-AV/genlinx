import fs from "fs-extra";
import path from "path";
import { defaultGlobalAppConfig } from "../../../config";
import { getGlobalAppConfigFilePath } from "../getGlobalAppConfigFilePath";
import { installDefaultGlobalAppConfig } from "../installDefaultGlobalAppConfig";
import pkg from "../../../package.json";

export function getGlobalAppConfig() {
    const configFilePath = getGlobalAppConfigFilePath();

    if (!fs.existsSync(configFilePath)) {
        console.log(
            `Could not locate config file for ${path.basename(
                pkg.name,
            )} at: ${configFilePath}`,
        );

        installDefaultGlobalAppConfig();
        return defaultGlobalAppConfig;
    }

    return JSON.parse(fs.readFileSync(configFilePath, "utf8"));
}

export default getGlobalAppConfig;
