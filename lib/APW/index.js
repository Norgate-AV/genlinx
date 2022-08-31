import fs from "fs-extra";
import path from "path";
import { File } from "./File";

export class APW {
    static fileType = {
        workspace: "Workspace",
        module: "Module",
        masterSrc: "MasterSrc",
        source: "Source",
        include: "Include",
        ir: "IR",
        tp4: "TP4",
        tp5: "TP5",
        tpd: "TPD",
        kpd: "KPD",
        axb: "AXB",
        tko: "TKO",
        irDb: "IRDB",
        irnDb: "IRNDB",
        duet: "Duet",
        tok: "TOK",
        tkn: "TKN",
        kpb: "KPB",
        xdd: "XDD",
        other: "Other",
    };

    static fileExtensions = {
        [APW.fileType.workspace]: ".apw",
        [APW.fileType.module]: ".axs",
        [APW.fileType.masterSrc]: ".axs",
        [APW.fileType.source]: ".axs",
        [APW.fileType.include]: ".axi",
        [APW.fileType.ir]: ".irl",
        [APW.fileType.tp4]: ".tp4",
        [APW.fileType.tp5]: ".tp5",
    };

    static compiledFileExtensions = {
        [APW.fileType.module]: ".tko",
        [APW.fileType.masterSrc]: ".tkn",
        [APW.fileType.source]: ".tkn",
    };

    #filePath;

    #id;

    #fileReferences = [];

    #uniqueFileReferences = [];

    constructor(filePath) {
        this.#filePath = path.isAbsolute(filePath)
            ? filePath
            : path.resolve(process.cwd(), filePath);

        const data = this.#read(this.#filePath);

        this.#id = this.#getId(data);

        this.#fileReferences = this.#getFileReferences(data).sort(
            APW.#sortFileReferences,
        );

        this.#uniqueFileReferences = this.#getUniqueFileReferences();
    }

    #read(filePath) {
        if (!fs.existsSync(filePath)) {
            throw new Error("File does not exist");
        }

        const data = fs.readFileSync(filePath, {
            encoding: "utf8",
        });

        if (!/<!DOCTYPE Workspace \[/.test(data)) {
            throw new Error("Not a Netlinx Workspace file");
        }

        return data;
    }

    #getId(data) {
        const pattern = /<Workspace.+\n.+>(?<id>.+)<.+>/m;
        const match = data.match(pattern).groups;

        if (!match) {
            throw new Error("No Workspace ID found");
        }

        return match.id;
    }

    #getFileReferences(data) {
        const pattern =
            /<File.+Type="(?<type>.+)".+\n.+>(?<id>.+)<.+\n.+>(?<path>.+)<.+\n.+\n.+(?:<DeviceMap.+\n.+\n.+\n)?.+File>/gm;
        const matches = data.matchAll(pattern);

        if (!matches) {
            throw new Error("No file references found");
        }

        const fileReferences = [];
        for (const match of matches) {
            const { type, id, path } = match.groups;
            fileReferences.push(new File(type, id, path));
        }

        return fileReferences;
    }

    #getUniqueFileReferences() {
        return [
            ...new Map(
                this.#fileReferences.map((file) => [file.path, file]),
            ).values(),
        ];
    }

    static #sortFileReferences(file, nextFile) {
        return file.path.localeCompare(nextFile.path);
    }

    get moduleFiles() {
        return this.#uniqueFileReferences.filter(
            (file) => file.type === APW.fileType.module,
        );
    }

    get masterSrcFiles() {
        return this.#uniqueFileReferences.filter(
            (file) => file.type === APW.fileType.masterSrc,
        );
    }

    get allFiles() {
        return [
            ...this.#uniqueFileReferences,
            new File(APW.fileType.workspace, this.#id, this.#filePath),
        ];
    }

    get includePath() {
        const includeFiles = this.#uniqueFileReferences.filter(
            (file) => file.type === APW.fileType.include,
        );

        const includeFileDirectories = includeFiles.map((file) => {
            const rootDirectory = path.dirname(this.#filePath);
            const absoluteFilePath = path.join(rootDirectory, file.path);
            return path.dirname(absoluteFilePath);
        });

        return [...new Set(includeFileDirectories)];
    }

    get id() {
        return this.#id.replace(" ", "-");
    }

    get totalFileCount() {
        return this.#fileReferences.length;
    }

    get uniqueFileCount() {
        return this.#uniqueFileReferences.length;
    }

    get missingFiles() {
        return this.#uniqueFileReferences.filter((file) => !file.exists);
    }
}

export default APW;
