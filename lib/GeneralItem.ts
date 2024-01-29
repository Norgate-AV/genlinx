import path from "path";
import chalk from "chalk";
import AdmZip from "adm-zip";
import { ArchiveConfig, ArchiveItem, File } from "./@types/index.js";

export class GeneralItem implements ArchiveItem {
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

    public addToArchive(): void {
        const { file } = this;

        this.addItem(file);
    }
}

export default GeneralItem;