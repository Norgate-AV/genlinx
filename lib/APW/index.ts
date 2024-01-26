import fs from "fs-extra";
import path from "path";
import { FileType } from "../@types/FileType";
import { FileReference } from "../@types/FileReference";
import { File } from "../@types/File";
import { FileId } from "../@types/FileId";

export const extensions: Record<FileType, string> = {
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
};

export const compiledExtensions: Record<FileType, string> = {
    [FileType.Workspace]: ".apw",
    [FileType.Module]: ".tko",
    [FileType.MasterSrc]: ".tkn",
    [FileType.Source]: ".tkn",
    [FileType.Include]: ".tkn",
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
};

export class APW {
    private _filePath: string | null = null;

    private _id: string | null = null;

    private fileReferences: Array<File> = [];

    private uniqueFileReferences: Array<File> = [];

    public constructor(filePath: string) {
        this.filePath = path.isAbsolute(filePath)
            ? filePath
            : path.resolve(process.cwd(), filePath);
    }

    public async load(): Promise<void> {
        try {
            const data = await this.read();

            this.id = APW.getId(data);
            this.fileReferences = await this.getFileReferences(data);
            this.uniqueFileReferences = this.getUniqueFileReferences();
        } catch (error: any) {
            throw new Error(`Failed to load APW file: ${error.message}`);
        }
    }

    private async read(): Promise<string> {
        const filePath = this.filePath;

        const exists = await fs.pathExists(filePath);

        if (!exists) {
            throw new Error("File does not exist");
        }

        const data = await fs.readFile(filePath, {
            encoding: "utf8",
        });

        if (!/<!DOCTYPE Workspace \[/.test(data)) {
            throw new Error("Not a Netlinx Workspace file");
        }

        return data;
    }

    private static getId(data: string): string {
        const pattern = /<Workspace.+\r?\n?.*?<Identifier>(?<id>.+)<.+>/m;
        const match = data.match(pattern)?.groups as FileId;

        if (!match) {
            throw new Error("No Workspace ID found");
        }

        return match.id;
    }

    private async getFileReferences(data: string): Promise<Array<File>> {
        const pattern =
            /<File.+Type="(?<type>.+)".+\r?\n?.*?<Identifier>(?<id>.+)<\/Identifier>\r?\n?.*?>(?<path>.+)<.+\r?\n?.*?\r?\n?.*?(?:<DeviceMap.+\r?\n?.*?\r?\n?.*?\r?\n?)?.*?<\/File>/gm;
        const matches = data.matchAll(pattern);

        const fileReferences: Array<File> = [];
        for (const match of matches) {
            if (!match.groups) {
                continue;
            }

            const { type, id, path } = match.groups as FileReference;

            const exists = await fs.pathExists(path);

            fileReferences.push({
                type,
                id,
                path,
                exists,
                isExtra: false,
            });
        }

        return fileReferences.sort(APW.sortFileReferences);
    }

    private getUniqueFileReferences(): Array<File> {
        return [
            ...new Map(
                this.fileReferences.map((file) => [file.path, file]),
            ).values(),
        ];
    }

    private static sortFileReferences(
        file: FileReference,
        nextFile: FileReference,
    ): number {
        return file.path.localeCompare(nextFile.path);
    }

    private isInWorkspace(id: string) {
        return this.uniqueFileReferences.find(
            (file) => file.id === id || file.path.includes(id),
        );
    }

    private async searchForExtraFileReferences(
        file: string,
        pattern: RegExp,
    ): Promise<Array<string>> {
        const data = await fs.readFile(file, { encoding: "utf-8" });

        const matches = data.matchAll(pattern);

        const fileReferences: Array<string> = [];
        for (const match of matches) {
            if (!match.groups) {
                continue;
            }

            const { id } = match.groups as FileId;

            if (this.isInWorkspace(id)) {
                continue;
            }

            fileReferences.push(id);
        }

        return [...new Set(fileReferences)];
    }

    public static getFileType(file: string): FileType | null {
        if (!(path.extname(file) in extensions)) {
            return null;
        }

        const fileType = Object.keys(extensions).find(
            (key) => extensions[key as FileType] === path.extname(file),
        );

        return fileType as FileType;
    }

    public static fileIsReadable(file: string): boolean {
        return path.extname(file) === ".axs" || path.extname(file) === ".axi";
    }

    public async getExtraFileReferencesFromFile(
        file: string,
    ): Promise<Array<string>> {
        const extraFiles: Array<string> = [];

        if (!APW.fileIsReadable(file)) {
            return extraFiles;
        }

        const extraIncludeFiles = await this.searchForExtraFileReferences(
            file,
            /#(?:include).+'(?<id>.+)'/gim,
        );

        const extraModuleFiles = await this.searchForExtraFileReferences(
            file,
            /^(?:define_module).+'(?<id>.+)'/gim,
        );

        extraFiles.push(...extraIncludeFiles, ...extraModuleFiles);

        return [...new Set(extraFiles)];
    }

    public async getExtraFileReferences(): Promise<Array<String>> {
        const extraFiles: Array<string> = [];

        for (const file of this.uniqueFileReferences) {
            if (!file.exists) {
                continue;
            }

            extraFiles.push(
                ...(await this.getExtraFileReferencesFromFile(file.path)),
            );
        }

        return [...new Set(extraFiles)];
    }

    get moduleFiles(): Array<FileReference> {
        return this.uniqueFileReferences.filter(
            (file) => file.type === FileType.Module,
        );
    }

    get masterSrcFiles(): Array<FileReference> {
        return this.uniqueFileReferences.filter(
            (file) => file.type === FileType.MasterSrc,
        );
    }

    get allFiles(): Array<File> {
        return [
            ...this.uniqueFileReferences,
            {
                type: FileType.Workspace,
                id: this.id,
                path: this.filePath,
                exists: true,
                isExtra: false,
            },
        ];
    }

    private getFileDirectories(files: Array<FileReference>): Array<string> {
        const directories = files.map((file) => {
            const rootDirectory = path.dirname(this.filePath);
            const absoluteFilePath = path.join(rootDirectory, file.path);
            return path.dirname(absoluteFilePath);
        });

        return [...new Set(directories)];
    }

    get includePath(): Array<string> {
        const includeFiles = this.uniqueFileReferences.filter(
            (file) => file.type === FileType.Include,
        );

        return this.getFileDirectories(includeFiles);
    }

    get masterSrcPath(): Array<string> {
        const masterSrcFiles = this.uniqueFileReferences.filter(
            (file) => file.type === FileType.MasterSrc,
        );

        return this.getFileDirectories(masterSrcFiles);
    }

    get modulePath(): Array<string> {
        const moduleFiles = this.uniqueFileReferences.filter(
            (file) =>
                file.type === FileType.Module ||
                file.type === FileType.Duet ||
                file.type === FileType.XDD,
        );

        return this.getFileDirectories(moduleFiles);
    }

    get id(): string {
        return this._id.replace(" ", "-");
    }

    set id(id: string) {
        this._id = id;
    }

    get totalFileCount(): number {
        return this.fileReferences.length;
    }

    get uniqueFileCount(): number {
        return this.uniqueFileReferences.length;
    }

    get missingFiles(): Array<FileReference> {
        return this.uniqueFileReferences.filter((file) => !file.exists);
    }

    get filePath(): string {
        return this._filePath;
    }

    set filePath(filePath: string) {
        this._filePath = filePath;
    }
}

export default APW;
