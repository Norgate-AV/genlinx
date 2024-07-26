import { FileReference } from "./index.js";

export type File = FileReference & {
    exists: boolean;
    isExtra: boolean;
    content?: string;
};
