import fs from "fs-extra";
import path from "path";

export class APW {
    static fileType = {
        module: "Module",
        masterSrc: "MasterSrc",
        source: "Source",
        include: "Include",
        tp4: "TP4",
        tp5: "TP5",
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
        this.#fileReferences = this.#getFileReferences(data);
        this.#uniqueFileReferences = this.#getUniqueFileReferences();
        this.#uniqueFileReferences = this.#sortFileReferences(
            this.#uniqueFileReferences,
        );
        this.#setExistingFileReferences();
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
        const pattern = /^<Workspace.+\n.+>(?<id>.+)<.+>$/gm;
        const match = data.match(pattern);

        if (!match) {
            throw new Error("No ID found");
        }

        return match.groups.id;
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
            fileReferences.push({ type, id, path, exists: false });
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

    #setExistingFileReferences() {
        this.#uniqueFileReferences.map((file) => {
            if (fs.existsSync(file.path)) {
                return { ...file, exists: true };
            }

            return file;
        });
    }

    #sortFileReferences(files) {
        return files.sort((a, b) => a.path.localeCompare(b.path));
    }

    static find(dir) {
        const result = [];

        const files = fs.readdirSync(dir);

        for (const file of files) {
            const filePath = path.join(dir, file);
            const stats = fs.statSync(filePath);

            if (stats.isDirectory()) {
                const apw = APW.find(filePath);

                if (apw) {
                    result.push(apw);
                }
            } else if (file.endsWith(".apw")) {
                result.push(filePath);
            }
        }

        return result;
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

    get includeFileDirectories() {
        const includeFiles = this.#uniqueFileReferences.filter(
            (file) => file.type === APW.fileType.include,
        );

        const includeFileDirectories = includeFiles.map((file) =>
            path.resolve(this.#filePath, path.dirname(file.path)),
        );

        return [
            ...new Map(
                includeFileDirectories.map((directory) => [
                    directory.path,
                    directory,
                ]),
            ).values(),
        ];
    }

    get id() {
        return this.#id;
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
