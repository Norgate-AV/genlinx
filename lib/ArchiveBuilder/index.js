import fs from "fs-extra";
import path from "path";
import chalk from "chalk";
import AdmZip from "adm-zip";
import { APW } from "..";
import { walkDirectory } from "../utils";
import { ArchiveItemFactory, FileType } from "../ArchiveItems";

export class ArchiveBuilder {
    extraFilesOnDisk = [];

    extraFileReferences = [];

    locatedExtraFileReferences = [];

    constructor(apw, options) {
        this.apw = apw;
        this.options = options;
        this.builder = new AdmZip();
    }

    #getExtraFilesOnDisk(locations) {
        console.log(chalk.blue("Searching known locations for extra files..."));

        const files = [];

        for (const location of locations) {
            console.log(chalk.cyan(`--> Searching ${location}`));
            files.push(...walkDirectory(location));
        }

        const filteredFiles = files.filter(
            (file) =>
                path.extname(file) ===
                    APW.fileExtensions[FileType.apw.module] ||
                APW.fileExtensions[FileType.apw.include] ||
                APW.fileExtensions[FileType.apw.duet] ||
                APW.fileExtensions[FileType.apw.xdd],
        );

        this.extraFilesOnDisk = filteredFiles;
        return this;
    }

    #displayExtraFileReferences(references) {
        console.log(
            chalk.green(`Found ${references.length} extra files referenced...`),
        );

        for (const reference of references) {
            console.log(chalk.cyan(`--> ${reference}`));
        }

        return this;
    }

    #addEnvFile(builder, apw, options) {
        console.log(chalk.blue("Adding env file to the archive..."));

        const masterSrcPath = path.join(
            "..",
            path.basename(apw.masterSrcPath[0]),
        );

        const file = {
            type: FileType.env,
            path: ".env",
            content: `SOURCE_DIRECTORY_RELATIVE_PATH=${masterSrcPath}`,
        };

        try {
            ArchiveItemFactory.create(builder, file, options).addToArchive();
        } catch (error) {
            console.log(chalk.red(error.message));
        }

        return this;
    }

    #addSymlinkScripts(builder, apw, options) {
        console.log(
            chalk.blue(
                "Adding symlink scripts for extra files to the archive...",
            ),
        );

        const scriptsPath = path.resolve(
            __dirname,
            path.join("..", "..", "scripts"),
        );

        const scripts = fs.readdirSync(scriptsPath);

        for (const script of scripts) {
            const file = {
                path: path.join(scriptsPath, script),
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

    #searchForExtraFiles(references) {
        const { extraFilesOnDisk, locatedExtraFileReferences } = this;

        const newLocatedReferences = [];

        for (const reference of references) {
            const filePath = extraFilesOnDisk.find((file) =>
                file.includes(reference),
            );

            if (!filePath) {
                console.log(chalk.red(`Could not find ${reference}`));
                continue;
            }

            console.log(chalk.green(`Found ${reference}: ${filePath}`));
            newLocatedReferences.push(filePath);
        }

        if (!newLocatedReferences.length) {
            return this;
        }

        locatedExtraFileReferences.push(...newLocatedReferences);
        return this.#getFileReferences(newLocatedReferences);
    }

    #getFileReferences(files) {
        const { extraFileReferences, apw } = this;

        for (const file of files) {
            if (!APW.fileIsReadable(file)) {
                continue;
            }

            console.log(chalk.blue(`Searching ${file} for references...`));
            const references = apw.getExtraFileReferencesFromFile(file);

            if (!references.length) {
                console.log(chalk.cyan("--> No references found"));
                continue;
            }

            const newReferences = references.filter(
                (reference) => !extraFileReferences.includes(reference),
            );

            if (!newReferences.length) {
                console.log(chalk.cyan("--> No new references found"));
                continue;
            }

            extraFileReferences.push(...newReferences);
            this.#displayExtraFileReferences(
                newReferences,
            ).#searchForExtraFiles(newReferences);
        }

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

        const { extraFileReferences } = this;
        extraFileReferences.push(...apw.getExtraFileReferences());

        if (!extraFileReferences.length) {
            console.log(
                chalk.blue(
                    "No extra file references found in the workspace file",
                ),
            );

            return this;
        }

        const extraFileSearchLocations = [
            ...options.extraFileSearchLocations,
            path.dirname(apw.filePath),
        ];

        this.#displayExtraFileReferences(extraFileReferences)
            .#getExtraFilesOnDisk(extraFileSearchLocations)
            .#searchForExtraFiles(extraFileReferences);

        const { locatedExtraFileReferences } = this;
        if (!locatedExtraFileReferences.length) {
            return this;
        }

        for (const reference of locatedExtraFileReferences) {
            const fileType = APW.getFileType(reference);

            const file = {
                type: fileType,
                path: reference,
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

        this.#addEnvFile(builder, apw, options).#addSymlinkScripts(
            builder,
            apw,
            options,
        );

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

    #createZip(builder, apw, options) {
        console.log(chalk.blue("Compressing files into a zip archive..."));

        const outputFile = `${apw.id}.${options.outputFileSuffix}`;
        builder.writeZip(outputFile);

        console.log(chalk.green.bold(`Created archive: ${outputFile}`));

        return this;
    }

    async build() {
        const { builder, apw, options } = this;

        console.log(chalk.blue("Creating archive..."));

        this.#addWorkspaceFiles(builder, apw, options)
            .#addExtraFiles(builder, apw, options)
            .#createZip(builder, apw, options);

        for (const entry of builder.getEntries()) {
            console.log(chalk.cyan(`--> ${entry.entryName}`));
        }
    }
}

export default ArchiveBuilder;
