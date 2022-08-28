import fs from "fs-extra";
import path from "path";
import { getAppConfigDirectory } from "../getAppConfigDirectory";
import pkg from "../../../package.json";

export function getAppConfig() {
    const configDir = getAppConfigDirectory();
    const configPath = path.join(configDir, "config.json");

    if (!fs.existsSync(configPath)) {
        console.log(
            `Could not locate config file for ${path.basename(
                pkg.name,
            )} at: ${configPath}`,
        );

        return null;
    }

    return JSON.parse(fs.readFileSync(configPath, "utf8"));
}

export default getAppConfig;
