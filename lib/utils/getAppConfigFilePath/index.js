import path from "path";
import { getAppConfigDirectory } from "../getAppConfigDirectory";

export function getAppConfigFilePath() {
    const configFileName = "config.json";
    const configDirectory = getAppConfigDirectory();
    return path.join(`${configDirectory}`, `${configFileName}`);
}

export default getAppConfigFilePath;
