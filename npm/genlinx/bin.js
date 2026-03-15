#!/usr/bin/env node
import { spawnSync } from "node:child_process";
import { createRequire } from "node:module";
import { arch, platform } from "node:os";
import { dirname, join } from "node:path";

const PLATFORMS = {
    "linux-x64": "@norgate-av/genlinx-linux-x64",
    "linux-arm64": "@norgate-av/genlinx-linux-arm64",
    "darwin-x64": "@norgate-av/genlinx-darwin-x64",
    "darwin-arm64": "@norgate-av/genlinx-darwin-arm64",
    "win32-x64": "@norgate-av/genlinx-win32-x64",
    "win32-arm64": "@norgate-av/genlinx-win32-arm64",
};

const key = `${platform()}-${arch()}`;
const pkg = PLATFORMS[key];

if (!pkg) {
    console.error(`genlinx: unsupported platform: ${key}`);
    process.exit(1);
}

const require = createRequire(import.meta.url);

let pkgDir;
try {
    pkgDir = dirname(require.resolve(`${pkg}/package.json`));
} catch {
    console.error(`genlinx: platform package ${pkg} is not installed.`);
    console.error(`  Try reinstalling: npm install @norgate-av/genlinx`);
    process.exit(1);
}

const binaryName = platform() === "win32" ? "genlinx.exe" : "genlinx";
const binaryPath = join(pkgDir, binaryName);

const result = spawnSync(binaryPath, process.argv.slice(2), {
    stdio: "inherit",
});
process.exit(result.status ?? 1);
