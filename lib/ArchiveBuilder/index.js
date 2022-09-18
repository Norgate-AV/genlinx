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

    static #fileIsOfInterest(file) {
        return (
            path.extname(file) === APW.fileExtensions[FileType.apw.module] ||
            path.extname(file) === APW.fileExtensions[FileType.apw.include] ||
            path.extname(file) === APW.fileExtensions[FileType.apw.duet] ||
            path.extname(file) === APW.fileExtensions[FileType.apw.xdd]
        );
    }

    #addItem(file) {
        const { builder, options } = this;

        ArchiveItemFactory.create(builder, file, options).addToArchive();

        return this;
    }

    async #getExtraFilesOnDisk(locations) {
        console.log(chalk.blue("Searching known locations for extra files..."));

        const files = [];

        for (const location of locations) {
            console.log(chalk.cyan(`--> Searching ${location}`));
            files.push(...(await walkDirectory(location)));
        }

        const filteredFiles = files.filter((file) =>
            ArchiveBuilder.#fileIsOfInterest(file),
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

    #addEnvFile() {
        const { apw } = this;

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
            this.#addItem(file);
        } catch (error) {
            console.log(chalk.red(error.message));
        }

        return this;
    }

    async #addSymlinkScripts() {
        console.log(
            chalk.blue(
                "Adding symlink scripts for extra files to the archive...",
            ),
        );

        const scriptsPath = path.resolve(
            __dirname,
            path.join("..", "..", "scripts"),
        );

        const scripts = await fs.readdir(scriptsPath);

        for (const script of scripts) {
            const file = {
                path: path.join(scriptsPath, script),
                extra: true,
            };

            try {
                this.#addItem(file);
            } catch (error) {
                console.log(chalk.red(error.message));
                continue;
            }
        }

        return this;
    }

    async #searchForExtraFiles(references) {
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

            const { ignoredFiles } = this.options;
            const isIgnored = ignoredFiles.includes(path.basename(filePath));

            if (isIgnored) {
                console.log(
                    chalk.blue(
                        `Ignoring ${path.basename(filePath)} as per config`,
                    ),
                );

                continue;
            }

            console.log(chalk.green(`Found ${reference}: ${filePath}`));
            newLocatedReferences.push(filePath);
        }

        if (newLocatedReferences.length === 0) {
            return this;
        }

        locatedExtraFileReferences.push(...newLocatedReferences);
        await this.#getFileReferences(newLocatedReferences);
        return this;
    }

    async #getFileReferences(files) {
        const { extraFileReferences, apw } = this;

        for (const file of files) {
            if (!APW.fileIsReadable(file)) {
                continue;
            }

            console.log(chalk.blue(`Searching ${file} for references...`));
            const references = await apw.getExtraFileReferencesFromFile(file);

            if (references.length === 0) {
                console.log(chalk.cyan("--> No references found"));
                continue;
            }

            const newReferences = references.filter(
                (reference) => !extraFileReferences.includes(reference),
            );

            if (newReferences.length === 0) {
                console.log(chalk.cyan("--> No new references found"));
                continue;
            }

            extraFileReferences.push(...newReferences);
            this.#displayExtraFileReferences(newReferences);
            await this.#searchForExtraFiles(newReferences);
        }

        return this;
    }

    async #addExtraFiles() {
        const { apw, options } = this;

        if (!options.includeFilesNotInWorkspace) {
            return this;
        }

        console.log(
            chalk.blue(
                "Searching for extra files that are not part of the workspace...",
            ),
        );

        const { extraFileReferences } = this;
        extraFileReferences.push(...(await apw.getExtraFileReferences()));

        if (extraFileReferences.length === 0) {
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

        this.#displayExtraFileReferences(extraFileReferences);
        await this.#getExtraFilesOnDisk(extraFileSearchLocations);
        await this.#searchForExtraFiles(extraFileReferences);

        const { locatedExtraFileReferences } = this;
        if (locatedExtraFileReferences.length === 0) {
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
                this.#addItem(file);
            } catch (error) {
                console.log(chalk.red(error.message));
                continue;
            }
        }

        this.#addEnvFile();
        await this.#addSymlinkScripts();

        return this;
    }

    #addWorkspaceFiles() {
        const { apw } = this;

        for (const file of apw.allFiles) {
            try {
                this.#addItem(file);
            } catch (error) {
                console.log(chalk.red(error.message));
                continue;
            }
        }

        return this;
    }

    #createZip() {
        const { builder, apw, options } = this;

        console.log(chalk.blue("Compressing files into a zip archive..."));

        const outputFile = `${apw.id}.${options.outputFileSuffix}`;
        builder.writeZip(outputFile);

        console.log(chalk.green.bold(`Created archive: ${outputFile}`));

        return this;
    }

    #displayZippedFiles() {
        const { builder } = this;

        for (const entry of builder.getEntries()) {
            console.log(chalk.cyan(`--> ${entry.entryName}`));
        }

        return this;
    }

    async build() {
        console.log(chalk.blue("Creating archive..."));

        this.#addWorkspaceFiles();

        await this.#addExtraFiles();

        this.#createZip().#displayZippedFiles();
    }
}

export default ArchiveBuilder;
