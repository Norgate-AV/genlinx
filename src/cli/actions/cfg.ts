import fs from "node:fs/promises";
import chalk from "chalk";
import { AmxExtensions } from "../../core/apw/index.js";
import { CfgBuilder } from "../../core/cfg/index.js";
import {
    getFilesByExtension,
    getAppConfig,
    loadAPW,
    printFiles,
    selectFiles,
    shouldPromptUser,
} from "../../utils/index.js";
import {
    CfgOptions,
    CfgCliArgs,
    AmxFileType as FileType,
} from "../../@types/index.js";

export const cfg = {
    async create(args: CfgCliArgs): Promise<void> {
        try {
            const config = await getAppConfig({
                cfg: {
                    ...args,
                },
            });

            const { workspaceFiles, verbose } = config.cfg as CfgOptions;

            if (workspaceFiles.length === 0) {
                verbose &&
                    console.log(chalk.blue("Searching for workspace files..."));

                const locatedWorkspaceFiles = await getFilesByExtension(
                    AmxExtensions[FileType.Workspace],
                );

                if (locatedWorkspaceFiles.length) {
                    verbose && printFiles(locatedWorkspaceFiles);
                    workspaceFiles.push(...locatedWorkspaceFiles);
                }
            }

            if (workspaceFiles.length === 0) {
                console.log(chalk.red("No workspace files found."));
                process.exit();
            }

            if (shouldPromptUser(config.cfg as CfgOptions, workspaceFiles)) {
                const selectedWorkspaceFiles =
                    await selectFiles(workspaceFiles);

                workspaceFiles.splice(0, workspaceFiles.length);
                workspaceFiles.push(...selectedWorkspaceFiles);
            }

            for (const workspaceFile of workspaceFiles) {
                verbose &&
                    console.log(
                        chalk.blue(`Generating CFG for ${workspaceFile}...`),
                    );

                const apw = await loadAPW(workspaceFile);

                const builder = new CfgBuilder(apw, config.cfg as CfgOptions);

                const cfg = await builder.build();

                const file = `${apw.id}.${config.cfg.outputFileSuffix}`;
                fs.writeFile(file, cfg);
            }
        } catch (error: any) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default cfg;
