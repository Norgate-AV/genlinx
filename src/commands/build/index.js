import execa from "execa";
import { NLRC, Options } from "../../../lib";
import { getAppConfig } from "../../../lib/utils";

export const build = {
    async build(filePath, cliOptions) {
        try {
            const config = getAppConfig();

            const options = Options.getBuildOptions(cliOptions, config.build);

            const command = NLRC.getCfgBuildCommand(filePath, options);

            const { executable } = options;
            await execa(
                executable.path,
                [...executable.args, command.path, ...command.args],
                {
                    cwd: process.cwd(),
                },
            );
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default build;
