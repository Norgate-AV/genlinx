import fs from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";
import chalk from "chalk";
import AdmZip from "adm-zip";
import { walkDirectory } from "./utils/index.js";
import { APW, ArchiveItemFactory, AmxExtensions } from "./index.js";
import {
    ArchiveOptions,
    ArchiveFileType as FileType,
    File,
} from "./@types/index.js";

export class ArchiveBuilder {
    private readonly apw: APW;
    private readonly options: ArchiveOptions;
    private readonly builder = new AdmZip();
    private readonly extraFilesOnDisk: Array<string> = [];
    private readonly extraFileReferences: Array<string> = [];
    private readonly locatedExtraFileReferences: Array<string> = [];

    public constructor(apw: APW, options: ArchiveOptions) {
        this.apw = apw;
        this.options = options;
    }

    private static fileIsOfInterest(file: string) {
        return (
            path.extname(file) === AmxExtensions[FileType.APW.Module] ||
            path.extname(file) === AmxExtensions[FileType.APW.Include] ||
            path.extname(file) === AmxExtensions[FileType.APW.Duet] ||
            path.extname(file) === AmxExtensions[FileType.APW.XDD]
        );
    }

    private addItem(file: File): this {
        const { builder, options } = this;

        ArchiveItemFactory.create(builder, file, options).addToArchive();

        return this;
    }

    private async getExtraFilesOnDisk(locations: Array<string>): Promise<this> {
        const { extraFilesOnDisk } = this;
        const { verbose } = this.options;

        verbose &&
            console.log(
                chalk.blue("Searching known locations for extra files..."),
            );

        const files: Array<string> = [];

        for (const location of locations) {
            verbose && console.log(chalk.cyan(`--> Searching ${location}`));
            files.push(...(await walkDirectory(location)));
        }

        extraFilesOnDisk.push(
            ...files.filter((file) => ArchiveBuilder.fileIsOfInterest(file)),
        );

        return this;
    }

    private displayExtraFileReferences(references: Array<string>): this {
        const { verbose } = this.options;

        verbose &&
            console.log(
                chalk.green(
                    `Found ${references.length} extra files referenced...`,
                ),
            );

        for (const reference of references) {
            verbose && console.log(chalk.cyan(`--> ${reference}`));
        }

        return this;
    }

    private addEnvFile(): this {
        const { apw } = this;
        const { verbose } = this.options;

        if (!apw.masterSrcPath[0]) {
            return this;
        }

        verbose && console.log(chalk.blue("Adding env file to the archive..."));

        const masterSrcPath = path.join(
            "..",
            path.basename(apw.masterSrcPath[0]),
        );

        const file: File = {
            id: "",
            type: FileType.Env,
            path: ".env",
            exists: true,
            isExtra: true,
            content: `SOURCE_DIRECTORY_RELATIVE_PATH=${masterSrcPath}`,
        };

        try {
            this.addItem(file);
        } catch (error: any) {
            console.log(chalk.red(error.message));
        }

        return this;
    }

    private async addSymlinkScripts(): Promise<this> {
        const { verbose } = this.options;

        verbose &&
            console.log(
                chalk.blue(
                    "Adding symlink scripts for extra files to the archive...",
                ),
            );

        const __dirname = path.dirname(fileURLToPath(import.meta.url));
        const scriptsPath = path.resolve(__dirname, path.join("..", "scripts"));

        const scripts = await fs.readdir(scriptsPath);

        for (const script of scripts) {
            const file: File = {
                id: "",
                type: FileType.Script,
                path: path.join(scriptsPath, script),
                exists: true,
                isExtra: true,
            };

            try {
                this.addItem(file);
            } catch (error: any) {
                console.log(chalk.red(error.message));
                continue;
            }
        }

        return this;
    }

