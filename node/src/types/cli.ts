import { CfgCliArgs } from "./cfg.js";
import { ArchiveCliArgs } from "./archive.js";
import { BuildCliArgs } from "./build.js";
import { ConfigCliArgs } from "./config.js";

export type CliArgs = CfgCliArgs | ArchiveCliArgs | BuildCliArgs;

export type CliOptions = {
    cfg?: CfgCliArgs;
    archive?: ArchiveCliArgs;
    build?: BuildCliArgs;
    config?: ConfigCliArgs;
};

export type FindCliArgs = {
    timeout: number;
    json?: boolean;
};
