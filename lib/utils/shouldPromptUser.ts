import { ArchiveConfig, BuildConfig, CfgConfig } from "../@types/index.js";

export function shouldPromptUser(
    options: ArchiveConfig | BuildConfig | CfgConfig,
    files: Array<string>,
) {
    return !options.all && files.length > 1;
}
