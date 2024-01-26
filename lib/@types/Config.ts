import { CliArchiveOptions, CliBuildOptions, CliCfgOptions } from ".";
import { ArchiveConfig } from "./ArchiveConfig";
import { BuildConfig } from "./BuildConfig";
import { CfgConfig } from "./CfgConfig";

export type Config = {
    cfg: CfgConfig & CliCfgOptions;
    archive: ArchiveConfig & CliArchiveOptions;
    build: BuildConfig & CliBuildOptions;
};

export type GlobalConfig = Partial<Config>;
export type LocalConfig = Partial<Config>;
