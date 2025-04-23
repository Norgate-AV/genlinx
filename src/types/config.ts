import { ArchiveConfig } from "./archive.js";
import { BuildConfig } from "./build.js";
import { CfgConfig } from "./cfg.js";

export type Config = {
    cfg: CfgConfig;
    archive: ArchiveConfig;
    build: BuildConfig;
};

export type ConfigCliArgs = {
    list?: boolean;
    edit?: boolean;
    local?: boolean;
    global?: boolean;
};

export type GlobalConfig = Partial<Config>;

export type LocalConfig = Partial<Config>;
