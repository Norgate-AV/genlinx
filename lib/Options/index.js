import path from "path";
import mergician from "mergician";

export class Options {
    static getCfgOptions(cliOptions, localConfig, globalConfig) {
        const options = mergician({
            appendArrays: true,
            dedupArrays: true,
            sortArrays: true,
        })(globalConfig, localConfig, cliOptions);

        console.log(`options: ${JSON.stringify(options, null, 4)}`);
        return {
            ...options,
            rootDirectory: path.resolve(options.rootDirectory),
        };
    }

    static getArchiveOptions(cliOptions, localConfig, globalConfig) {
        const options = mergician({
            appendArrays: true,
            dedupArrays: true,
            sortArrays: true,
        })(globalConfig, localConfig, cliOptions);

        console.log(`options: ${JSON.stringify(options, null, 4)}`);
        return options;
    }

    static getBuildOptions(cliOptions, localConfig, globalConfig) {
        const options = mergician({
            appendArrays: true,
            dedupArrays: true,
            sortArrays: true,
        })(globalConfig, localConfig, cliOptions);

        console.log(`options: ${JSON.stringify(options, null, 4)}`);
        return options;
    }
}

export default Options;
