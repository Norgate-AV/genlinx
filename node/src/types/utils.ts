export type ObjectTypes<T> = T[keyof T];

export type ShellCommand = {
    path: string;
    args: Array<string>;
};
