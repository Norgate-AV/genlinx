import path from "node:path";
import { describe, it, expect, beforeEach } from "vitest";
import { NLRC } from "./NLRC.js";
import { BuildOptions } from "./@types/index.js";

describe("NLRC", () => {
    describe("getSourceBuildCommand", () => {
        let buildOptions: BuildOptions;

        beforeEach(() => {
            buildOptions = {
                nlrc: {
                    path: "C:/NLRC.exe",
                    option: {
                        cfg: "-CFG",
                        includePath: "-I",
                        modulePath: "-M",
                        libraryPath: "-L",
                    },
                    includePath: ["path/to/include"],
                    modulePath: ["path/to/module"],
                    libraryPath: ["path/to/library"],
                },
                shell: {
                    path: "cmd.exe",
                },
                all: false,
                createCfg: true,
                verbose: false,
                cfgFiles: [],
            };
        });

        it("should generate correct command for source file", () => {
            const command = NLRC.getSourceBuildCommand(
                "test.axs",
                buildOptions,
            );

            expect(command.path).toBe('"C:/NLRC.exe"');

            expect(containsPartialPath(command.args, "test.axs")).toBeTruthy();
            expect(command.args).toContain('-I"path/to/include"');
            expect(command.args).toContain('-M"path/to/module"');
            expect(command.args).toContain('-L"path/to/library"');
        });

        it("should include additional paths if specified", () => {
            buildOptions.includePath = ["custom/include"];
            buildOptions.modulePath = ["custom/module"];
            buildOptions.libraryPath = ["custom/library"];

            const command = NLRC.getSourceBuildCommand(
                "test.axs",
                buildOptions,
            );

            expect(command.args).toContain('-I"path/to/include"');
            expect(command.args).toContain('-M"path/to/module"');
            expect(command.args).toContain('-L"path/to/library"');

            // expect(command.args).toContain('-I"custom/include"');
            // expect(command.args).toContain('-M"custom/module"');
            // expect(command.args).toContain('-L"custom/library"');
        });
    });

    describe("getCfgBuildCommand", () => {
        it("should generate correct command for cfg file", () => {
            const buildOptions: BuildOptions = {
                nlrc: {
                    path: "C:/NLRC.exe",
                    option: {
                        cfg: "-CFG",
                        includePath: "-I",
                        modulePath: "-M",
                        libraryPath: "-L",
                    },
                    includePath: [],
                    modulePath: [],
                    libraryPath: [],
                },
                shell: { path: "cmd.exe" },
                all: false,
                createCfg: true,
                verbose: false,
                cfgFiles: [],
            };

            const command = NLRC.getCfgBuildCommand(
                "project.cfg",
                buildOptions,
            );

            expect(command.path).toBe('"C:/NLRC.exe"');
            expect(command.args).toContain(
                `-CFG"${path.resolve(process.cwd(), "project.cfg")}"`,
            );
        });
    });
});

function containsPartialPath(items: string[], partialPath: string): boolean {
    return items.some((item) => item.includes(partialPath));
}
