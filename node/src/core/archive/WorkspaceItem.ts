import path from "node:path";
import chalk from "chalk";
import AdmZip from "adm-zip";
import { ArchiveOptions, ArchiveItem, File } from "../../types/index.js";

export class WorkspaceItem implements ArchiveItem {
    private readonly file: File;
    private readonly builder: AdmZip;
    private readonly options: ArchiveOptions;

    public constructor(builder: AdmZip, file: File, options: ArchiveOptions) {
        this.file = file;
        this.builder = builder;
        this.options = options;
    }

    private addItem(file: File): void {
        const { builder } = this;
        const { verbose } = this.options;

        builder.addLocalFile(file.path);
        verbose &&
            console.log(chalk.green(`Added file: ${path.basename(file.path)}`));
    }

    public addToArchive(): void {
        const { file } = this;

        this.addItem(file);
    }
}
export default WorkspaceItem;
