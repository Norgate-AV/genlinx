import path from "node:path";
import chalk from "chalk";
import AdmZip from "adm-zip";
import { AmxExtensions, AmxCompiledExtensions } from "./index.js";
import {
    ArchiveOptions,
    ArchiveItem,
    ArchiveFileType as FileType,
    File,
} from "./@types/index.js";

export class ModuleItem implements ArchiveItem {
    private readonly file: File;
    private readonly builder: AdmZip;
    private readonly options: ArchiveOptions;
    private readonly archivePath: string;

    public constructor(builder: AdmZip, file: File, options: ArchiveOptions) {
        this.file = file;
        this.builder = builder;
        this.options = options;
        this.archivePath = file.isExtra
            ? options.extraFileArchiveLocation
            : path.dirname(file.path);
    }

    private addItem(file: File): void {
        const { builder, archivePath } = this;
        const { verbose } = this.options;

        builder.addLocalFile(file.path, archivePath);
        verbose && console.log(chalk.green(`Added file: ${file.path}`));
    }

    private static getCompiledFile(file: File): File {
        const compiledFilePath = file.path.replace(
            AmxExtensions[FileType.APW.Module],
            AmxCompiledExtensions[FileType.APW.Module],
        );

        return {
            ...file,
            path: compiledFilePath,
        };
    }

    public addToArchive(): void {
        const { file } = this;
        const { includeCompiledModuleFiles } = this.options;

        this.addItem(file);

        if (!includeCompiledModuleFiles) {
            return;
        }

        const compiledFile = ModuleItem.getCompiledFile(file);
        this.addItem(compiledFile);
    }
}
export default ModuleItem;
