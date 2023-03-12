import mergician from "mergician";
import { defaultGlobalAppConfig } from "../../../config";

export function getOptions(cliOptions, localConfig, globalConfig) {
    const options = mergician({
        appendArrays: true,
        dedupArrays: true,
        sortArrays: true,
    })(defaultGlobalAppConfig, globalConfig, localConfig, cliOptions);

    return options;
}

export default getOptions;
