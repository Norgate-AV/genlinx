import fs from "fs-extra";
import path from "path";
import { getLocalAppConfigFilePath } from "../getLocalAppConfigFilePath";

export function getLocalAppConfig(filePath) {
    const configFilePath = getLocalAppConfigFilePath(path.dirname(filePath));

    if (!fs.existsSync(configFilePath)) {
        return {};
    }

    return JSON.parse(fs.readFileSync(configFilePath, "utf8"));
}

export default getLocalAppConfig;
