import fs from "fs-extra";
import path from "path";
import os from "os";
import util from "util";
import { mergician } from "mergician";
import { findUp } from "find-up";
import pkg from "../../package.json";
import defaultConfig from "../../config/index.js";
import {
    CliArchiveOptions,
    CliBuildOptions,
    CliCfgOptions,
    Config,
    GlobalConfig,
    LocalConfig,
} from "../@types/index.js";

function getGlobalAppConfigFilePath(): string {
    const file = "config.json";
    const directory =
        process.env.GENLINX_CONFIG_DIR ||
        path.join(os.homedir(), ".config", `${path.basename(pkg.name)}`);

    return path.join(`${directory}`, `${file}`);
}

async function getGlobalAppConfig(): Promise<GlobalConfig> {
    const filePath = getGlobalAppConfigFilePath();

    const exists = await fs.pathExists(filePath);

    if (!exists) {
        return {};
    }

    const config = JSON.parse(await fs.readFile(filePath, "utf8"));

    return config || {};
}

function resolvePaths(root: string, config: LocalConfig): LocalConfig {
    if (config.cfg?.includePath) {
        config.cfg.includePath = config.cfg.includePath.map((p) =>
            path.resolve(root, p),
        );
    }

    if (config.cfg?.modulePath) {
        config.cfg.modulePath = config.cfg.modulePath.map((p) =>
            path.resolve(root, p),
        );
    }

    if (config.cfg?.libraryPath) {
        config.cfg.libraryPath = config.cfg.libraryPath.map((p) =>
            path.resolve(root, p),
        );
    }

    if (config.archive?.extraFileSearchLocations) {
        config.archive.extraFileSearchLocations =
            config.archive.extraFileSearchLocations.map((p) =>
                path.resolve(root, p),
            );
    }

    if (config.build?.nlrc?.includePath) {
        config.build.nlrc.includePath = config.build.nlrc.includePath.map((p) =>
            path.resolve(root, p),
        );
    }

    if (config.build?.nlrc?.modulePath) {
        config.build.nlrc.modulePath = config.build.nlrc.modulePath.map((p) =>
            path.resolve(root, p),
        );
    }

    if (config.build?.nlrc?.libraryPath) {
        config.build.nlrc.libraryPath = config.build.nlrc.libraryPath.map((p) =>
            path.resolve(root, p),
        );
    }

    return config;
}

async function getLocalAppConfig(): Promise<LocalConfig> {
    const filePath = await findUp([".genlinxrc.json", ".genlinxrc"], {
        cwd: process.cwd(),
    });

    if (!filePath) {
        return {};
    }

    console.log(`Found local config file at ${filePath}`);

    const root = path.dirname(filePath);
    const config = resolvePaths(
        root,
        JSON.parse(await fs.readFile(filePath, "utf8")),
    );

    return config || {};
}

export async function getAppConfig(
    cliOptions: CliCfgOptions | CliArchiveOptions | CliBuildOptions,
): Promise<Config> {
    // console.log("GLOBAL-----------------------------------------------");
    const global = await getGlobalAppConfig();
    // console.log(
    //     util.inspect(global, { showHidden: false, depth: null, colors: true }),
    // );

    // console.log("LOCAL-----------------------------------------------");
    const local = await getLocalAppConfig();
    // console.log(
    //     util.inspect(local, { showHidden: false, depth: null, colors: true }),
    // );

    const config = {
        default: defaultConfig,
        global: global || {},
        local: local || {},
        cli: cliOptions || {},
    };

    // console.log("FULL CONFIG-----------------------------------------");
    // console.log(
    //     util.inspect(config, { showHidden: false, depth: null, colors: true }),
    // );

    return mergician({
        appendArrays: true,
        dedupArrays: true,
        sortArrays: true,
    })(config.default, config.global, config.local, config.cli);
}

export default getAppConfig;
