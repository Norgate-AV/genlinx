import mergician from "mergician";

export function getOptions(cliOptions, localConfig, globalConfig) {
    const options = mergician({
        appendArrays: true,
        dedupArrays: true,
        sortArrays: true,
        filter({ srcVal }) {
            return Boolean(srcVal) || srcVal === 0;
        },
    })(globalConfig, localConfig, cliOptions);

    console.log(options);

    return options;
}

export default getOptions;
