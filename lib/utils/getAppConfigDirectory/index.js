import os from "os";
import path from "path";
import pkg from "../../../package.json";

const defaultAppConfigDirectory = path.join(
    os.homedir(),
    ".config",
    `${path.basename(pkg.name)}`,
);

export function getAppConfigDirectory() {
    return process.env.GENLINX_CONFIG_DIR || defaultAppConfigDirectory;
}

export default getAppConfigDirectory;
