import path from "path";
import chalk from "chalk";

export class EnvItem {
    constructor(builder, file, options) {
        this.file = file;
        this.builder = builder;
        this.options = options;
    }

    #addItem(file) {
        const { builder, options } = this;
        const { extraFileArchiveLocation } = options;

        const filePath = path.join(extraFileArchiveLocation, file.path);
        const content = Buffer.from(file.content, "utf8");

        builder.addFile(filePath, content);
        console.log(chalk.green(`Added file: ${filePath}`));
    }

    addToArchive() {
        const { file } = this;

        this.#addItem(file);
    }
}
export default EnvItem;
