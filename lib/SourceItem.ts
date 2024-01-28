import path from "path";
import chalk from "chalk";
import AdmZip from "adm-zip";
import { AmxExtensions, AmxCompiledExtensions } from "./index.js";
import {
    ArchiveConfig,
    ArchiveItem,
    ArchiveFileType as FileType,
    File,
} from "./@types/index.js";

export class SourceItem implements ArchiveItem {
    private readonly file: File;
    private readonly builder: AdmZip;
    private readonly options: ArchiveConfig;
    private readonly archivePath: string;

    public constructor(builder: AdmZip, file: File, options: ArchiveConfig) {
        this.file = file;
        this.builder = builder;
        this.options = options;
        this.archivePath = file.isExtra
            ? options.extraFileArchiveLocation
            : path.dirname(file.path);
    }

    private addItem(file: File): void {
        const { builder, archivePath } = this;

        builder.addLocalFile(file.path, archivePath);
        console.log(chalk.green(`Added file: ${file.path}`));
    }

    private static getCompiledFile(file: File): File {
        const compiledFilePath = file.path.replace(
            AmxExtensions[FileType.APW.Source],
            AmxCompiledExtensions[FileType.APW.Source],
        );

        return {
            ...file,
            path: compiledFilePath,
        };
    }

    public addToArchive(): void {
        const { file } = this;
        const { includeCompiledSourceFiles } = this.options;

        this.addItem(file);

        if (!includeCompiledSourceFiles) {
            return;
        }

        const compiledFile = SourceItem.getCompiledFile(file);
        this.addItem(compiledFile);
    }
}
export default SourceItem;
