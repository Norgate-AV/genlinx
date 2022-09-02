import path from "path";
import chalk from "chalk";
import AdmZip from "adm-zip";
import { APW } from "../APW";
import { walkDirectory } from "../utils";

export class ArchiveBuilder {
    constructor(apw, options) {
        this.apw = apw;
        this.options = options;
        this.builder = new AdmZip();
    }

    #addWorkspaceFiles(builder, apw, options) {
        for (const file of apw.allFiles) {
            try {
                builder.addLocalFile(file.path, path.dirname(file.path));
                if (options.includeCompiledSourceFiles) {
                    if (file.type === APW.fileType.masterSrc) {
                        const compiledFilePath = file.path.replace(
                            APW.fileExtensions[APW.fileType.masterSrc],
                            APW.compiledFileExtensions[APW.fileType.masterSrc],
                        );

                        builder.addLocalFile(
                            compiledFilePath,
                            path.dirname(compiledFilePath),
                        );
                    }
                }

                if (options.includeCompiledModuleFiles) {
                    if (file.type === APW.fileType.module) {
                        const compiledFilePath = file.path.replace(
                            APW.fileExtensions[APW.fileType.module],
                            APW.compiledFileExtensions[APW.fileType.module],
                        );

                        builder.addLocalFile(
                            compiledFilePath,
                            path.dirname(compiledFilePath),
                        );
                    }
                }
            } catch (error) {
                console.log(chalk.red(error.message));
                continue;
            }

            console.log(chalk.green(`Added file: ${file.path}`));
        }

        return this;
    }

    async build() {
        const { builder, apw, options } = this;

        console.log(chalk.blue("Creating archive..."));

        this.#addWorkspaceFiles(builder, apw, options);

        if (options.includeFilesNotInWorkspace) {
            console.log(
                chalk.blue(
                    "Searching for extra files that are not part of the workspace...",
                ),
            );

            const extraFileReferences = apw.getExtraFileReferences();

            if (extraFileReferences) {
                console.log(
                    chalk.green(
                        `Found ${extraFileReferences.length} extra files referenced...`,
                    ),
                );

                for (const fileReference of extraFileReferences) {
                    console.log(chalk.cyan(`--> ${fileReference}}`));
                }

                console.log(
                    chalk.blue("Searching known locations for extra files..."),
                );

                const files = [];
                for (const location of options.extraFileLocations) {
                    console.log(chalk.cyan(`--> ${location}`));
                    files.push(...walkDirectory(location));
                }

                for (const fileReference of extraFileReferences) {
                    const filePath = files.find((filePath) =>
                        filePath.includes(fileReference),
                    );

                    if (!filePath) {
                        console.log(
                            chalk.red(`Could not find ${fileReference}`),
                        );

                        continue;
                    }

                    console.log(
                        chalk.green(`Found ${fileReference}: ${filePath}`),
                    );

                    const fileType = APW.getFileType();
                    console.log(fileType);

                    switch (path.extname(filePath)) {
                        case APW.fileExtensions[APW.fileType.module]:
                            builder.addLocalFile(filePath, ".genlinx");

                            console.log(chalk.green(`Added file: ${filePath}`));

                            if (!options.includeCompiledModuleFiles) {
                                break;
                            }

                            const compiledFilePath = filePath.replace(
                                APW.fileExtensions[APW.fileType.module],
                                APW.compiledFileExtensions[APW.fileType.module],
                            );

                            if (!compiledFilePath) {
                                break;
                            }

                            builder.addLocalFile(compiledFilePath, ".genlinx");

                            console.log(
                                chalk.green(`Added file: ${compiledFilePath}`),
                            );

                            break;
                        case APW.fileExtensions[APW.fileType.duet]:
                        case APW.fileExtensions[APW.fileType.include]:
                            builder.addLocalFile(filePath, ".genlinx");
                            console.log(chalk.green(`Added file: ${filePath}`));
                            break;
                        default:
                    }
                }
            }
        }

        console.log(chalk.blue("Compressing files into a zip archive..."));
        builder.writeZip(options.outputFile);

        console.log(chalk.green.bold(`Created archive: ${options.outputFile}`));

        for (const entry of builder.getEntries()) {
            console.log(chalk.cyan(`--> ${entry.entryName}`));
        }
    }
}

export default ArchiveBuilder;
