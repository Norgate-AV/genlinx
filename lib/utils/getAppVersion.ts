import { getPackageJson } from "./index.js";

export async function getAppVersion(): Promise<string> {
    return (await getPackageJson())?.version || "Unknown";
}
