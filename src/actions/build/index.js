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
    const pattern =
        /(?<log>ERROR: (?<path>.+)\((?<line>\d+)\): (?<code>.+): (?<message>.+))/gm;
    const matches = data.matchAll(pattern);

    const errors = [];

    for (const match of matches) {
        const { log, path, line, code, message } = match.groups;
        errors.push({ log, path, line, code, message });
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
        console.log(chalk.red(error.log));
    }

    throw new Error(
        `The build process failed with a total of ${errorCount} error(s).`,
    );
}

function getWarningCount(data) {
    const pattern = /(?<count>[1-9]+) warning\(s\)/gm;
    const matches = data.matchAll(pattern);

    let count = 0;

    for (const match of matches) {
        count += Number(match.groups.count);
    }

    return count;
}

function getWarnings(data) {
    const pattern =
        /(?<log>WARNING: (?<path>.+)\((?<line>\d+)\)(?:\w+)?: (?:'(?<warning>.+)'$|((?<code>\w+\d+): (?<message>.+)$)))/gm;
    const matches = data.matchAll(pattern);

    const warnings = [];

    for (const match of matches) {
        const { log, path, line, warning, code, message } = match.groups;
        warnings.push({ log, path, line, warning, code, message });
    }

    return warnings;
}

function printWarnings(data) {
    if (!/WARNING:/gm.test(data)) {
        return;
    }

    const warningCount = getWarningCount(data);
    console.log(
        chalk.yellow(`A total of ${warningCount} warning(s) occurred.`),
    );

    const warnings = getWarnings(data);
    for (const warning of warnings) {
        console.log(chalk.yellow(warning.log));
    }
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
            printWarnings(stdout);
            catchErrors(stdout);
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default build;
