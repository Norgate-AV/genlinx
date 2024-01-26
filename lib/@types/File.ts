import { FileReference } from "./FileReference";

export type File = FileReference & {
    exists: boolean;
    isExtra: boolean;
};
