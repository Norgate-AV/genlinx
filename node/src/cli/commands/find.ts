import { Command, Option } from "commander";
import { actions } from "../actions/index.js";

export function find(): Command {
    const command = new Command("find");

    command
        .description("find NetLinx devices on a local broadcast subnet")
        .addOption(
            new Option(
                "-t, --timeout <number>",
                "timeout in milliseconds",
            ).default(6000),
        )
        .addOption(new Option("-j, --json", "output as JSON"))
        .action(async (options) => {
            await actions.find.execute(options);
        });

    return command;
}

export default find;
