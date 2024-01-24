import fs from "fs-extra";
import path from "path";
import { FileType } from "./FileType";

export type FileReference = {
    type: string;
    id: string;
    path: string;
    exists?: boolean;
    extra?: boolean;
};

const extensions: Record<string, string> = {
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
};

const compiledExtensions: Record<string, string> = {
    [FileType.Module]: ".tko",
    [FileType.MasterSrc]: ".tkn",
    [FileType.Source]: ".tkn",
};

const fileTypeFromExtension: Record<string, string> = {
    [extensions[FileType.Module]]: FileType.Module,
    [extensions[FileType.Include]]: FileType.Include,
    [extensions[FileType.Duet]]: FileType.Duet,
    [extensions[FileType.XDD]]: FileType.XDD,
    [extensions[FileType.Workspace]]: FileType.Workspace,
    [extensions[FileType.IR]]: FileType.IR,
    [extensions[FileType.TP4]]: FileType.TP4,
    [extensions[FileType.TP5]]: FileType.TP5,
    [extensions[FileType.TPD]]: FileType.TPD,
    [compiledExtensions[FileType.Module]]: FileType.Module,
    [compiledExtensions[FileType.Source]]: FileType.Source,
};

export class APW {
    private _filePath: string;

    private _id: string;

    private fileReferences: Array<FileReference> = [];

    private uniqueFileReferences: Array<FileReference> = [];

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
        } catch (error) {
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
        const match = data.match(pattern)?.groups;

        if (!match) {
            throw new Error("No Workspace ID found");
        }

        return match.id;
    }

    private async getFileReferences(
        data: string,
    ): Promise<Array<FileReference>> {
        const pattern =
            /<File.+Type="(?<type>.+)".+\r?\n?.*?<Identifier>(?<id>.+)<\/Identifier>\r?\n?.*?>(?<path>.+)<.+\r?\n?.*?\r?\n?.*?(?:<DeviceMap.+\r?\n?.*?\r?\n?.*?\r?\n?)?.*?<\/File>/gm;
        const matches = data.matchAll(pattern);

        const fileReferences: Array<FileReference> = [];
        for (const match of matches) {
            if (!match.groups) {
                continue;
            }

            const { type, id, path } = match.groups;

            const exists = await fs.pathExists(path);

            fileReferences.push({
                type,
                id,
                path,
                exists,
                extra: false,
            });
        }

        return fileReferences.sort(APW.sortFileReferences);
    }

    private getUniqueFileReferences(): Array<FileReference> {
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

            const { id } = match.groups;

            if (this.isInWorkspace(id)) {
                continue;
            }

            fileReferences.push(id);
        }

        return [...new Set(fileReferences)];
    }

    public static getFileType(file: string): string {
        return fileTypeFromExtension[path.extname(file)];
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

    get allFiles(): Array<FileReference> {
        return [
            ...this.uniqueFileReferences,
            {
                type: FileType.Workspace,
                id: this.id,
                path: this.filePath,
                exists: true,
                extra: false,
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
