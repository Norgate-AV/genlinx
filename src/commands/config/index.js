import { Command } from "commander";
import { actions } from "../../actions";

export function config() {
    const command = new Command("config");

    command
        .description("edit configuration properties for genlinx")
        .option("-g, --global", "edit the global configuration")
        .option("-l, --local", "edit the local configuration");

    command
        .command("set")
        .argument("key <string>", "key to set")
        .argument("value(s) <string...>", "value(s) to set")
        .description("set configuration properties for genlinx")
        .action(() => actions.config.set());

    command
        .command("get")
        .argument("key <string>", "key to get")
        .description("get configuration properties for genlinx")
        .action(() => actions.config.get());

    return command;
}

export default config;
