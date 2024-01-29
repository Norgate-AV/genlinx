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
