import os from "os";
import chalk from "chalk";
import { Command, Option } from "commander";
import { actions } from "../actions/index.js";

export function build(): Command {
    const command = new Command("build");

    command
        .description("build a NetLinx workspace or individual source file")
        .addOption(
            new Option(
                "-c, --cfg-files <string...>",
                "cfg file(s) to build from",
            )
                .default([], "search for CFG files in current directory")
                .conflicts([
                    "sourceFile",
                    "includePath",
                    "modulePath",
                    "libraryPath",
                ]),
        )
        .addOption(
            new Option(
                "-s, --source-file <string>",
                "axs source file to build",
            ).conflicts(["cfgFiles"]),
        )
        .addOption(
            new Option(
                "-i, --include-path <string...>",
                "add additional include paths",
            ).conflicts(["cfgFiles"]),
        )
        .addOption(
            new Option(
                "-m, --module-path <string...>",
                "add additional module paths",
            ).conflicts(["cfgFiles"]),
        )
        .addOption(
            new Option(
                "-l, --library-path <string...>",
                "add additional library paths",
            ).conflicts(["cfgFiles"]),
        )
        .addOption(
            new Option(
                "-a, --all",
                "select all cfg files without prompting",
            ).conflicts([
                "sourceFile",
                "includePath",
                "modulePath",
                "libraryPath",
            ]),
        )
        .addOption(
            new Option(
                "-A, --no-all",
                "select multiple cfg files with a prompt",
            ).conflicts([
                "sourceFile",
                "includePath",
                "modulePath",
                "libraryPath",
            ]),
        )
        .addOption(new Option("--verbose", "verbose output").default(false))
        .action((options) => {
            if (os.platform() !== "win32") {
                console.log(
                    chalk.red(
                        "The build command is only supported on Windows.",
                    ),
                );

                process.exit(1);
            }

            actions.build.execute(options);
        });

    return command;
}

export default build;
