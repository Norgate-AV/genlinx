import path from "path";

export class Options {
    static #getCfgLocalConfig(config) {
        return {
            outputFile: "outputFile" in config ? config.outputFile : null,
            outputLogFile:
                "outputLogFile" in config ? config.outputLogFile : null,
            outputLogFileOption:
                "outputLogFileOption" in config
                    ? config.outputLogFileOption
                    : null,
            outputLogConsoleOption:
                "outputLogConsoleOption" in config
                    ? config.outputLogConsoleOption
                    : null,
            buildWithDebugInformation:
                "buildWithDebugInformation" in config
                    ? config.buildWithDebugInformation
                    : null,
            buildWithSource:
                "buildWithSource" in config ? config.buildWithSource : null,
            includePath: "includePath" in config ? config.includePath : null,
            modulePath: "modulePath" in config ? config.modulePath : null,
            libraryPath: "libraryPath" in config ? config.libraryPath : null,
        };
    }

    static getCfgOptions(apw, cliOptions, localConfig, globalConfig) {
        localConfig = this.#getCfgLocalConfig(localConfig);

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
                cliOptions.outputFile || localConfig.outputFile
                    ? `${apw.id}.${localConfig.outputFile}`
                    : null || `${apw.id}.${globalConfig.outputFile}`,

            outputLogFile:
                cliOptions.outputLogFile || localConfig.outputLogFile
                    ? `${apw.id}.${localConfig.outputLogFile}`
                    : null || `${apw.id}.${globalConfig.outputLogFile}`,

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

    static #getArchiveLocalConfig(config) {
        return {
            outputFile: "outputFile" in config ? config.outputFile : null,
            includeCompiledSourceFiles:
                "includeCompiledSourceFiles" in config
                    ? config.includeCompiledSourceFiles
                    : null,
            includeCompiledModuleFiles:
                "includeCompiledModuleFiles" in config
                    ? config.includeCompiledModuleFiles
                    : null,
            includeFilesNotInWorkspace:
                "includeFilesNotInWorkspace" in config
                    ? config.includeFilesNotInWorkspace
                    : null,
            extraFileSearchLocations:
                "extraFileSearchLocations" in config
                    ? config.extraFileSearchLocations
                    : null,
            extraFileArchiveLocation:
                "extraFileArchiveLocation" in config
                    ? config.extraFileArchiveLocation
                    : null,
        };
    }

    static getArchiveOptions(apw, cliOptions, localConfig, globalConfig) {
        localConfig = this.#getArchiveLocalConfig(localConfig);

        const extraFileSearchLocations = localConfig.extraFileSearchLocations
            ? [
                  ...globalConfig.extraFileSearchLocations,
                  ...localConfig.extraFileSearchLocations,
                  path.dirname(apw.filePath),
              ]
            : [
                  ...globalConfig.extraFileSearchLocations,
                  path.dirname(apw.filePath),
              ];

        return {
            outputFile:
                cliOptions.outputFile || localConfig.outputFile
                    ? `${apw.id}.${localConfig.outputFile}`
                    : null || `${apw.id}.${globalConfig.outputFile}`,

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

    static #getBuildLocalConfig(config) {
        return {
            nlrc: "nlrc" in config ? config.nlrc : null,
            shell: "shell" in config ? config.shell : null,
        };
    }

    static getBuildOptions(cliOptions, localConfig, globalConfig) {
        localConfig = this.#getBuildLocalConfig(localConfig);

        return {
            ...cliOptions,
            nlrc: localConfig.nlrc || globalConfig.nlrc,
            shell: localConfig.shell || globalConfig.shell,
        };
    }
}

export default Options;
