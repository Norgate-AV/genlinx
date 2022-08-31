import execa from "execa";
import { NLRC, Options } from "../../../lib";
import { getAppConfig } from "../../../lib/utils";

export const build = {
    async build(filePath, cliOptions) {
        try {
            const config = getAppConfig();

            const options = Options.getBuildOptions(cliOptions, config.build);

            const command = NLRC.getCfgBuildCommand(filePath, options);

            const { shell } = options;
            execa(command.path, [...command.args], {
                shell: shell.path,
                windowsVerbatimArguments: true,
            }).stdout.pipe(process.stdout);
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default build;
