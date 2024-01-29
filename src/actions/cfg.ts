import fs from "fs-extra";
import chalk from "chalk";
import { CfgBuilder, AmxExtensions } from "../../lib/index.js";
import {
    getFilesByExtension,
    getAppConfig,
    loadAPW,
    printFiles,
    selectFiles,
    shouldPromptUser,
} from "../../lib/utils/index.js";
import {
    CliCfgOptions,
    AmxFileType as FileType,
} from "../../lib/@types/index.js";

export const cfg = {
    async create(cliOptions: CliCfgOptions): Promise<void> {
        try {
            const config = await getAppConfig({
                cfg: {
                    ...cliOptions,
                },
            });

            const { workspaceFiles } = config.cfg;

            if (workspaceFiles.length === 0) {
                console.log(chalk.blue("Searching for workspace files..."));

                const locatedWorkspaceFiles = await getFilesByExtension(
                    AmxExtensions[FileType.Workspace],
                );

                if (locatedWorkspaceFiles.length) {
                    printFiles(locatedWorkspaceFiles);
                    workspaceFiles.push(...locatedWorkspaceFiles);
                }
            }

            if (workspaceFiles.length === 0) {
                console.log(chalk.red("No workspace files found."));
                process.exit();
            }

            if (shouldPromptUser(config.cfg, workspaceFiles)) {
                const selectedWorkspaceFiles =
                    await selectFiles(workspaceFiles);

                workspaceFiles.splice(0, workspaceFiles.length);
                workspaceFiles.push(...selectedWorkspaceFiles);
            }

            for (const workspaceFile of workspaceFiles) {
                console.log(
                    chalk.blue(`Generating CFG for ${workspaceFile}...`),
                );

                const apw = await loadAPW(workspaceFile);

                const cfgBuilder = new CfgBuilder(apw, config.cfg);
                const cfg = cfgBuilder.build();

                const outputFile = `${apw.id}.${config.cfg.outputFileSuffix}`;
                fs.writeFile(outputFile, cfg);
            }
        } catch (error: any) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default cfg;
