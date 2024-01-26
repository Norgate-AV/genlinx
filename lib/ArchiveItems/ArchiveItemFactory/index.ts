import { FileType } from "../FileType";
import { EnvItem } from "../EnvItem";
import { GeneralItem } from "../GeneralItem";
import { ModuleItem } from "../ModuleItem";
import { SourceItem } from "../SourceItem";
import { WorkspaceItem } from "../WorkspaceItem";
import { ArchiveItem } from "../../@types/ArchiveItem";
import AdmZip from "adm-zip";
import { FileReference } from "../../@types/FileReference";
import { ArchiveConfig } from "../../@types/ArchiveConfig";

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
