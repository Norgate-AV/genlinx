import { ArchiveCliArgs } from "./archive.js";
import { BuildCliArgs } from "./build.js";

export type CliArgs = ArchiveCliArgs | BuildCliArgs;

export type CliOptions = {
    archive?: ArchiveCliArgs;
    build?: BuildCliArgs;
};
