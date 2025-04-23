export type BuildLogCollection = {
    error: Array<string>;
    warning: Array<string>;
};

export type BuildLogCollectionKey = keyof BuildLogCollection;
