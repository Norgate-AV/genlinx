import path from "path";

export function getLocalAppConfigFilePath(directory) {
    const configFileName = ".genlinxrc.json";
    return path.join(`${directory}`, `${configFileName}`);
}

export default getLocalAppConfigFilePath;
