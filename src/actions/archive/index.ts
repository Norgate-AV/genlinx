import chalk from "chalk";
import { ArchiveBuilder } from "../../../lib";
import { extensions } from "../../../lib/APW";
import {
    getFilesByExtension,
    getAppConfig,
    loadAPW,
    printFiles,
    selectFiles,
} from "../../../lib/utils";
import {
    ArchiveConfig,
    CliArchiveOptions,
    FileType,
} from "../../../lib/@types";

function shouldPromptUser(options: ArchiveConfig, files: Array<string>) {
    return !options.all && files.length > 1;
}

export const archive = {
    async create(cliOptions: CliArchiveOptions): Promise<void> {
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

            if (shouldPromptUser(config.archive, workspaceFiles)) {
                const selectedWorkspaceFiles =
                    await selectFiles(workspaceFiles);

                workspaceFiles.splice(0, workspaceFiles.length);
                workspaceFiles.push(...selectedWorkspaceFiles);
            }

            for (const workspaceFile of workspaceFiles) {
                console.log(
                    chalk.blue(`Generating archive for ${workspaceFile}...`),
                );

                const apw = await loadAPW(workspaceFile);

                const builder = new ArchiveBuilder(apw, config.archive);
                builder.build();
            }
        } catch (error: any) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default archive;