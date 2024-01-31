import { ArchiveConfig, BuildConfig, CfgConfig } from "./index.js";

export type Config = {
    cfg: CfgConfig;
    archive: ArchiveConfig;
    build: BuildConfig;
};
