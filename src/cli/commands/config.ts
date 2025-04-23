import { Argument, Command, Option } from "commander";
import { actions } from "../actions/index.js";

export function config(): Command {
    const command = new Command("config");

    command
        .description("view/edit configuration properties for genlinx")
        .addArgument(new Argument("[key]", "configuration key"))
        .addOption(
            new Option("--global", "use global configuration").conflicts(
                "local",
            ),
        )
        .addOption(
            new Option("--local", "use local configuration").conflicts(
                "global",
            ),
        )
        .addOption(
            new Option(
                "-l, --list",
                "display the configuration in stdout",
            ).conflicts("edit"),
        )
        .addOption(
            new Option(
                "-e, --edit",
                "edit the configuration with default text editor",
            ).conflicts("list"),
        )
        .action(async (key, options) => {
            await actions.config.process(key, options);
        });

    return command;
}

export default config;
