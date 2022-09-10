import execa from "execa";
import chalk from "chalk";
import { NLRC, Options } from "../../../lib";
import { getGlobalAppConfig, getLocalAppConfig } from "../../../lib/utils";

function getErrorCount(data) {
    const pattern = /(?<count>[1-9]+) error\(s\)/gm;
    const matches = data.matchAll(pattern);

    let count = 0;

    for (const match of matches) {
        count += Number(match.groups.count);
    }

    return count;
}

function getErrors(data) {
    const pattern = /(ERROR: (?<message>.+))/gm;
    const matches = data.matchAll(pattern);

    const errors = [];

    for (const match of matches) {
        errors.push(match.groups.message);
    }

    return errors;
}

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
                reject: false,
            });

            childProcess.stdout.pipe(process.stdout);

            const { stdout } = await childProcess;
            if (/ERROR:/gm.test(stdout)) {
                const errorCount = getErrorCount(stdout);
                console.log(
                    chalk.red(`A total of ${errorCount} error(s) occurred.`),
                );

                const errors = getErrors(stdout);
                for (const error of errors) {
                    console.log(chalk.red(error));
                }

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
