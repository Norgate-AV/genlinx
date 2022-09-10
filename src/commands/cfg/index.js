import { Command, Option } from "commander";
import { actions } from "../../actions";

export function cfg() {
    const command = new Command("cfg");

    command
        .description("generate NetLinx build CFG files")
        .argument("apw file <string>", "apw file to generate the CFG from")
        .option(
            "-d, --root-directory <string>",
            "root directory reference",
            ".",
        )
        .option("-o, --output-file <string>", "output file name")
        .option("-L, --output-log-file <string>", "output log file name")
        .addOption(
            new Option(
                "-O, --output-log-file-option <string>",
                "output log file option",
            ).choices(["A", "N"]),
        )
        .option("-C, --output-log-console-option", "output log to console")
        .option(
            "-D, --build-with-debug-information",
            "build with debug information",
        )
        .option("-S, --build-with-source", "build with source")
        .option("-i, --include-path <string...>", "additional include paths")
        .option("-m, --module-path <string...>", "additional module paths")
        .option("-l, --library-path <string...>", "additional library paths")
        .action((apw, options) => actions.cfg.create(apw, options));

    return command;
}

export default cfg;
