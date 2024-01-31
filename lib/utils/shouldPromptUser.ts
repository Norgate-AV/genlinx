import { ArchiveOptions, BuildOptions, CfgOptions } from "../@types/index.js";

export function shouldPromptUser(
    options: ArchiveOptions | BuildOptions | CfgOptions,
    files: Array<string>,
): boolean {
    return !options.all && files.length > 1;
}
