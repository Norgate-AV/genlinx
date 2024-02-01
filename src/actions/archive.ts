import chalk from "chalk";
import { ArchiveBuilder, AmxExtensions } from "../../lib/index.js";
import {
    getFilesByExtension,
    getAppConfig,
    loadAPW,
    printFiles,
    selectFiles,
    shouldPromptUser,
} from "../../lib/utils/index.js";
import {
    ArchiveOptions,
    ArchiveCliArgs,
    AmxFileType as FileType,
} from "../../lib/@types/index.js";

export const archive = {
    async create(args: ArchiveCliArgs): Promise<void> {
        try {
            const config = await getAppConfig({
                archive: {
                    ...args,
                },
            });

            const { workspaceFiles, verbose } =
                config.archive as ArchiveOptions;

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

            if (
                shouldPromptUser(
                    config.archive as ArchiveOptions,
                    workspaceFiles,
                )
            ) {
                const selectedWorkspaceFiles =
                    await selectFiles(workspaceFiles);

                workspaceFiles.splice(0, workspaceFiles.length);
                workspaceFiles.push(...selectedWorkspaceFiles);
            }

            for (const workspaceFile of workspaceFiles) {
                verbose &&
                    console.log(
                        chalk.blue(
                            `Generating archive for ${workspaceFile}...`,
                        ),
                    );

                const apw = await loadAPW(workspaceFile);

                const builder = new ArchiveBuilder(
                    apw,
                    config.archive as ArchiveOptions,
                );
                builder.build();
            }
        } catch (error: any) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default archive;
