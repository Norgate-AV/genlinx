import path from "path";

export class Options {
    static getCfgOptions(apw, cliOptions, config) {
        return {
            rootDirectory: path.resolve(cliOptions.rootDirectory),
            outputFile: cliOptions.outputFile || `${apw.id}.build.cfg`,

            outputLogFile: cliOptions.outputLogFile || config.outputLogFile,
            outputLogFileOption:
                cliOptions.outputLogFileOption || config.outputLogFileOption,

            outputLogConsoleOption:
                cliOptions.outputLogConsoleOption ||
                config.outputLogConsoleOption,

            buildWithDebugInformation:
                cliOptions.buildWithDebugInformation ||
                config.buildWithDebugInformation,

            buildWithSource:
                cliOptions.buildWithSource || config.buildWithSource,

            includePath: cliOptions.includePath
                ? [...config.includePath, ...cliOptions.includePath]
                : config.includePath,

            modulePath: cliOptions.modulePath
                ? [...config.modulePath, ...cliOptions.modulePath]
                : config.modulePath,

            libraryPath: cliOptions.libraryPath
                ? [...config.libraryPath, ...cliOptions.libraryPath]
                : config.libraryPath,
        };
    }

    static getBuildOptions(cliOptions, config) {
        return {
            ...cliOptions,
            nlrc: config.nlrc,
            executable: config.executable,
        };
    }
}

export default Options;
