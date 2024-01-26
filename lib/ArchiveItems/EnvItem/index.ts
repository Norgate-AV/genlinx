import path from "path";
import chalk from "chalk";
import { ArchiveItem } from "../../@types/ArchiveItem";
import { FileReference } from "../../@types/FileReference";
import AdmZip from "adm-zip";
import { ArchiveConfig } from "../../@types";

export class EnvItem implements ArchiveItem {
    private readonly file: FileReference;
    private readonly builder: AdmZip;
    private readonly options: ArchiveConfig;

    public constructor(
        builder: AdmZip,
        file: FileReference,
        options: ArchiveConfig,
    ) {
        this.file = file;
        this.builder = builder;
        this.options = options;
    }

    private addItem(file: FileReference): void {
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