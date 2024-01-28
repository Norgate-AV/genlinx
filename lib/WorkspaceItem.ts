import path from "path";
import chalk from "chalk";
import AdmZip from "adm-zip";
import { ArchiveConfig, ArchiveItem, FileReference } from "./@types/index.js";

export class WorkspaceItem implements ArchiveItem {
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
        const { builder } = this;

        builder.addLocalFile(file.path);
        console.log(chalk.green(`Added file: ${path.basename(file.path)}`));
    }

    public addToArchive(): void {
        const { file } = this;

        this.addItem(file);
    }
}
export default WorkspaceItem;
