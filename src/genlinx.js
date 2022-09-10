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

    program.parse(args);

    if (!args[2]) {
        program.outputHelp();
    }
}

export default cli;
