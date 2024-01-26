import { FileId } from "./FileId";

export type FileReference = FileId & {
    type: string;
    path: string;
};
