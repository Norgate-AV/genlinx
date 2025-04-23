import chalk from "chalk";
import { execa } from "execa";
import { NLRC } from "../../compiler/index.js";
import {
    getFilesByExtension,
    getAppConfig,
    printFiles,
    selectFiles,
    shouldPromptUser,
} from "../../utils/index.js";
import {
    BuildConfig,
    BuildCliArgs,
    BuildLog,
    BuildLogCollection,
    BuildLogCollectionKey,
    ShellCommand,
    BuildOptions,
} from "../../@types/index.js";

export function getBuildLogs(data: string): BuildLogCollection {
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
        error: [...new Set<string>(logs.error)],
        warning: [...new Set<string>(logs.warning)],
    };
}

export function catchAllErrors(errors: Array<string>) {
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

export function printAllWarnings(warnings: Array<string>) {
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

export async function runBuildProcess(
    command: ShellCommand,
    options: BuildConfig,
): Promise<string> {
    const { shell } = options;

    const childProcess = execa(command.path, [...command.args], {
        shell: shell.path,
        windowsVerbatimArguments: true,
        reject: false,
    });

    childProcess.stdout?.pipe(process.stdout);

    const { stdout } = await childProcess;
    return stdout;
}

export async function buildFile(
    file: string,
    command: ShellCommand,
    options: BuildOptions,
): Promise<string> {
    options.verbose &&
        console.log(chalk.blue(`Executing build for ${file}...`));

    return await runBuildProcess(command, options);
}

export async function executeSourceBuild(
    files: Array<string>,
    options: BuildOptions,
): Promise<void> {
    for (const file of files) {
        const command = NLRC.getSourceBuildCommand(file, options);

        const buildResult = await buildFile(file, command, options);

        const logs = getBuildLogs(buildResult);

        printAllWarnings(logs.warning);
        catchAllErrors(logs.error);
    }
}

export async function executeCfgBuild(
    files: Array<string>,
    options: BuildOptions,
): Promise<void> {
    const { verbose } = options;

    if (files.length === 0) {
        verbose && console.log(chalk.blue("Searching for CFG files..."));

        const located = await getFilesByExtension(".cfg");

        if (located.length) {
            verbose && printFiles(located);
            files.push(...located);
        }
    }

    if (files.length === 0) {
        console.log(chalk.red("No CFG files found."));
        process.exit();
    }

    if (shouldPromptUser(options, files)) {
        const selected = await selectFiles(files);

        files.splice(0, files.length);
        files.push(...selected);
    }

    for (const file of files) {
        const command = NLRC.getCfgBuildCommand(file, options);

        const buildResult = await buildFile(file, command, options);

        const logs = getBuildLogs(buildResult);

        printAllWarnings(logs.warning);
        catchAllErrors(logs.error);
    }
}

export const build = {
    async execute(args: BuildCliArgs): Promise<void> {
        try {
            const config = await getAppConfig({
                build: {
                    ...args,
                },
            });

            const { cfgFiles, sourceFiles } = config.build as BuildOptions;

            if (sourceFiles && sourceFiles.length > 0) {
                await executeSourceBuild(
                    sourceFiles,
                    config.build as BuildOptions,
                );

                process.exit();
            }

            if (cfgFiles && cfgFiles.length > 0) {
                await executeCfgBuild(cfgFiles, config.build as BuildOptions);
                process.exit();
            }

            console.log("No source or CFG files specified.");
        } catch (error: any) {
            process.exit(1);
        }
    },
};

export default build;
