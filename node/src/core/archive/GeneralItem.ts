import path from "node:path";
import chalk from "chalk";
import AdmZip from "adm-zip";
import { ArchiveOptions, ArchiveItem, File } from "../../types/index.js";

export class GeneralItem implements ArchiveItem {
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

    public addToArchive(): void {
        const { file } = this;

        this.addItem(file);
    }
}

export default GeneralItem;
