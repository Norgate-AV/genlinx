import { Command, /*Argument,*/ Option } from "commander";
import { actions } from "../actions/index.js";

export function config(): Command {
    const command = new Command("config");

    command
        .description("view/edit configuration properties for genlinx")
        // .addArgument(
        //     new Argument("[key]", "the configuration property to edit"),
        // )
        // .addArgument(
        //     new Argument(
        //         "[value...]",
        //         "the value to set the configuration property to",
        //     ),
        // )
        // .addOption(
        //     new Option("--global", "edit the global configuration").conflicts(
        //         "local",
        //     ),
        // )
        // .addOption(
        //     new Option("--local", "edit the local configuration").conflicts(
        //         "global",
        //     ),
        // )
        .addOption(
            new Option(
                "-l, --list",
                "display the current configuration in stdout",
            ).conflicts("edit"),
        )
        // .addOption(new Option("--add", "add value(s) to an array"))
        // .addOption(new Option("--remove", "remove value(s) from an array"))
        .addOption(
            new Option(
                "-e, --edit",
                "edit the configuration with default text editor",
            ).conflicts("list"),
        )
        .action(async (options) => {
            await actions.config.process(options);
        });

    return command;
}

export default config;
