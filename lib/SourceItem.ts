import path from "path";
import chalk from "chalk";
import AdmZip from "adm-zip";
import { extensions, compiledExtensions } from "./index.js";
import {
    ArchiveConfig,
    ArchiveItem,
    FileReference,
    FileType,
} from "./@types/index.js";

export class SourceItem implements ArchiveItem {
    private readonly file: FileReference;
    private readonly builder: AdmZip;
    private readonly options: ArchiveConfig;
    private readonly archivePath: string;

    public constructor(
        builder: AdmZip,
        file: FileReference,
        options: ArchiveConfig,
    ) {
        this.file = file;
        this.builder = builder;
        this.options = options;
        this.archivePath = file.extra
            ? options.extraFileArchiveLocation
            : path.dirname(file.path);
    }

    private addItem(file: FileReference): void {
        const { builder, archivePath } = this;

        builder.addLocalFile(file.path, archivePath);
        console.log(chalk.green(`Added file: ${file.path}`));
    }

    private static getCompiledFile(file: FileReference): FileReference {
        const compiledFilePath = file.path.replace(
            extensions[FileType.Source],
            compiledExtensions[FileType.Source],
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
