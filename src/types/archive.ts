import { AmxFileType, ObjectTypes } from "./index.js";

export interface ArchiveItem {
    addToArchive(): void;
}

export const ArchiveFileType = {
    APW: AmxFileType,
    Script: "Script",
    Env: "Env",
} as const;

export type ArchiveFileType = ObjectTypes<typeof ArchiveFileType>;

export type ArchiveOptions = ArchiveConfig & ArchiveCliArgs;

export type ArchiveCliArgs = {
    workspaceFiles: Array<string>;
    outputFileSuffix?: string;
    includeCompiledSourceFiles?: boolean;
    includeCompiledModuleFiles?: boolean;
    includeFilesNotInWorkspace?: boolean;
    extraFileSearchLocations?: Array<string>;
    extraFileArchiveLocation?: string;
    verbose: boolean;
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
