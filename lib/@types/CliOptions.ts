import { CfgCliArgs, ArchiveCliArgs, BuildCliArgs } from "./index.js";

export type CliOptions = {
    cfg?: CfgCliArgs;
    archive?: ArchiveCliArgs;
    build?: BuildCliArgs;
};
