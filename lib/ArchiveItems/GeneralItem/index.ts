import path from "path";
import chalk from "chalk";

export class GeneralItem {
    constructor(builder, file, options) {
        this.file = file;
        this.builder = builder;
        this.options = options;
        this.archivePath = file.extra
            ? options.extraFileArchiveLocation
            : path.dirname(file.path);
    }

    #addItem(file) {
        const { builder, archivePath } = this;

        builder.addLocalFile(file.path, archivePath);
        console.log(chalk.green(`Added file: ${file.path}`));
    }

    addToArchive() {
        const { file } = this;

        this.#addItem(file);
    }
}

export default GeneralItem;
