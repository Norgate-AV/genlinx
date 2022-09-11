import fs from "fs-extra";
import path from "path";
import chalk from "chalk";
import { APW, CfgBuilder, Options } from "../../../lib";
import { getGlobalAppConfig, getLocalAppConfig } from "../../../lib/utils";

async function getWorkspaceFiles() {
    const entities = await fs.readdir(process.cwd());
    const workspaceFiles = entities.filter(
        (entity) =>
            path.extname(entity) === APW.fileExtensions[APW.fileType.workspace],
    );

    return workspaceFiles;
}

export const cfg = {
    async create(cliOptions) {
        try {
            const { workspaceFiles } = cliOptions;

            if (!workspaceFiles.length) {
                console.log(chalk.blue("Searching for workspace files..."));
                workspaceFiles.push(...(await getWorkspaceFiles()));
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
