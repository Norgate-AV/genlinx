import chalk from "chalk";
import { APW, ArchiveBuilder } from "../../../lib";
import {
    getFilesByExtension,
    getAppConfig,
    loadAPW,
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

            if (workspaceFiles.length === 0) {
                console.log(chalk.blue("Searching for workspace files..."));

                const locatedWorkspaceFiles = await getFilesByExtension(
                    APW.fileExtensions[APW.fileType.workspace],
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

            if (shouldPromptUser(cliOptions, workspaceFiles)) {
                const selectedWorkspaceFiles =
                    await selectFiles(workspaceFiles);

                workspaceFiles.splice(0, workspaceFiles.length);
                workspaceFiles.push(...selectedWorkspaceFiles);
            }

            const config = await getAppConfig(cliOptions);

            for (const workspaceFile of workspaceFiles) {
                console.log(
                    chalk.blue(`Generating archive for ${workspaceFile}...`),
                );

                const apw = await loadAPW(workspaceFile);

                const builder = new ArchiveBuilder(apw, config.archive);
                builder.build();
            }
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default archive;
