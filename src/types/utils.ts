export type ObjectTypes<T> = T[keyof T];

export type ShellCommand = {
    path: string;
    args: Array<string>;
};

export interface Date {
    numeric: string;
    text: string;
}
