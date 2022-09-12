import fs from "fs-extra";
import chalk from "chalk";
import { APW, CfgBuilder } from "../../../lib";
import {
    getFilesByExtension,
    getGlobalAppConfig,
    getLocalAppConfig,
    getOptions,
    printFiles,
    selectFiles,
} from "../../../lib/utils";

function shouldPromptUser(options, files) {
    return !options.all && files.length > 1;
}

export const cfg = {
    async create(cliOptions) {
        try {
            const { workspaceFiles } = cliOptions;

            if (!workspaceFiles.length) {
                console.log(chalk.blue("Searching for workspace files..."));
                workspaceFiles.push(
                    ...(await getFilesByExtension(
                        APW.fileExtensions[APW.fileType.workspace],
                    )),
                );

                if (workspaceFiles.length) {
                    printFiles(workspaceFiles);
                }
            }

            if (!workspaceFiles.length) {
                console.log(chalk.red("No workspace files found."));
                process.exit();
            }

            if (shouldPromptUser(cliOptions, workspaceFiles)) {
                const selectedWorkspaceFiles = await selectFiles(
                    workspaceFiles,
                );

                workspaceFiles.splice(0, workspaceFiles.length);
                workspaceFiles.push(...selectedWorkspaceFiles);
            }

            const globalConfig = getGlobalAppConfig();

            for (const workspaceFile of workspaceFiles) {
                console.log(
                    chalk.blue(`Generating CFG for ${workspaceFile}...`),
                );

                const localConfig = getLocalAppConfig(workspaceFile);

                const apw = new APW(workspaceFile);

                const options = getOptions(
                    cliOptions,
                    localConfig.cfg,
                    globalConfig.cfg,
                );

                const cfgBuilder = new CfgBuilder(apw, options);
                const cfg = cfgBuilder.build();

                const outputFile = `${apw.id}.${options.outputFileSuffix}`;
                fs.writeFile(outputFile, cfg);
            }
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default cfg;
