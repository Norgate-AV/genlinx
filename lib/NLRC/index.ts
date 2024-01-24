import path from "path";
import { BuildConfig } from "../@types";

export class NLRC {
    static getCfgBuildCommand(file: string, options: BuildConfig) {
        const filePath = path.isAbsolute(file)
            ? file
            : path.resolve(process.cwd(), file);

        const { nlrc } = options;

        return {
            path: `"${nlrc.path}"`,
            args: [`${nlrc.option.cfg}"${filePath}"`],
        };
    }

    static getSourceBuildCommand(file: string, options: BuildConfig) {
        const filePath = path.isAbsolute(file)
            ? file
            : path.resolve(process.cwd(), file);

        const { nlrc } = options;

        const args = [`"${filePath}"`];

        if (nlrc.includePath) {
            args.push(
                `${nlrc.option.includePath}"${nlrc.includePath.join(";")}"`,
            );
        }

        if (nlrc.modulePath) {
            args.push(
                `${nlrc.option.modulePath}"${nlrc.modulePath.join(";")}"`,
            );
        }

        if (nlrc.libraryPath) {
            args.push(
                `${nlrc.option.libraryPath}"${nlrc.libraryPath.join(";")}"`,
            );
        }

        return {
            path: `"${nlrc.path}"`,
            args: [...args],
        };
    }
}

export default NLRC;
