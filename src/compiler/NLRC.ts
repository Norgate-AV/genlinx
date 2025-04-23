import path from "node:path";
import { BuildConfig, ShellCommand } from "../@types/index.js";

export class NLRC {
    private static readonly COMPILER_FLAGS = {
        cfg: "-CFG",
        includePath: "-I",
        modulePath: "-M",
        libraryPath: "-L",
    };

    static getCfgBuildCommand(
        file: string,
        options: BuildConfig,
    ): ShellCommand {
        const filePath = path.isAbsolute(file)
            ? file
            : path.resolve(process.cwd(), file);

        const { nlrc } = options;

        return {
            path: `"${nlrc.path}"`,
            args: [`${NLRC.COMPILER_FLAGS.cfg}"${filePath}"`],
        };
    }

    static getSourceBuildCommand(
        file: string,
        options: BuildConfig,
    ): ShellCommand {
        const filePath = path.isAbsolute(file)
            ? file
            : path.resolve(process.cwd(), file);

        const { nlrc } = options;

        const args = [`"${filePath}"`];

        args.push(
            `${NLRC.COMPILER_FLAGS.includePath}"${nlrc.includePath.join(";")}"`,
        );
        args.push(
            `${NLRC.COMPILER_FLAGS.modulePath}"${nlrc.modulePath.join(";")}"`,
        );
        args.push(
            `${NLRC.COMPILER_FLAGS.libraryPath}"${nlrc.libraryPath.join(";")}"`,
        );

        return {
            path: `"${nlrc.path}"`,
            args: [...args],
        };
    }
}

export default NLRC;
