import { FileType } from "../FileType";
import { GeneralItem } from "../GeneralItem";
import { ModuleItem } from "../ModuleItem";
import { SourceItem } from "../SourceItem";
import { WorkspaceItem } from "../WorkspaceItem";

export class ArchiveItemFactory {
    static create(builder, file, options) {
        switch (file.type) {
            case FileType.apw.workspace:
                return new WorkspaceItem(builder, file, options);
            case FileType.apw.module:
                return new ModuleItem(builder, file, options);
            case FileType.apw.source:
            case FileType.apw.masterSrc:
                return new SourceItem(builder, file, options);
            default:
                return new GeneralItem(builder, file, options);
        }
    }
}

export default ArchiveItemFactory;
