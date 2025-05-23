import { AmxFileType as FileType } from "../../types/index.js";

export const AmxExtensions: Record<FileType, string> = {
    [FileType.Workspace]: ".apw",
    [FileType.Module]: ".axs",
    [FileType.MasterSrc]: ".axs",
    [FileType.Source]: ".axs",
    [FileType.Include]: ".axi",
    [FileType.IR]: ".irl",
    [FileType.TP4]: ".tp4",
    [FileType.TP5]: ".tp5",
    [FileType.TPD]: ".tpd",
    [FileType.Duet]: ".jar",
    [FileType.XDD]: ".xdd",
    [FileType.KPD]: ".kpd",
    [FileType.AXB]: ".axb",
    [FileType.TKO]: ".tko",
    [FileType.IRDB]: ".irdb",
    [FileType.IRNDB]: ".irndb",
    [FileType.TOK]: ".tok",
    [FileType.TKN]: ".tkn",
    [FileType.KPB]: ".kpb",
    [FileType.Other]: "",
} as const;