    private async searchForExtraFiles(
        references: Array<string>,
    ): Promise<this> {
        const { extraFilesOnDisk, locatedExtraFileReferences } = this;
        const { verbose } = this.options;

        const newLocatedReferences: Array<string> = [];

        for (const reference of references) {
            const filePath = extraFilesOnDisk.find((file) =>
                file.includes(reference),
            );

            if (!filePath) {
                verbose &&
                    console.log(chalk.red(`Could not find ${reference}`));
                continue;
            }

            const { ignoredFiles } = this.options;
            const isIgnored = ignoredFiles.includes(path.basename(filePath));

            if (isIgnored) {
                verbose &&
                    console.log(
                        chalk.blue(
                            `Ignoring ${path.basename(filePath)} as per config`,
                        ),
                    );

                continue;
            }

            verbose &&
                console.log(chalk.green(`Found ${reference}: ${filePath}`));
            newLocatedReferences.push(filePath);
        }

        if (newLocatedReferences.length === 0) {
            return this;
        }

        locatedExtraFileReferences.push(...newLocatedReferences);
        await this.getFileReferences(newLocatedReferences);
        return this;
    }

    private async getFileReferences(files: Array<string>): Promise<this> {
        const { extraFileReferences, apw } = this;
        const { verbose } = this.options;

        for (const file of files) {
            if (!APW.fileIsReadable(file)) {
                continue;
            }

            verbose &&
                console.log(chalk.blue(`Searching ${file} for references...`));
            const references = await apw.getExtraFileReferencesFromFile(file);

            if (references.length === 0) {
                verbose && console.log(chalk.cyan("--> No references found"));
                continue;
            }

            const newReferences = references.filter(
                (reference) => !extraFileReferences.includes(reference),
            );

            if (newReferences.length === 0) {
                verbose &&
                    console.log(chalk.cyan("--> No new references found"));
                continue;
            }

            extraFileReferences.push(...newReferences);
            this.displayExtraFileReferences(newReferences);
            await this.searchForExtraFiles(newReferences);
        }

        return this;
    }

    private async addExtraFiles(): Promise<this> {
        const { apw, options } = this;
        const { verbose } = options;

        if (!options.includeFilesNotInWorkspace) {
            return this;
        }

        verbose &&
            console.log(
                chalk.blue(
                    "Searching for extra files that are not part of the workspace...",
                ),
            );

        const { extraFileReferences } = this;
        extraFileReferences.push(...(await apw.getExtraFileReferences()));

        if (extraFileReferences.length === 0) {
            verbose &&
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

        this.displayExtraFileReferences(extraFileReferences);
        await this.getExtraFilesOnDisk(extraFileSearchLocations);
        await this.searchForExtraFiles(extraFileReferences);

        const { locatedExtraFileReferences } = this;
        if (locatedExtraFileReferences.length === 0) {
            return this;
        }

        for (const reference of locatedExtraFileReferences) {
            const file: File = {
                id: "",
                type: APW.getFileType(reference),
                path: reference,
                exists: true,
                isExtra: true,
            };

            try {
                this.addItem(file);
            } catch (error: any) {
                console.log(chalk.red(error.message));
                continue;
            }
        }

        this.addEnvFile();
        await this.addSymlinkScripts();

        return this;
    }

    private addWorkspaceFiles(): this {
        const { apw } = this;

        for (const file of apw.allFiles) {
            try {
                this.addItem(file);
            } catch (error: any) {
                console.log(chalk.red(error.message));
                continue;
            }
        }

        return this;
    }

    private createZip(): this {
        const { builder, apw, options } = this;
        const { verbose } = options;

        verbose &&
            console.log(chalk.blue("Compressing files into a zip archive..."));

        const outputFile = `${apw.id}.${options.outputFileSuffix}`;
        builder.writeZip(outputFile);

        console.log(chalk.green.bold(`Created archive: ${outputFile}`));

        return this;
    }

    private displayZippedFiles(): this {
        const { builder } = this;

        for (const entry of builder.getEntries()) {
            console.log(chalk.cyan(`--> ${entry.entryName}`));
        }

        return this;
    }

    public async build(): Promise<void> {
        const { verbose } = this.options;

        verbose && console.log(chalk.blue("Creating archive..."));

        this.addWorkspaceFiles();

        await this.addExtraFiles();

        this.createZip();

        verbose && this.displayZippedFiles();
    }
}

export default ArchiveBuilder;
