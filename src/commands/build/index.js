import execa from "execa";
import chalk from "chalk";
import { NLRC, Options } from "../../../lib";
import { getGlobalAppConfig, getLocalAppConfig } from "../../../lib/utils";

export const build = {
    async build(filePath, cliOptions) {
        try {
            const globalConfig = getGlobalAppConfig();
            const localConfig = getLocalAppConfig(filePath);

            const options = Options.getBuildOptions(
                cliOptions,
                localConfig.build,
                globalConfig.build,
            );

            const command = NLRC.getCfgBuildCommand(filePath, options);

            const { shell } = options;
            const childProcess = execa(command.path, [...command.args], {
                shell: shell.path,
                windowsVerbatimArguments: true,
            });

            childProcess.stdout.pipe(process.stdout);

            const { stdout } = await childProcess;
            if (/ERROR:/gm.test(stdout)) {
                const errorPattern = /(ERROR: (?<message>.+))/gm;
                const errorMatches = stdout.matchAll(errorPattern);

                for (const match of errorMatches) {
                    console.log(chalk.red(match.groups.message));
                }

                const errorCountPattern = /(?<count>[1-9]+) error\(s\)/gm;
                const errorCountMatches = stdout.matchAll(errorCountPattern);

                let errorCount = 0;

                for (const match of errorCountMatches) {
                    errorCount += Number(match.groups.count);
                }

                console.log(chalk.red(`${errorCount} error(s)`));
                throw new Error(
                    `The build process failed with a total of ${errorCount} error(s).`,
                );
            }
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default build;
