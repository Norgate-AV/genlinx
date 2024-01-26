import { CliArchiveOptions } from "./CliArchiveOptions";
import { CliBuildOptions } from "./CliBuildOptions";
import { CliCfgOptions } from "./CliCfgOptions";

export type CliConfig = {
    cfg?: CliCfgOptions;
    archive?: CliArchiveOptions;
    build?: CliBuildOptions;
};
