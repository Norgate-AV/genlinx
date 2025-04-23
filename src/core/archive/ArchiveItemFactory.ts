import AdmZip from "adm-zip";
import {
    EnvItem,
    GeneralItem,
    ModuleItem,
    SourceItem,
    WorkspaceItem,
} from "./index.js";
import {
    ArchiveFileType as FileType,
    ArchiveItem,
    File,
    ArchiveOptions,
} from "../../@types/index.js";

export class ArchiveItemFactory {
    public static create(
        builder: AdmZip,
        file: File,
        options: ArchiveOptions,
    ): ArchiveItem {
        switch (file.type) {
            case FileType.APW.Workspace:
                return new WorkspaceItem(builder, file, options);
            case FileType.APW.Module:
                return new ModuleItem(builder, file, options);
            case FileType.APW.Source:
            case FileType.APW.MasterSrc:
                return new SourceItem(builder, file, options);
            case FileType.Env:
                return new EnvItem(builder, file, options);
            default:
                return new GeneralItem(builder, file, options);
        }
    }
}

export default ArchiveItemFactory;
