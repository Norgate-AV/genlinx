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

    try {
        configFilePath = await getLocalAppConfigFilePath(cwd);

        if (!configFilePath) {
            throw new Error("No local config file found.");
        }

        localConfig = JSON.parse(await fs.readFile(configFilePath, "utf8"));

        // if (!Object.prototype.hasOwnProperty.call(localConfig, "cfg")) {
        //     localConfig.cfg = defaultLocalAppConfig.cfg;
        // }

        // if (!Object.prototype.hasOwnProperty.call(localConfig, "archive")) {
        //     localConfig.archive = defaultLocalAppConfig.archive;
        // }

        // if (!Object.prototype.hasOwnProperty.call(localConfig, "build")) {
        //     localConfig.build = defaultLocalAppConfig.build;
        // }

        // return {
        //     path: configFilePath,
        //     config: localConfig,
        // };
    } catch (error) {
        // console.log(error);
    }

    return {
        path: configFilePath,
        config: localConfig,
    };
}

export default getLocalAppConfig;
