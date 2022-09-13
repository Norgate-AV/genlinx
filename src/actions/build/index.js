import execa from "execa";
import chalk from "chalk";
import { NLRC } from "../../../lib";
import {
    getGlobalAppConfig,
    getLocalAppConfig,
    getOptions,
} from "../../../lib/utils";

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
    const pattern = /(?<fullError>ERROR: (?<message>.+))/gm;
    const matches = data.matchAll(pattern);

    const errors = [];

    for (const match of matches) {
        const { fullError, message } = match.groups;
        errors.push({ fullError, message });
    }

    return errors;
}

function catchErrors(data) {
    if (!/ERROR:/gm.test(data)) {
        return;
    }

    const errorCount = getErrorCount(data);
    console.log(chalk.red(`A total of ${errorCount} error(s) occurred.`));

    const errors = getErrors(data);
    for (const error of errors) {
        console.log(chalk.red(error.fullError));
    }

    throw new Error(
        `The build process failed with a total of ${errorCount} error(s).`,
    );
}

export const build = {
    async execute(filePath, cliOptions) {
        try {
            const globalConfig = getGlobalAppConfig();
            const localConfig = getLocalAppConfig(filePath);

            const options = getOptions(
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
            catchErrors(stdout);
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default build;
