import { CliCfgOptions, CliArchiveOptions, CliBuildOptions } from "./index.js";

export type CliConfig = {
    cfg?: CliCfgOptions;
    archive?: CliArchiveOptions;
    build?: CliBuildOptions;
};
