import fs from "fs-extra";
import chalk from "chalk";
import { CfgBuilder } from "../../../lib";
import { extensions } from "../../../lib/APW";
import {
    getFilesByExtension,
    getAppConfig,
    loadAPW,
    printFiles,
    selectFiles,
} from "../../../lib/utils";
import { CfgConfig, CliCfgOptions, FileType } from "../../../lib/@types";

function shouldPromptUser(options: CfgConfig, files: Array<string>): boolean {
    return !options.all && files.length > 1;
}

export const cfg = {
    async create(cliOptions: CliCfgOptions): Promise<void> {
        try {
            const { workspaceFiles } = cliOptions;

            const config = await getAppConfig(cliOptions);

            if (workspaceFiles.length === 0) {
                console.log(chalk.blue("Searching for workspace files..."));

                const locatedWorkspaceFiles = await getFilesByExtension(
                    extensions[FileType.Workspace],
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
