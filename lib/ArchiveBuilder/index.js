import path from "path";
import chalk from "chalk";
import AdmZip from "adm-zip";
import { APW } from "../APW";
import { walkDirectory } from "../utils";

export class ArchiveBuilder {
    #extraFilesDirectory = ".genlinx";

    constructor(apw, options) {
        this.apw = apw;
        this.options = options;
        this.builder = new AdmZip();
    }

    #addFile(builder, filePath, archiveFilePath = undefined) {
        builder.addLocalFile(filePath, archiveFilePath);

        console.log(chalk.green(`Added file: ${filePath}`));

        return this;
    }

    static #getExtraFilesOnDisk(locations) {
        console.log(chalk.blue("Searching known locations for extra files..."));

        const files = [];

        for (const location of locations) {
            console.log(chalk.cyan(`--> Searching ${location}`));
            files.push(...walkDirectory(location));
        }

        return files;
    }

    #listExtraFileReferences(references) {
        console.log(
            chalk.green(`Found ${references.length} extra files referenced...`),
        );

        for (const reference of references) {
            console.log(chalk.cyan(`--> ${reference}`));
        }

        return this;
    }

    static #getLocatedFileReferences(files, references) {
        const locatedFileReferences = files.filter((file) =>
            references.includes(path.basename(file, path.extname(file))),
        );

        console.log(
            chalk.green(
                `Found ${locatedFileReferences.length} of ${references.length} extra files on disk...`,
            ),
        );

        return locatedFileReferences;
    }

    #addCompiledFile(builder, fileType, filePath, references) {
        const compiledFile = filePath.replace(
            APW.fileExtensions[fileType],
            APW.compiledFileExtensions[fileType],
        );

        const compiledFilePath = references.find((filePath) =>
            filePath.includes(compiledFile),
        );

        if (!compiledFilePath) {
            console.log(
                chalk.red(`Could not find ${path.basename(compiledFile)}`),
            );

            return this;
        }

        console.log(
            chalk.cyan(
                `--> Found ${path.basename(compiledFile)}: ${compiledFilePath}`,
            ),
        );

        this.#addFile(builder, compiledFilePath, this.#extraFilesDirectory);
        return this;
    }

    #addExtraFiles(builder, apw, options) {
        if (!options.includeFilesNotInWorkspace) {
            return this;
        }

        console.log(
            chalk.blue(
                "Searching for extra files that are not part of the workspace...",
            ),
        );

        const extraFileReferences = apw.getExtraFileReferences();

        if (!extraFileReferences.length) {
            return this;
        }

        this.#listExtraFileReferences(extraFileReferences);

        const { extraFileLocations } = options;
        const extraFiles =
            ArchiveBuilder.#getExtraFilesOnDisk(extraFileLocations);

        const locatedFileReferences = ArchiveBuilder.#getLocatedFileReferences(
            extraFiles,
            extraFileReferences,
        );

        for (const fileReference of extraFileReferences) {
            const filePath = locatedFileReferences.find((filePath) =>
                filePath.includes(fileReference),
            );

            if (!filePath) {
                console.log(chalk.red(`Could not find ${fileReference}`));
                continue;
            }

            console.log(chalk.green(`Found ${fileReference}: ${filePath}`));

            switch (path.extname(filePath)) {
                case APW.fileExtensions[APW.fileType.module]:
                    this.#addFile(builder, filePath, this.#extraFilesDirectory);

                    if (!options.includeCompiledModuleFiles) {
                        break;
                    }

                    this.#addCompiledFile(
                        builder,
                        APW.fileType.module,
                        filePath,
                        locatedFileReferences,
                    );

                    break;
                default:
                    this.#addFile(builder, filePath, this.#extraFilesDirectory);
            }
        }

        return this;
    }

    #addWorkspaceFiles(builder, apw, options) {
        for (const file of apw.allFiles) {
            try {
                if (file.type === APW.fileType.workspace) {
                    this.#addFile(builder, file.path);

                    continue;
                }

                this.#addFile(builder, file.path, path.dirname(file.path));

                if (options.includeCompiledSourceFiles) {
                    if (file.type === APW.fileType.masterSrc) {
                        const compiledFilePath = file.path.replace(
                            APW.fileExtensions[APW.fileType.masterSrc],
                            APW.compiledFileExtensions[APW.fileType.masterSrc],
                        );

                        this.#addFile(
                            builder,
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

                        this.#addFile(
                            builder,
                            compiledFilePath,
                            path.dirname(compiledFilePath),
                        );
                    }
                }
            } catch (error) {
                console.log(chalk.red(error.message));
                continue;
            }
        }

        return this;
    }

    #createZip(builder, options) {
        console.log(chalk.blue("Compressing files into a zip archive..."));

        builder.writeZip(options.outputFile);

        console.log(chalk.green.bold(`Created archive: ${options.outputFile}`));

        return this;
    }

    async build() {
        const { builder, apw, options } = this;

        console.log(chalk.blue("Creating archive..."));

        this.#addWorkspaceFiles(builder, apw, options)
            .#addExtraFiles(builder, apw, options)
            .#createZip(builder, options);

        for (const entry of builder.getEntries()) {
            console.log(chalk.cyan(`--> ${entry.entryName}`));
        }
    }
}

export default ArchiveBuilder;
