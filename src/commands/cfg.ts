import { Command, Option } from "commander";
import { actions } from "../actions/index.js";

export function cfg(): Command {
    const command = new Command("cfg");

    command
        .description("generate NetLinx build CFG files")
        .addOption(
            new Option(
                "-w, --workspace-files <string...>",
                "workspace file(s) to generate a CFG for",
            ).default([], "search for workspace files in current directory"),
        )
        .addOption(
            new Option(
                "-r, --root-directory <string>",
                "root directory reference",
            ).default("", "use parent directory of workspace file as root"),
        )
        .addOption(
            new Option(
                "-o, --output-file-suffix <string>",
                "output file suffix",
            ),
        )
        .addOption(
            new Option(
                "-f, --output-log-file-suffix <string>",
                "output log file suffix",
            ),
        )
        .addOption(
            new Option(
                "-k, --output-log-file-option <string>",
                `
output log file option
A - append
N - overwrite`,
            ).choices(["A", "N"]),
        )
        .addOption(
            new Option(
                "-c, --output-log-console-option",
                "output log to console",
            ),
        )
        .addOption(
            new Option(
                "-C, --no-output-log-console-option",
                "do not output log to console",
            ),
        )
        .addOption(
            new Option(
                "-d, --build-with-debug-information",
                "build with debug information",
            ),
        )
        .addOption(
            new Option(
                "-D, --no-build-with-debug-information",
                "do not build with debug information",
            ),
        )
        .addOption(new Option("-s, --build-with-source", "build with source"))
        .addOption(
            new Option(
                "-S, --no-build-with-source",
                "do not build with source",
            ),
        )
        .addOption(
            new Option(
                "-i, --include-path <string...>",
                "add additional include paths",
            ),
        )
        .addOption(
            new Option(
                "-m, --module-path <string...>",
                "add additional module paths",
            ),
        )
        .addOption(
            new Option(
                "-l, --library-path <string...>",
                "add additional library paths",
            ),
        )
        .addOption(
            new Option(
                "-a, --all",
                `
if no workspace files are specified with the -w option and more than one
workspace file is found in the current directory, select all of them
without prompting`,
            ),
        )
        .addOption(
            new Option(
                "-A, --no-all",
                `
if no workspace files are specified with the -w option and more than one
workspace file is found in the current directory, prompt to select which
workspace files to use`,
            ),
        )
        .action((options) => actions.cfg.create(options))
        .addHelpText(
            "after",
            `
the CFG file is used to configure the NetLinx Studio build process

if an option is omitted, genlinx will:
    1. look for a .genlinxrc.json file in the root directory of the project
       if one is found, any options defined will be set to the value in the file
    2. use the default values from the global config file for any remaining options

the output file will be named the same as the workspace ID combined with the suffix.
the default suffix is "build.cfg"

for example, if the workspace ID is "SomeAwesomeProject", the output file will be
"SomeAwesomeProject.build.cfg", unless the -o option is used to specify a suffix

Examples:
    $ genlinx cfg -w workspace.apw                                                      # generate CFG for workspace.apw
    $ genlinx cfg -w workspace.apw -s                                                   # build with source
    $ genlinx cfg -w workspace.apw -D                                                   # do not build with debug information
    $ genlinx cfg -a                                                                    # search for and automatically select all workspace files
    $ genlinx cfg -A                                                                    # search for and prompt to select workspace files
    $ genlinx cfg -w workspace.apw -i \\path\\to\\includes \\path\\to\\more\\includes          # add additional include paths
    $ genlinx cfg -w workspace.apw -m \\path\\to\\modules \\path\\to\\more\\modules            # add additional module paths
    $ genlinx cfg -w workspace.apw -l \\path\\to\\libraries \\path\\to\\more\\libraries        # add additional library paths`,
        );

    return command;
}

export default cfg;
