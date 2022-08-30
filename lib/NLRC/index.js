import path from "path";

export class NLRC {
    static getCfgBuildCommand(cfgFilePath, options) {
        const filePath = path.isAbsolute(cfgFilePath)
            ? cfgFilePath
            : path.resolve(process.cwd(), cfgFilePath);

        const { nlrc } = options;

        return {
            path: `"${nlrc.path}"`,
            args: [...nlrc.args, `"${filePath}"`],
        };
    }
}

export default NLRC;
