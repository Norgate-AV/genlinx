import path from "node:path";
import { mergician } from "mergician";
import defaultConfig from "../config/default.js";
import { getGlobalConfig, getLocalConfig } from "./index.js";
import { CliOptions, Config } from "../@types/index.js";

function normalizePaths(config: Config): Config {
    config.cfg.includePath = config.cfg.includePath.map((p) =>
        path.normalize(p),
    );

    config.cfg.modulePath = config.cfg.modulePath.map((p) => path.normalize(p));
    config.cfg.libraryPath = config.cfg.libraryPath.map((p) =>
        path.normalize(p),
    );

    config.archive.extraFileSearchLocations =
        config.archive.extraFileSearchLocations.map((p) => path.normalize(p));

    config.build.nlrc.path = path.normalize(config.build.nlrc.path);
    config.build.nlrc.includePath = config.build.nlrc.includePath.map((p) =>
        path.normalize(p),
    );

    config.build.nlrc.modulePath = config.build.nlrc.modulePath.map((p) =>
        path.normalize(p),
    );

    config.build.nlrc.libraryPath = config.build.nlrc.libraryPath.map((p) =>
        path.normalize(p),
    );

    config.build.shell.path = path.normalize(config.build.shell.path);

    return config;
}

export async function getAppConfig(cliOptions: CliOptions): Promise<Config> {
    const global = await getGlobalConfig();
    const local = await getLocalConfig(cliOptions);

    const config = {
        default: defaultConfig,
        global: global ? global.config : {},
        local: local ? local.config : {},
        cli: cliOptions || {},
    };

    const merged: Config = mergician({
        prependArrays: true,
        dedupArrays: true,
        sortArrays: false,
    })(config.default, config.global, config.local, config.cli);

    return normalizePaths(merged);
}

export default getAppConfig;
