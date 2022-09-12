import { Command, Option } from "commander";
import { actions } from "../../actions";

export function cfg() {
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
                "-d, --root-directory <string>",
                "root directory reference",
            ).default(".", "use current directory as root"),
        )
        .addOption(
            new Option(
                "-o, --output-file-suffix <string>",
                "output file suffix",
            ),
        )
        .addOption(
            new Option(
                "-L, --output-log-file-suffix <string>",
                "output log file suffix",
            ),
        )
        .addOption(
            new Option(
                "-O, --output-log-file-option <string>",
                "output log file option",
            ).choices(["A", "N"]),
        )
        .addOption(
            new Option(
                "-C, --output-log-console-option",
                "output log to console",
            ),
        )
        .addOption(
            new Option(
                "-D, --build-with-debug-information",
                "build with debug information",
            ),
        )
        .addOption(new Option("-S, --build-with-source", "build with source"))
        .addOption(
            new Option(
                "-i, --include-path <string...>",
                "additional include paths",
            ),
        )
        .addOption(
            new Option(
                "-m, --module-path <string...>",
                "additional module paths",
            ),
        )
        .addOption(
            new Option(
                "-l, --library-path <string...>",
                "additional library paths",
            ),
        )
        .addOption(
            new Option(
                "-a, --all",
                "select all workspace files in current directory without prompting",
            ),
        )
        .action((apw, options) => actions.cfg.create(apw, options));

    return command;
}

export default cfg;
