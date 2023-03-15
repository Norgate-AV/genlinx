import path from "path";
import Conf from "conf";
import { getLocalAppConfigFilePath } from "../getLocalAppConfigFilePath";
import { schema } from "../../../config/schema";

export async function getLocalAppConfig(cwd) {
    const configFilePath = await getLocalAppConfigFilePath(cwd);

    if (!configFilePath) {
        return {};
    }

    const localConfig = new Conf({
        schema,
        configName: path.basename(configFilePath, ".json"),
        cwd: path.dirname(configFilePath),
    });

    return {
        path: configFilePath,
        config: localConfig,
    };
}

export default getLocalAppConfig;
