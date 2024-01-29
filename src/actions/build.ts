import util from "util";
import chalk from "chalk";
import { NLRC } from "../../lib/index.js";
import {
    getFilesByExtension,
    getAppConfig,
    printFiles,
    selectFiles,
    shouldPromptUser,
} from "../../lib/utils/index.js";
import {
    BuildConfig,
    Config,
    CliBuildOptions,
    BuildLog,
    BuildLogCollection,
    BuildLogCollectionKey,
    ShellCommand,
} from "../../lib/@types/index.js";

function getBuildLogs(data: string): BuildLogCollection {
    const pattern = /(?<log>(?<level>ERROR|WARNING): .+)/gm;
    const matches = data.matchAll(pattern);

    const logs: BuildLogCollection = {
        error: [],
        warning: [],
    };

    for (const match of matches) {
        if (!match.groups) {
            continue;
        }

        const { log, level } = match.groups as BuildLog;
        logs[level.toLowerCase() as BuildLogCollectionKey].push(log);
    }

    return {
        error: [...new Set(logs.error)],
        warning: [...new Set(logs.warning)],
    };
}

function catchAllErrors(errors: Array<string>) {
    if (errors.length === 0) {
        return;
    }

    console.log(chalk.red(`A total of ${errors.length} error(s) occurred.`));

    for (const error of errors) {
        console.log(chalk.red(error));
    }

    throw new Error(
        `The build process failed with a total of ${errors.length} error(s).`,
    );
}

function printAllWarnings(warnings: Array<string>) {
    if (warnings.length === 0) {
        return;
    }

    console.log(
        chalk.yellow(`A total of ${warnings.length} warning(s) occurred.`),
    );

    for (const warning of warnings) {
        console.log(chalk.yellow(warning));
    }
}

async function runBuildProcess(
    command: ShellCommand,
    options: BuildConfig,
): Promise<string> {
    const { shell } = options;

    const { execa } = await import("execa");
    const childProcess = execa(command.path, [...command.args], {
        shell: shell.path,
        windowsVerbatimArguments: true,
        reject: false,
    });

    childProcess.stdout?.pipe(process.stdout);

    const { stdout } = await childProcess;
    return stdout;
}

async function buildFile(
    file: string,
    command: ShellCommand,
    options: BuildConfig,
): Promise<string> {
    console.log(chalk.blue(`Executing build for ${file}...`));

    const buildResult = await runBuildProcess(command, options);

    return buildResult;
}

async function executeSourceBuild(
    sourceFile: string,
    config: Config,
): Promise<void> {
    const command = NLRC.getSourceBuildCommand(sourceFile, config.build);

    const buildResult = await buildFile(sourceFile, command, config.build);

    const logs = getBuildLogs(buildResult);

    printAllWarnings(logs.warning);
    catchAllErrors(logs.error);
}

async function executeCfgBuild(
    cfgFiles: Array<string>,
    config: Config,
): Promise<void> {
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

    if (shouldPromptUser(config.build, cfgFiles)) {
        const selectedCfgFiles = await selectFiles(cfgFiles);

        cfgFiles.splice(0, cfgFiles.length);
        cfgFiles.push(...selectedCfgFiles);
    }

    for (const cfgFile of cfgFiles) {
        const command = NLRC.getCfgBuildCommand(cfgFile, config.build);

        const buildResult = await buildFile(cfgFile, command, config.build);

        const logs = getBuildLogs(buildResult);

        printAllWarnings(logs.warning);
        catchAllErrors(logs.error);
    }
}

export const build = {
    async execute(cliOptions: CliBuildOptions): Promise<void> {
        try {
            const config = await getAppConfig({
                build: {
                    ...cliOptions,
                },
            });

            const { cfgFiles, sourceFile } = config.build;

            if (sourceFile) {
                await executeSourceBuild(sourceFile, config);
                process.exit();
            }

            if (cfgFiles.length > 0) {
                await executeCfgBuild(cfgFiles, config);
                process.exit();
            }

            console.log(chalk.red("No source or CFG files specified."));
        } catch (error: any) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default build;
