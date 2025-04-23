export type BuildConfig = {
    nlrc: {
        path: string;
        includePath: Array<string>;
        modulePath: Array<string>;
        libraryPath: Array<string>;
    };
    shell: {
        path: string;
    };
    all: boolean;
    createCfg: boolean;
};
