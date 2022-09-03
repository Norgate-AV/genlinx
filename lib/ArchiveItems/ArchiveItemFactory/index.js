import { APW } from "../..";
import { GeneralItem } from "../GeneralItem";
import { ModuleItem } from "../ModuleItem";
import { SourceItem } from "../SourceItem";
import { WorkspaceItem } from "../WorkspaceItem";

export class ArchiveItemFactory {
    static create(builder, file, options) {
        switch (file.type) {
            case APW.fileType.workspace:
                return new WorkspaceItem(builder, file, options);
            case APW.fileType.module:
                return new ModuleItem(builder, file, options);
            case APW.fileType.source:
            case APW.fileType.masterSrc:
                return new SourceItem(builder, file, options);
            default:
                return new GeneralItem(builder, file, options);
        }
    }
}

export default ArchiveItemFactory;
