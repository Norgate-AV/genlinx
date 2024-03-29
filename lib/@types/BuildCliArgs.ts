export type BuildCliArgs = {
    cfgFiles: Array<string>;
    sourceFiles?: Array<string>;
    includePath?: Array<string>;
    modulePath?: Array<string>;
    libraryPath?: Array<string>;
    all?: boolean;
    noAll?: boolean;
    verbose: boolean;
};
