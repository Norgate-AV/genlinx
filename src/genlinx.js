import figlet from "figlet";
import StringBuilder from "string-builder";
import { Command } from "commander";
import { version } from "../package.json";
import { archive, build, cfg, config } from "./commands";

export function cli(args) {
    const program = new Command();

    program
        .name("genlinx")
        .description("cli helper utility for NetLinx projects 🚀🚀🚀")
        .version(version, "-v, --version");

    program
        .addCommand(archive())
        .addCommand(build())
        .addCommand(cfg())
        .addCommand(config());

    program.addHelpText("beforeAll", () => {
        const builder = new StringBuilder();

        builder
            .appendLine(figlet.textSync("genlinx"))
            .appendLine()
            .appendLine("Open source CLI tool for NetLinx projects")
            .appendLine("Copyright (c) 2022, Norgate AV Solutions Ltd")
            .appendLine("https://github.com/Norgate-AV-Solutions-Ltd/genlinx")
            .appendLine()
            .appendLine("===================================================")
            .appendLine();

        return builder.toString();
    });

    if (!args[2]) {
        program.outputHelp();
        process.exit();
    }

    program.parse(args);
}

export default cli;
