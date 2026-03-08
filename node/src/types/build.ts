export type BuildOptions = BuildConfig & BuildCliArgs;

export type BuildCliArgs = {
    cfgFiles?: Array<string>;
    sourceFiles?: Array<string>;
    includePath?: Array<string>;
    modulePath?: Array<string>;
    libraryPath?: Array<string>;
    all?: boolean;
    noAll?: boolean;
    verbose: boolean;
};

export type BuildConfig = {
    nlrc: {
        path: string;
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

export type BuildLog = {
    log: string;
    level: string;
};

export type BuildLogCollection = {
    error: Array<string>;
    warning: Array<string>;
};

export type BuildLogCollectionKey = keyof BuildLogCollection;
