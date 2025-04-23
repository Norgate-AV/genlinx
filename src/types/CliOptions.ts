import {
    CfgCliArgs,
    ArchiveCliArgs,
    BuildCliArgs,
    ConfigCliArgs,
} from "./index.js";

export type CliOptions = {
    cfg?: CfgCliArgs;
    archive?: ArchiveCliArgs;
    build?: BuildCliArgs;
    config?: ConfigCliArgs;
};
