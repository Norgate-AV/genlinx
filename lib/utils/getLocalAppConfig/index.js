import fs from "fs-extra";
import { getLocalAppConfigFilePath } from "../getLocalAppConfigFilePath";

const defaultLocalAppConfig = {
    cfg: {},
    archive: {},
    build: {},
};

export async function getLocalAppConfig(cwd) {
    let configFilePath = null;
    let localConfig = defaultLocalAppConfig;
    let result = null;

    try {
        configFilePath = await getLocalAppConfigFilePath(cwd);

        if (!configFilePath) {
            throw new Error("No local config file found.");
        }

        localConfig = JSON.parse(await fs.readFile(configFilePath, "utf8"));

        result = {
            path: configFilePath,
            config: localConfig,
        };
    } catch (error) {
        result = null;
    }

    return result;
}

export default getLocalAppConfig;
