import chalk from "chalk";
import { APW, ArchiveBuilder, Options } from "../../../lib";
import {
    getGlobalAppConfig,
    getLocalAppConfig,
    getFilesByExtension,
    printFiles,
    selectFiles,
} from "../../../lib/utils";

function shouldPromptUser(options, files) {
    return !options.all && files.length > 1;
}

export const archive = {
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
                    chalk.blue(`Generating archive for ${workspaceFile}...`),
                );

                const localConfig = getLocalAppConfig(workspaceFile);

                const apw = new APW(workspaceFile);

                const options = Options.getArchiveOptions(
                    apw,
                    cliOptions,
                    localConfig.archive,
                    globalConfig.archive,
                );

                console.log(`options: ${JSON.stringify(options, null, 4)}`);

                const builder = new ArchiveBuilder(apw, options);
                builder.build();
            }
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default archive;
