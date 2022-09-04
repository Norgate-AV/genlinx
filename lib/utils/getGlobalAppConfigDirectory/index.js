import os from "os";
import path from "path";
import pkg from "../../../package.json";

const defaultGlobalAppConfigDirectory = path.join(
    os.homedir(),
    ".config",
    `${path.basename(pkg.name)}`,
);

export function getGlobalAppConfigDirectory() {
    return process.env.GENLINX_CONFIG_DIR || defaultGlobalAppConfigDirectory;
}

export default getGlobalAppConfigDirectory;
