import { CliBuildOptions } from "./index.js";

export type BuildConfig = {
    nlrc: {
        path: string;
        option: {
            cfg: "-CFG";
            includePath: "-I";
            modulePath: "-M";
            libraryPath: "-L";
        };
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

export type BuildOptions = BuildConfig & CliBuildOptions;
