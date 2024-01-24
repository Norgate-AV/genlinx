export type CfgConfig = {
    outputFileSuffix: string;
    outputLogFileSuffix: string;
    outputLogFileOption: "N" | "Y";
    outputLogConsoleOption: boolean;
    buildWithDebugInformation: boolean;
    buildWithSource: boolean;
    includePath: Array<string>;
    modulePath: Array<string>;
    libraryPath: Array<string>;
    all: boolean;
};

export type ArchiveConfig = {
    outputFileSuffix: string;
    includeCompiledSourceFiles: boolean;
    includeCompiledModuleFiles: boolean;
    includeFilesNotInWorkspace: boolean;
    extraFileSearchLocations: Array<string>;
    extraFileArchiveLocation: string;
    all: boolean;
    ignoredFiles: Array<string>;
};

export type BuildConfig = {
    nlrc: {
        path: string;
        option: {
            cfg: "-CFG";
            includePath: "-I";
            modulePath: "-M";
            libraryPath: "-L";
        };
        includePath: Array<string>;
        modulePath: Array<string>;
        libraryPath: Array<string>;
    };
    shell: {
        path: string;
    };
    all: boolean;
    createCfg: boolean;
};

export type Config = {
    cfg: CfgConfig;
    archive: ArchiveConfig;
    build: BuildConfig;
};

export type GlobalConfig = Partial<Config>;
export type LocalConfig = Partial<Config>;
export type CliConfig = Partial<Config>;
