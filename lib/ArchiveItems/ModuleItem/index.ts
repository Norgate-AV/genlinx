import path from "path";
import chalk from "chalk";
import { APW } from "../..";

export class ModuleItem {
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

    static #getCompiledFile(file) {
        const compiledFilePath = file.path.replace(
            APW.fileExtensions[APW.fileType.module],
            APW.compiledFileExtensions[APW.fileType.module],
        );

        return {
            ...file,
            path: compiledFilePath,
        };
    }

    addToArchive() {
        const { file } = this;
        const { includeCompiledModuleFiles } = this.options;

        this.#addItem(file);

        if (!includeCompiledModuleFiles) {
            return;
        }

        const compiledFile = ModuleItem.#getCompiledFile(file);
        this.#addItem(compiledFile);
    }
}
export default ModuleItem;
