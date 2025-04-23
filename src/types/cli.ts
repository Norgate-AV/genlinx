import {
    CfgCliArgs,
    ArchiveCliArgs,
    BuildCliArgs,
    ConfigCliArgs,
} from "./index.js";

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
