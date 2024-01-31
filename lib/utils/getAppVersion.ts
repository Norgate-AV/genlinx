import fs from "fs-extra";
import { fileURLToPath } from "url";
import { findUp } from "find-up";
import { PackageJson } from "type-fest";

export async function getAppVersion(): Promise<string> {
    const packageJsonPath = await findUp("package.json", {
        cwd: fileURLToPath(import.meta.url),
    });

    if (!packageJsonPath) {
        return "Unknown version";
    }

    const packageJson: PackageJson = JSON.parse(
        await fs.readFile("package.json", "utf-8"),
    );

    return packageJson.version || "Unknown version";
}
