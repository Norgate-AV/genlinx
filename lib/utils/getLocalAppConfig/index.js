import fs from "fs-extra";
import path from "path";
import findUp from "find-up";

const defaultLocalAppConfig = {
    cfg: {},
    archive: {},
    build: {},
};

export async function getLocalAppConfig(filePath) {
    try {
        const configFilePath = await findUp(".genlinxrc.json", {
            cwd: path.dirname(filePath),
        });

        if (!configFilePath) {
            return defaultLocalAppConfig;
        }

        const localConfig = JSON.parse(
            await fs.readFile(configFilePath, "utf8"),
        );

        return localConfig;
    } catch (error) {
        console.log(error);
    }

    return defaultLocalAppConfig;
}

export default getLocalAppConfig;
