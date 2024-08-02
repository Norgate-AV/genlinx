import path from "node:path";
import { getPackageJson } from "./index.js";

export async function getModuleName(): Promise<string> {
    return path.basename((await getPackageJson())?.name || "Unknown");
}
