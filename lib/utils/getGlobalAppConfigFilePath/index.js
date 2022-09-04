import path from "path";
import { getGlobalAppConfigDirectory } from "../getGlobalAppConfigDirectory";

export function getGlobalAppConfigFilePath() {
    const configFileName = "config.json";
    const configDirectory = getGlobalAppConfigDirectory();
    return path.join(`${configDirectory}`, `${configFileName}`);
}

export default getGlobalAppConfigFilePath;
