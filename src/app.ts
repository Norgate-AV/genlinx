#!/usr/bin/env node

import figlet from "figlet";
import StringBuilder from "string-builder";
import { Command } from "commander";
import { version } from "../package.json";
import { archive, build, cfg } from "./commands";

const args = process.argv;
const program = new Command();

program
    .name("genlinx")
    .description("cli helper utility for NetLinx projects 🚀🚀🚀")
    .version(version, "-v, --version");

program.addCommand(archive()).addCommand(build()).addCommand(cfg());

program.addHelpText("beforeAll", () => {
    const builder = new StringBuilder();

    builder
        .appendLine(figlet.textSync("genlinx"))
        .appendLine()
        .appendLine("Open source CLI tool for NetLinx projects")
        .appendLine("Copyright (c) 2023, Norgate AV Services Limited")
        .appendLine("https://github.com/Norgate-AV/genlinx")
        .appendLine()
        .appendLine("===================================================")
        .appendLine();

    return builder.toString();
});

program.addHelpText("afterAll", () => {
    const builder = new StringBuilder();

    builder
        .appendLine("===================================================")
        .appendLine()
        .appendLine("For more help, make sure to check out the man page:")
        .appendLine("    $ man genlinx");

    return builder.toString();
});

if (!args[2]) {
    program.outputHelp();
    process.exit();
}

program.parse(args);