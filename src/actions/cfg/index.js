import fs from "fs-extra";
import chalk from "chalk";
import { APW, CfgBuilder, Options } from "../../../lib";
import {
    getGlobalAppConfig,
    getLocalAppConfig,
    getFilesByExtension,
    printFiles,
} from "../../../lib/utils";

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

            const globalConfig = getGlobalAppConfig();

            for (const workspaceFile of workspaceFiles) {
                console.log(
                    chalk.blue(`Generating CFG for ${workspaceFile}...`),
                );

                const localConfig = getLocalAppConfig(workspaceFile);

                const apw = new APW(workspaceFile);

                const options = Options.getCfgOptions(
                    apw,
                    cliOptions,
                    localConfig.cfg,
                    globalConfig.cfg,
                );

                const cfgBuilder = new CfgBuilder(apw, options);
                const cfg = cfgBuilder.build();

                fs.writeFile(options.outputFile, cfg);
            }
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default cfg;
