import fs from "fs-extra";
import path from "path";
import os from "os";
import chalk from "chalk";
import { mergician } from "mergician";
import { findUp } from "find-up";
import pkg from "../../package.json";
import defaultConfig from "../../config/default.js";
import {
    ArchiveCliArgs,
    BuildCliArgs,
    CfgCliArgs,
    CliOptions,
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

async function getLocalAppConfig(options: CliOptions): Promise<LocalConfig> {
    const filePath = await findUp([".genlinxrc.json", ".genlinxrc"], {
        cwd: process.cwd(),
    });

    if (!filePath) {
        return {};
    }

    const verbose =
        (options.cfg as CfgCliArgs)?.verbose ||
        (options.build as BuildCliArgs)?.verbose ||
        (options.archive as ArchiveCliArgs)?.verbose;

    verbose &&
        console.log(chalk.blue(`Using local config file at ${filePath}`));

    const root = path.dirname(filePath);
    const config = resolvePaths(
        root,
        JSON.parse(await fs.readFile(filePath, "utf8")),
    );

    return config || {};
}

export async function getAppConfig(cliOptions: CliOptions): Promise<Config> {
    const global = await getGlobalAppConfig();
    const local = await getLocalAppConfig(cliOptions);

    const config = {
        default: defaultConfig,
        global: global || {},
        local: local || {},
        cli: cliOptions || {},
    };

    return mergician({
        appendArrays: true,
        dedupArrays: true,
        sortArrays: true,
    })(config.default, config.global, config.local, config.cli);
}

export default getAppConfig;
