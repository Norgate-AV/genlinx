import AdmZip from "adm-zip";
import {
    EnvItem,
    FileType,
    GeneralItem,
    ModuleItem,
    SourceItem,
    WorkspaceItem,
} from "./index.js";
import { ArchiveConfig, ArchiveItem, FileReference } from "./@types/index.js";

export class ArchiveItemFactory {
    public static create(
        builder: AdmZip,
        file: FileReference,
        options: ArchiveConfig,
    ): ArchiveItem {
        switch (file.type) {
            case FileType.apw.Workspace:
                return new WorkspaceItem(builder, file, options);
            case FileType.apw.Module:
                return new ModuleItem(builder, file, options);
            case FileType.apw.Source:
            case FileType.apw.MasterSrc:
                return new SourceItem(builder, file, options);
            case FileType.env:
                return new EnvItem(builder, file, options);
            default:
                return new GeneralItem(builder, file, options);
        }
    }
}

export default ArchiveItemFactory;
