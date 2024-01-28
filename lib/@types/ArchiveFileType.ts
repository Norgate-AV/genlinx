import { AmxFileType, ObjectTypes } from "./index.js";

export const ArchiveFileType = {
    APW: AmxFileType,
    Script: "Script",
    Env: "Env",
} as const;

export type ArchiveFileType = ObjectTypes<typeof ArchiveFileType>;
