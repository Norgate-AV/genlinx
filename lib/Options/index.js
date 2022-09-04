import path from "path";

export class Options {
    static getCfgOptions(apw, cliOptions, localConfig, globalConfig) {
        const includePath = localConfig.includePath
            ? [...globalConfig.includePath, ...localConfig.includePath]
            : globalConfig.includePath;

        const modulePath = localConfig.modulePath
            ? [...globalConfig.modulePath, ...localConfig.modulePath]
            : globalConfig.modulePath;

        const libraryPath = localConfig.libraryPath
            ? [...globalConfig.libraryPath, ...localConfig.libraryPath]
            : globalConfig.libraryPath;

        return {
            rootDirectory: path.resolve(cliOptions.rootDirectory),
            outputFile:
                cliOptions.outputFile ||
                `${apw.id}.${localConfig.outputFile}` ||
                `${apw.id}.${globalConfig.outputFile}`,

            outputLogFile:
                cliOptions.outputLogFile ||
                `${apw.id}.${localConfig.outputLogFile}` ||
                `${apw.id}.${globalConfig.outputLogFile}`,

            outputLogFileOption:
                cliOptions.outputLogFileOption ||
                localConfig.outputLogFileOption ||
                globalConfig.outputLogFileOption,

            outputLogConsoleOption:
                cliOptions.outputLogConsoleOption ||
                localConfig.outputLogConsoleOption ||
                globalConfig.outputLogConsoleOption,

            buildWithDebugInformation:
                cliOptions.buildWithDebugInformation ||
                localConfig.buildWithDebugInformation ||
                globalConfig.buildWithDebugInformation,

            buildWithSource:
                cliOptions.buildWithSource ||
                localConfig.buildWithSource ||
                globalConfig.buildWithSource,

            includePath: cliOptions.includePath
                ? [...includePath, ...cliOptions.includePath]
                : includePath,

            modulePath: cliOptions.modulePath
                ? [...modulePath, ...cliOptions.modulePath]
                : modulePath,

            libraryPath: cliOptions.libraryPath
                ? [...libraryPath, ...cliOptions.libraryPath]
                : libraryPath,
        };
    }

    static getArchiveOptions(apw, cliOptions, localConfig, globalConfig) {
        const extraFileSearchLocations = localConfig.extraFileSearchLocations
            ? [
                  ...localConfig.extraFileSearchLocations,
                  ...globalConfig.extraFileSearchLocations,
                  path.dirname(apw.filePath),
              ]
            : [
                  ...globalConfig.extraFileSearchLocations,
                  path.dirname(apw.filePath),
              ];

        return {
            ...cliOptions,
            outputFile:
                cliOptions.outputFile ||
                `${apw.id}.${localConfig.outputFile}` ||
                `${apw.id}.${globalConfig.outputFile}`,

            includeCompiledSourceFiles:
                cliOptions.includeCompiledSourceFiles ||
                localConfig.includeCompiledSourceFiles ||
                globalConfig.includeCompiledSourceFiles,

            includeCompiledModuleFiles:
                cliOptions.includeCompiledModuleFiles ||
                localConfig.includeCompiledModuleFiles ||
                globalConfig.includeCompiledModuleFiles,

            includeFilesNotInWorkspace:
                cliOptions.includeFilesNotInWorkspace ||
                localConfig.includeFilesNotInWorkspace ||
                globalConfig.includeFilesNotInWorkspace,

            extraFileSearchLocations: cliOptions.extraFileSearchLocations
                ? [
                      ...extraFileSearchLocations,
                      ...cliOptions.extraFileSearchLocations,
                  ]
                : extraFileSearchLocations,

            extraFileArchiveLocation:
                cliOptions.extraFileArchiveLocation ||
                localConfig.extraFileArchiveLocation ||
                globalConfig.extraFileArchiveLocation,
        };
    }

    static getBuildOptions(cliOptions, localConfig, globalConfig) {
        return {
            ...cliOptions,
            nlrc: localConfig.nlrc || globalConfig.nlrc,
            shell: localConfig.shell || globalConfig.shell,
        };
    }
}

export default Options;
