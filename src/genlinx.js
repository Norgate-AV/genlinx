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
        .description("Generates a Netlinx build CFG file")
        .argument(
            "APW File <string>",
            "The APW file to generate the build CFG for",
        )
        .option("-d, --dir <string>", "Main AXS Root Directory Reference", ".")
        .option("-o, --output <string>", "The output file", "build.cfg")
        .option("-lf, --log-file <string>", "Output log file name", "")
        .addOption(
            new Option(
                "-lo, --log-file-option <string>",
                "Output log file option",
            )
                .choices(["A", "N"])
                .default("N"),
        )
        .option("-lc, --log-console-option", "Output log to the console", false)
        .option(
            "-bd, --build-with-debug",
            "Build with debug information",
            false,
        )
        .option("-bs, --build-with-source", "Build with source", false)
        .option(
            "-i, --include-path <string...>",
            "Additional include paths",
            "",
        )
        .option("-m, --module-path <string...>", "Additional module paths", "")
        .option(
            "-l, --library-path <string...>",
            "Additional library paths",
            "",
        )
        .action((apw, options) => cfg.create(apw, options));

    program
        .command("archive")
        .description("Generates a Netlinx workspace archive")
        .argument(
            "APW File <string>",
            "The APW file to generate the archive for",
        )
        .option(
            "-d, --dir <string>",
            "The root directory the archive should use",
            ".",
        )
        .option("-o, --output <string>", "The output file", "project.zip")
        .action((apw, options) => archive.create(apw, options));

    program
        .command("build")
        .description("Builds a Netlinx project from a CFG file")
        .argument(
            "Netlinx build CFG File <string>",
            "The CFG file to build from",
        )
        .action((cfg, options) => build.build(cfg, options));

    program.parse(args);

    if (!args[2]) {
        program.outputHelp();
    }
}

export default cli;
