import { Command } from "commander";
import { actions } from "../../actions";

export function archive() {
    const command = new Command("archive");

    command
        .description("generate a NetLinx workspace zip archive")
        .argument("apw file <string>", "apw file to generate the archive from")
        .option("-o, --output-file <string>", "output file name")
        .option(
            "-S, --include-compiled-source-files",
            "include compiled source files",
        )
        .option(
            "-M, --include-compiled-module-files",
            "include compiled module files",
        )
        .option(
            "-N, --include-files-not-in-workspace",
            "include files not in workspace",
        )
        .option(
            "-l, --extra-file-search-locations <string...>",
            "extra file locations to search",
        )
        .option(
            "-p, --extra-file-archive-location <string>",
            "location to place extra files in the archive",
        )
        .action((apw, options) => actions.archive.create(apw, options));

    return command;
}

export default archive;
