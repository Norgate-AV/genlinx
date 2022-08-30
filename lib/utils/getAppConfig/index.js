import fs from "fs-extra";
import path from "path";
import { defaultAppConfig } from "../../../config";
import { getAppConfigFilePath } from "../getAppConfigFilePath";
import { installDefaultAppConfig } from "../installDefaultAppConfig";
import pkg from "../../../package.json";

export function getAppConfig() {
    const configFilePath = getAppConfigFilePath();

    if (!fs.existsSync(configFilePath)) {
        console.log(
            `Could not locate config file for ${path.basename(
                pkg.name,
            )} at: ${configFilePath}`,
        );

        installDefaultAppConfig();
        return defaultAppConfig;
    }

    return JSON.parse(fs.readFileSync(configFilePath, "utf8"));
}

export default getAppConfig;
