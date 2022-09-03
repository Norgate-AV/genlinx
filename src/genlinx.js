import { program, Option } from "commander";
import { version } from "../package.json";
import { archive, build, cfg } from "./commands";

export function cli(args) {
    program
        .name("genlinx")
        .description("cli helper utility for NetLinx projects 🚀🚀🚀")
        .version(version, "-v, --version");

    program
        .command("cfg")
        .description("generate a NetLinx build CFG file")
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
        .action((apw, options) => cfg.create(apw, options));

    program
        .command("archive")
        .description("generate a NetLinx workspace zip archive")
        .argument("apw file <string>", "apw file to generate the archive from")
        .option("-o, --output-file <string>", "output file name")
        .option(
            "-S, --include-compiled-source-files",
            "include compiled source files",
        )
        .option(
            "-M, --include-compiled-module-files",
            "include compiled module files",
        )
        .option(
            "-N, --include-files-not-in-workspace",
            "include files not in workspace",
        )
        .option(
            "-l, --extra-file-search-locations <string...>",
            "extra file locations to search",
        )
        .action((apw, options) => archive.create(apw, options));

    program
        .command("build")
        .description("build a NetLinx project from a CFG file")
        .argument("cfg file <string>", "cfg file to build from")
        .action((cfg, options) => build.build(cfg, options));

    program.parse(args);

    if (!args[2]) {
        program.outputHelp();
    }
}

export default cli;
