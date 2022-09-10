import { Command } from "commander";
import { actions } from "../../actions";

export function build() {
    const command = new Command("build");

    command
        .description("build a NetLinx project from a CFG file")
        .argument("cfg file <string>", "cfg file to build from")
        .action((cfg, options) => actions.build.build(cfg, options));

    return command;
}

export default build;
