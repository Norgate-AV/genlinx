import path from "path";
import chalk from "chalk";
import AdmZip from "adm-zip";
import { ArchiveConfig, ArchiveItem, FileReference } from "./@types/index.js";

export class GeneralItem implements ArchiveItem {
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

    public addToArchive(): void {
        const { file } = this;

        this.addItem(file);
    }
}

export default GeneralItem;
