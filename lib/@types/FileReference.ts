import { FileId } from "./index.js";

export type FileReference = FileId & {
    type: string;
    path: string;
};
