export class NLRC {
    #cfgFilePath;

    #compilerPath;

    constructor(cfgFilePath, config) {
        this.#cfgFilePath = cfgFilePath;
        this.#compilerPath = config.nlrc.path;
    }

    get buildCommand() {
        return `${this.#compilerPath} -C"${this.#cfgFilePath}"`;
    }
}

export default NLRC;
