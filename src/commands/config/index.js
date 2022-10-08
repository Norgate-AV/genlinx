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
        .addOption(new Option("-d, --delete", "delete the given key"))
        .action(async (key, value, options) => {
            if (!key && !value) {
                if (options.list) {
                    actions.config.show(options);
                    process.exit();
                }

                if (options.edit) {
                    await actions.config.edit(options);
                    process.exit();
                }

                process.exit();
            }

            if (!value) {
                actions.config.get(key, options);
                process.exit();
            }

            await actions.config.set(key, value, options);
        });

    return command;
}

export default config;
