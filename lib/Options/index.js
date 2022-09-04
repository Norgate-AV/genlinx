import path from "path";

export class Options {
    static getCfgOptions(apw, cliOptions, globalConfig) {
        return {
            rootDirectory: path.resolve(cliOptions.rootDirectory),
            outputFile:
                cliOptions.outputFile || `${apw.id}.${globalConfig.outputFile}`,

            outputLogFile:
                cliOptions.outputLogFile ||
                `${apw.id}.${globalConfig.outputLogFile}`,

            outputLogFileOption:
                cliOptions.outputLogFileOption ||
                globalConfig.outputLogFileOption,

            outputLogConsoleOption:
                cliOptions.outputLogConsoleOption ||
                globalConfig.outputLogConsoleOption,

            buildWithDebugInformation:
                cliOptions.buildWithDebugInformation ||
                globalConfig.buildWithDebugInformation,

            buildWithSource:
                cliOptions.buildWithSource || globalConfig.buildWithSource,

            includePath: cliOptions.includePath
                ? [...globalConfig.includePath, ...cliOptions.includePath]
                : globalConfig.includePath,

            modulePath: cliOptions.modulePath
                ? [...globalConfig.modulePath, ...cliOptions.modulePath]
                : globalConfig.modulePath,

            libraryPath: cliOptions.libraryPath
                ? [...globalConfig.libraryPath, ...cliOptions.libraryPath]
                : globalConfig.libraryPath,
        };
    }

    static getArchiveOptions(apw, cliOptions, globalConfig) {
        return {
            ...cliOptions,
            outputFile:
                cliOptions.outputFile || `${apw.id}.${globalConfig.outputFile}`,

            includeCompiledSourceFiles:
                cliOptions.includeCompiledSourceFiles ||
                globalConfig.includeCompiledSourceFiles,

            includeCompiledModuleFiles:
                cliOptions.includeCompiledModuleFiles ||
                globalConfig.includeCompiledModuleFiles,

            includeFilesNotInWorkspace:
                cliOptions.includeFilesNotInWorkspace ||
                globalConfig.includeFilesNotInWorkspace,

            extraFileSearchLocations: cliOptions.extraFileSearchLocations
                ? [
                      ...globalConfig.extraFileSearchLocations,
                      ...cliOptions.extraFileSearchLocations,
                      path.dirname(apw.filePath),
                  ]
                : [
                      ...globalConfig.extraFileSearchLocations,
                      path.dirname(apw.filePath),
                  ],

            extraFileArchiveLocation:
                cliOptions.extraFileArchiveLocation ||
                globalConfig.extraFileArchiveLocation,
        };
    }

    static getBuildOptions(cliOptions, globalConfig) {
        return {
            ...cliOptions,
            nlrc: globalConfig.nlrc,
            shell: globalConfig.shell,
        };
    }
}

export default Options;
