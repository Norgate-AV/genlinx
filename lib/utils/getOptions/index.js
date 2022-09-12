import mergician from "mergician";

export function getOptions(cliOptions, localConfig, globalConfig) {
    const options = mergician({
        appendArrays: true,
        dedupArrays: true,
        sortArrays: true,
    })(globalConfig, localConfig, cliOptions);

    return options;
}

export default getOptions;
