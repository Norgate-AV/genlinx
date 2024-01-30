export type CfgCliArgs = {
    workspaceFiles: Array<string>;
    rootDirectory: string;
    outputFileSuffix?: string;
    outputLogFileSuffix?: string;
    outputLogFileOption?: "A" | "N";
    outputLogConsoleOption?: boolean;
    noOutputLogConsoleOption?: boolean;
    buildWithDebugInformation?: boolean;
    noBuildWithDebugInformation?: boolean;
    buildWithSource?: boolean;
    noBuildWithSource?: boolean;
    includePath?: Array<string>;
    modulePath?: Array<string>;
    libraryPath?: Array<string>;
    all?: boolean;
    noAll?: boolean;
};
