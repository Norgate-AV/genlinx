export type CliArchiveOptions = {
    workspaceFiles: Array<string>;
    outputFileSuffix?: string;
    includeCompiledSourceFiles?: boolean;
    includeCompiledModuleFiles?: boolean;
    includeFilesNotInWorkspace?: boolean;
    extraFileSearchLocations?: Array<string>;
    extraFileArchiveLocation?: string;
};
