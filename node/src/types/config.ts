import { ArchiveConfig } from "./archive.js";

export type Config = {
    archive: ArchiveConfig;
};

export type GlobalConfig = Partial<Config>;

export type LocalConfig = Partial<Config>;
