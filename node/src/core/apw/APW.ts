import fs from "node:fs/promises";
import path from "node:path";
import {
    File,
    FileId,
    FileReference,
    AmxFileType as FileType,
} from "../../types/index.js";
import { AmxExtensions } from "./index.js";
import { pathExists } from "../../utils/index.js";

export class APW {
    private _filePath: string | null = null;

    private _id: string | null = null;

    private readonly fileReferences: Array<File> = [];

    private readonly uniqueFileReferences: Array<File> = [];

    public constructor(filePath: string) {
        this.filePath = path.isAbsolute(filePath)
            ? filePath
            : path.resolve(process.cwd(), filePath);
    }

    public async load(): Promise<void> {
        try {
            const data = await this.read();

            this.id = APW.getId(data);
            this.fileReferences.push(...(await this.getFileReferences(data)));
            this.uniqueFileReferences.push(...this.getUniqueFileReferences());
        } catch (error: any) {
            throw new Error(`Failed to load APW file: ${error.message}`);
        }
    }

    private async read(): Promise<string> {
        const { filePath } = this;

        const exists = await pathExists(filePath);

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

            const exists = await pathExists(path);

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

        return [...new Set<string>(fileReferences)];
    }

    public static getFileType(file: string): FileType {
        const extension = path.extname(file);

        if (!(extension in Object.values<string>(AmxExtensions))) {
            return FileType.Other;
        }

        const fileType = Object.keys(AmxExtensions).find(
            (key) => AmxExtensions[key as FileType] === extension,
        );

        return (fileType as FileType) || FileType.Other;
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

        return [...new Set<string>(extraFiles)];
    }

    public async getExtraFileReferences(): Promise<Array<string>> {
        const extraFiles: Array<string> = [];

        for (const file of this.uniqueFileReferences) {
            if (!file.exists) {
                continue;
            }

            extraFiles.push(
                ...(await this.getExtraFileReferencesFromFile(file.path)),
            );
        }

        return [...new Set<string>(extraFiles)];
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

    private getFileDirectories(files: Array<File>): Array<string> {
        const directories = files.map((file) => {
            const rootDirectory = path.dirname(this.filePath);
            const absoluteFilePath = path.join(rootDirectory, file.path);
            return path.dirname(absoluteFilePath);
        });

        return [...new Set<string>(directories)];
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
        return this._id ? this._id.replace(" ", "-") : "";
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
        return this._filePath || "";
    }

    set filePath(filePath: string) {
        this._filePath = filePath;
    }
}

export default APW;
