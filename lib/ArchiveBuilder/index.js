import fs from "fs-extra";
import path from "path";
import chalk from "chalk";
import AdmZip from "adm-zip";
import { APW } from "..";
import { walkDirectory } from "../utils";
import { ArchiveItemFactory } from "../ArchiveItems";

export class ArchiveBuilder {
    constructor(apw, options) {
        this.apw = apw;
        this.options = options;
        this.builder = new AdmZip();
    }

    static #getExtraFilesOnDisk(locations) {
        console.log(chalk.blue("Searching known locations for extra files..."));

        const files = [];

        for (const location of locations) {
            console.log(chalk.cyan(`--> Searching ${location}`));
            files.push(...walkDirectory(location));
        }

        const filteredFiles = files.filter(
            (file) =>
                path.extname(file) ===
                    APW.fileExtensions[APW.fileType.module] ||
                APW.fileExtensions[APW.fileType.include] ||
                APW.fileExtensions[APW.fileType.duet] ||
                APW.fileExtensions[APW.fileType.xdd],
        );

        return filteredFiles;
    }

    #displayExtraWorkspaceFileReferences(references) {
        console.log(
            chalk.green(`Found ${references.length} extra files referenced...`),
        );

        for (const reference of references) {
            console.log(chalk.cyan(`--> ${reference}`));
        }

        return this;
    }

    #addSymlinkScripts(builder, options) {
        console.log(
            chalk.blue(
                "Adding symlink scripts for extra files to the archive...",
            ),
        );

        const scriptPaths = fs.readdirSync(
            path.resolve(__dirname, path.join("..", "..", "scripts")),
        );

        for (const scriptPath of scriptPaths) {
            const file = {
                path: scriptPath,
                extra: true,
            };

            try {
                ArchiveItemFactory.create(
                    builder,
                    file,
                    options,
                ).addToArchive();
            } catch (error) {
                console.log(chalk.red(error.message));
                continue;
            }
        }
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

        const extraWorkspaceFileReferences = apw.getExtraFileReferences();

        if (!extraWorkspaceFileReferences.length) {
            console.log(
                chalk.blue(
                    "No extra file references found in the workspace file",
                ),
            );

            return this;
        }

        this.#displayExtraWorkspaceFileReferences(extraWorkspaceFileReferences);

        const { extraFileSearchLocations } = options;
        const extraFilesOnDisk = ArchiveBuilder.#getExtraFilesOnDisk(
            extraFileSearchLocations,
        );

        for (const workspaceFileReference of extraWorkspaceFileReferences) {
            const filePath = extraFilesOnDisk.find((file) =>
                file.includes(workspaceFileReference),
            );

            if (!filePath) {
                console.log(
                    chalk.red(`Could not find ${workspaceFileReference}`),
                );

                continue;
            }

            console.log(
                chalk.green(`Found ${workspaceFileReference}: ${filePath}`),
            );

            const fileType = APW.getExtraFileType(filePath);
            console.log(
                chalk.blue(
                    `File ${workspaceFileReference} is of type "${fileType}"`,
                ),
            );

            const file = {
                type: fileType,
                path: filePath,
                extra: true,
            };

            try {
                ArchiveItemFactory.create(
                    builder,
                    file,
                    options,
                ).addToArchive();
            } catch (error) {
                console.log(chalk.red(error.message));
                continue;
            }
        }

        this.#addSymlinkScripts(builder, options);

        return this;
    }

    #addWorkspaceFiles(builder, apw, options) {
        for (const file of apw.allFiles) {
            try {
                ArchiveItemFactory.create(
                    builder,
                    file,
                    options,
                ).addToArchive();
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
