/* eslint-disable no-continue */
import path from "path";
import chalk from "chalk";
import AdmZip from "adm-zip";
import { APW } from "../APW";

export class ArchiveBuilder {
    constructor(apw, options) {
        this.apw = apw;
        this.options = options;
        this.builder = new AdmZip();
    }

    build() {
        const { builder, apw, options } = this;

        for (const file of apw.allFiles) {
            try {
                builder.addLocalFile(file.path, path.dirname(file.path));
                if (options.includeCompiledSourceFiles) {
                    if (file.type === APW.fileType.masterSrc) {
                        const compiledFilePath = file.path.replace(
                            APW.fileExtensions[APW.fileType.masterSrc],
                            APW.compiledFileExtensions[APW.fileType.masterSrc],
                        );

                        builder.addLocalFile(
                            compiledFilePath,
                            path.dirname(compiledFilePath),
                        );
                    }
                }

                if (options.includeCompiledModuleFiles) {
                    if (file.type === APW.fileType.module) {
                        const compiledFilePath = file.path.replace(
                            APW.fileExtensions[APW.fileType.module],
                            APW.compiledFileExtensions[APW.fileType.module],
                        );

                        builder.addLocalFile(
                            compiledFilePath,
                            path.dirname(compiledFilePath),
                        );
                    }
                }
            } catch (error) {
                console.log(chalk.red(error.message));
                continue;
            }

            console.log(chalk.green(`Added file: ${file.path}`));
        }

        builder.writeZip(options.outputFile);

        console.log(chalk.green.bold(`Created archive: ${options.outputFile}`));

        for (const entry of builder.getEntries()) {
            console.log(chalk.cyan(`--> ${entry.entryName}`));
        }
    }
}

export default ArchiveBuilder;
