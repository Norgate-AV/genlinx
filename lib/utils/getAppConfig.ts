import { mergician } from "mergician";
import defaultConfig from "../../config/default.js";
import { getGlobalConfig, getLocalConfig } from "./index.js";
import { CliOptions, Config } from "../@types/index.js";

export async function getAppConfig(cliOptions: CliOptions): Promise<Config> {
    const global = await getGlobalConfig();
    const local = await getLocalConfig(cliOptions);

    const config = {
        default: defaultConfig,
        global: global ? global.config : {},
        local: local ? local.config : {},
        cli: cliOptions || {},
    };

    return mergician({
        prependArrays: true,
        dedupArrays: true,
        sortArrays: false,
    })(config.default, config.global, config.local, config.cli);
}

export default getAppConfig;
