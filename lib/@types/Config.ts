import {
    ArchiveConfig,
    BuildConfig,
    CfgConfig,
    CliArchiveOptions,
    CliBuildOptions,
    CliCfgOptions,
} from "./index.js";

export type Config = {
    cfg: CfgConfig & CliCfgOptions;
    archive: ArchiveConfig & CliArchiveOptions;
    build: BuildConfig & CliBuildOptions;
};

export type GlobalConfig = Partial<Config>;
export type LocalConfig = Partial<Config>;
