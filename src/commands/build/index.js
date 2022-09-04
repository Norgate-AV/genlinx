import execa from "execa";
import { NLRC, Options } from "../../../lib";
import { getGlobalAppConfig, getLocalAppConfig } from "../../../lib/utils";

export const build = {
    async build(filePath, cliOptions) {
        try {
            const globalConfig = getGlobalAppConfig();
            const localConfig = getLocalAppConfig(filePath);

            const options = Options.getBuildOptions(
                cliOptions,
                localConfig.build ? localConfig.build : {},
                globalConfig.build,
            );

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
