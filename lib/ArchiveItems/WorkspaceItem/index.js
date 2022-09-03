import chalk from "chalk";

export class WorkspaceItem {
    constructor(builder, file, options) {
        this.file = file;
        this.builder = builder;
        this.options = options;
    }

    #addItem(file) {
        const { builder } = this;

        builder.addLocalFile(file.path);
        console.log(chalk.green(`Added file: ${file.path}`));
    }

    addToArchive() {
        const { file } = this;

        this.#addItem(file);
    }
}
export default WorkspaceItem;
