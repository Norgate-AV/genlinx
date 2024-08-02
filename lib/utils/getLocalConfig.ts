import path from "node:path";
import chalk from "chalk";
import { cosmiconfig, CosmiconfigResult } from "cosmiconfig";
import {
    ArchiveCliArgs,
    BuildCliArgs,
    CfgCliArgs,
    CliOptions,
    LocalConfig,
} from "../@types/index.js";
import { getPackageJson } from "./index.js";

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

export async function getLocalConfig(
    options: CliOptions,
): Promise<CosmiconfigResult> {
    const packageJson = await getPackageJson();

    if (!packageJson) {
        throw new Error("package.json not found");
    }

    if (!packageJson.name) {
        throw new Error("package.json missing name field");
    }

    const explorer = cosmiconfig(path.basename(packageJson.name));
    const result = await explorer.search();

    if (!result) {
        return null;
    }

    const verbose =
        (options.cfg as CfgCliArgs)?.verbose ||
        (options.build as BuildCliArgs)?.verbose ||
        (options.archive as ArchiveCliArgs)?.verbose;

    verbose &&
        console.log(
            chalk.blue(`Using local config file at ${result.filepath}`),
        );

    const root = path.dirname(result.filepath);
    result.config = resolvePaths(root, result.config as LocalConfig);

    return result;
}
