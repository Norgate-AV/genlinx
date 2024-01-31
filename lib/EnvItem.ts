import path from "path";
import chalk from "chalk";
import AdmZip from "adm-zip";
import { ArchiveOptions, ArchiveItem, File } from "./@types/index.js";

export class EnvItem implements ArchiveItem {
    private readonly file: File;
    private readonly builder: AdmZip;
    private readonly options: ArchiveOptions;

    public constructor(builder: AdmZip, file: File, options: ArchiveOptions) {
        this.file = file;
        this.builder = builder;
        this.options = options;
    }

    private addItem(file: File): void {
        if (!file.content) {
            throw new Error("File content is not defined.");
        }

        const { builder, options } = this;
        const { extraFileArchiveLocation } = options;

        const filePath = path.join(extraFileArchiveLocation, file.path);
        const content = Buffer.from(file.content, "utf8");

        builder.addFile(filePath, content);
        console.log(chalk.green(`Added file: ${filePath}`));
    }

    addToArchive() {
        const { file } = this;

        this.addItem(file);
    }
}

export default EnvItem;
