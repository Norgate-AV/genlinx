export type CfgOptions = CfgConfig & CfgCliArgs;

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
    verbose: boolean;
};

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
