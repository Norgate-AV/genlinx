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
