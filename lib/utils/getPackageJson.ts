import fs from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { findUp } from "find-up";
import { PackageJson } from "type-fest";

export async function getPackageJson(): Promise<PackageJson | null> {
    const filepath = await findUp("package.json", {
        cwd: path.dirname(fileURLToPath(import.meta.url)),
    });

    if (!filepath) {
        return null;
    }

    return JSON.parse(await fs.readFile(filepath, "utf-8"));
}
