import { Command, Argument, Option } from "commander";
import { actions } from "../../actions";

export function config() {
    const command = new Command("config");

    command
        .description("edit configuration properties for genlinx")
        .addArgument(
            new Argument("[key]", "the configuration property to edit"),
        )
        .addArgument(
            new Argument(
                "[value...]",
                "the value to set the configuration property to",
            ),
        )
        .addOption(
            new Option("--global", "edit the global configuration").conflicts(
                "local",
            ),
        )
        .addOption(
            new Option("--local", "edit the local configuration").conflicts(
                "global",
            ),
        )
        .addOption(
            new Option("-l, --list", "display the configuration").conflicts(
                "edit",
            ),
        )
        .addOption(new Option("--add", "add value(s) to an array"))
        .addOption(new Option("--remove", "remove value(s) from an array"))
        .addOption(
            new Option(
                "-e, --edit",
                "edit the configuration with default text editor",
            ).conflicts("list"),
        )
        .action(async (key, value, options) => {
            await actions.config.process(key, value, options);
        });

    return command;
}

export default config;
