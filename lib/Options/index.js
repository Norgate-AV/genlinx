import path from "path";

export class Options {
    static getCfgOptions(cliOptions, config) {
        return {
            rootDirectory: path.resolve(cliOptions.rootDirectory),

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
}

export default Options;
