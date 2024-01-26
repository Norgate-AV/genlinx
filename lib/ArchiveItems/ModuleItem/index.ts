import path from "path";
import chalk from "chalk";
import { APW } from "../..";
import { extensions, compiledExtensions } from "../../APW";
import { ArchiveItem } from "../../@types/ArchiveItem";
import { ArchiveConfig, FileReference } from "../../@types";
import { FileType } from "../../@types/FileType";
import AdmZip from "adm-zip";

export class ModuleItem implements ArchiveItem {
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
            extensions[FileType.Module],
            compiledExtensions[FileType.Module],
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
