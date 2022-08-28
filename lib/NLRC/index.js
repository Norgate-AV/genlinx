import path from "path";

export class NLRC {
    static getCfgBuildCommand(cfgFilePath, options, config) {
        const filePath = path.isAbsolute(cfgFilePath)
            ? cfgFilePath
            : path.resolve(process.cwd(), cfgFilePath);

        const { nlrc } = config.build;

        return `"${nlrc.path}" -C"${filePath}"`;
    }
}

export default NLRC;
