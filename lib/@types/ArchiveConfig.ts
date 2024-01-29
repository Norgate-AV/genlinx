import { CliArchiveOptions } from "./index.js";

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

export type ArchiveOptions = ArchiveConfig & CliArchiveOptions;
