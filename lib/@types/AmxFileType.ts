import { ObjectTypes } from "./index.js";

export const AmxFileType = {
    Workspace: "Workspace",
    Module: "Module",
    MasterSrc: "MasterSrc",
    Source: "Source",
    Include: "Include",
    IR: "IR",
    TP4: "TP4",
    TP5: "TP5",
    TPD: "TPD",
    KPD: "KPD",
    AXB: "AXB",
    TKO: "TKO",
    IRDB: "IRDB",
    IRNDB: "IRNDB",
    Duet: "Duet",
    TOK: "TOK",
    TKN: "TKN",
    KPB: "KPB",
    XDD: "XDD",
    Other: "Other",
} as const;

export type AmxFileType = ObjectTypes<typeof AmxFileType>;
