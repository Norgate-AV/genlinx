import execa from "execa";
import chalk from "chalk";
import { NLRC } from "../../../lib";
import {
    getFilesByExtension,
    getGlobalAppConfig,
    getLocalAppConfig,
    getOptions,
    printFiles,
    selectFiles,
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
        /(?<log>ERROR: (?<path>.+)\((?<line>\d+)\)(?:(?<error>[^:].+: .+$)|: (?:(?<code>\w+\d+): (?<message>.+)$)))/gm;
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

function shouldPromptUser(options, files) {
    return !options.all && files.length > 1;
}

async function runBuildProcess(command, options) {
    const { shell } = options;
    const childProcess = execa(command.path, [...command.args], {
        shell: shell.path,
        windowsVerbatimArguments: true,
        reject: false,
    });

    childProcess.stdout.pipe(process.stdout);

    const { stdout } = await childProcess;
    return stdout;
}

export const build = {
    async execute(cliOptions) {
        try {
            const { cfgFiles, sourceFile } = cliOptions;

            const globalConfig = getGlobalAppConfig();

            if (sourceFile) {
                console.log(chalk.blue(`Executing build for ${sourceFile}...`));

                const localConfig = getLocalAppConfig(sourceFile);

                const options = getOptions(
                    cliOptions,
                    localConfig.build,
                    globalConfig.build,
                );

                const command = NLRC.getSourceBuildCommand(sourceFile, options);

                const buildResult = await runBuildProcess(command, options);
                printWarnings(buildResult);
                catchErrors(buildResult);

                process.exit();
            }

            if (cfgFiles.length === 0) {
                console.log(chalk.blue("Searching for CFG files..."));

                const locatedCfgFiles = await getFilesByExtension(".cfg");

                if (locatedCfgFiles.length) {
                    printFiles(locatedCfgFiles);
                    cfgFiles.push(...locatedCfgFiles);
                }
            }

            if (cfgFiles.length === 0) {
                console.log(chalk.red("No CFG files found."));
                process.exit();
            }

            if (shouldPromptUser(cliOptions, cfgFiles)) {
                const selectedCfgFiles = await selectFiles(cfgFiles);

                cfgFiles.splice(0, cfgFiles.length);
                cfgFiles.push(...selectedCfgFiles);
            }

            for (const cfgFile of cfgFiles) {
                console.log(chalk.blue(`Executing build for ${cfgFile}...`));

                const localConfig = getLocalAppConfig(cfgFile);

                const options = getOptions(
                    cliOptions,
                    localConfig.build,
                    globalConfig.build,
                );

                const command = NLRC.getCfgBuildCommand(cfgFile, options);

                const buildResult = await runBuildProcess(command, options);
                printWarnings(buildResult);
                catchErrors(buildResult);
            }
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default build;
