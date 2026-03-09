import { ArchiveConfig } from "./archive.js";
import { BuildConfig } from "./build.js";

export type Config = {
    archive: ArchiveConfig;
    build: BuildConfig;
};

export type GlobalConfig = Partial<Config>;

export type LocalConfig = Partial<Config>;
