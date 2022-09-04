import fs from "fs-extra";
import path from "path";
import { getLocalAppConfigFilePath } from "../getLocalAppConfigFilePath";

const defaultLocalAppConfig = {
    cfg: {},
    archive: {},
    build: {},
};

export function getLocalAppConfig(filePath) {
    const configFilePath = getLocalAppConfigFilePath(path.dirname(filePath));

    if (!fs.existsSync(configFilePath)) {
        return defaultLocalAppConfig;
    }

    const localConfig = JSON.parse(fs.readFileSync(configFilePath, "utf8"));

    if (!Object.prototype.hasOwnProperty.call(localConfig, "cfg")) {
        localConfig.cfg = defaultLocalAppConfig.cfg;
    }

    if (!Object.prototype.hasOwnProperty.call(localConfig, "archive")) {
        localConfig.archive = defaultLocalAppConfig.archive;
    }

    if (!Object.prototype.hasOwnProperty.call(localConfig, "build")) {
        localConfig.build = defaultLocalAppConfig.build;
    }

    return localConfig;
}

export default getLocalAppConfig;
