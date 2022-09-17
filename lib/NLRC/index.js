import path from "path";

export class NLRC {
    static getCfgBuildCommand(file, options) {
        const filePath = path.isAbsolute(file)
            ? file
            : path.resolve(process.cwd(), file);

        const { nlrc } = options;

        return {
            path: `"${nlrc.path}"`,
            args: [`${nlrc.option.cfg}"${filePath}"`],
        };
    }

    static getSourceBuildCommand(file, options) {
        const filePath = path.isAbsolute(file)
            ? file
            : path.resolve(process.cwd(), file);

        const { nlrc, includePath, modulePath, libraryPath } = options;

        const args = [`"${filePath}"`];

        if (includePath) {
            args.push(`${nlrc.option.includePath}"${includePath.join(";")}"`);
        }

        if (modulePath) {
            args.push(`${nlrc.option.modulePath}"${modulePath.join(";")}"`);
        }

        if (libraryPath) {
            args.push(`${nlrc.option.libraryPath}"${libraryPath.join(";")}"`);
        }

        return {
            path: `"${nlrc.path}"`,
            args: [...args],
        };
    }
}

export default NLRC;
