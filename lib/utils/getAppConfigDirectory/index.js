import os from "os";
import path from "path";
import pkg from "../../../package.json";

const defaultAppConfigDir = path.join(
    os.homedir(),
    ".config",
    `${path.basename(pkg.name)}`,
);

export function getAppConfigDirectory() {
    return process.env.GENLINX_CONFIG_DIR || defaultAppConfigDir;
}

export default getAppConfigDirectory;
