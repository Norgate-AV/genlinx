import { Config } from "../lib/@types/index.js";

export const config: Config = {
    cfg: {
        outputFileSuffix: "build.cfg",
        outputLogFileSuffix: "build.log",
        outputLogFileOption: "N",
        outputLogConsoleOption: true,
        buildWithDebugInformation: false,
        buildWithSource: false,
        includePath: ["C:/Program Files (x86)/Common Files/AMXShare/AXIs"],
        modulePath: [
            "C:/Program Files (x86)/Common Files/AMXShare/Duet/bundle",
            "C:/Program Files (x86)/Common Files/AMXShare/Duet/lib",
            "C:/Program Files (x86)/Common Files/AMXShare/Duet/module",
        ],
        libraryPath: ["C:/Program Files (x86)/Common Files/AMXShare/SYCs"],
        all: false,
    },
    archive: {
        outputFileSuffix: "archive.zip",
        includeCompiledSourceFiles: true,
        includeCompiledModuleFiles: true,
        includeFilesNotInWorkspace: true,
        extraFileSearchLocations: [
            "C:/Program Files (x86)/Common Files/AMXShare",
        ],
        extraFileArchiveLocation: ".genlinx",
        all: false,
        ignoredFiles: [
            "G4API.axi",
            "NetLinx.axi",
            "SNAPI.axi",
            "UnicodeLib.axi",
            "componentssdk.jar",
            "componentssdkrt.jar",
            "DeviceDriverEngine.jar",
            "devicesdkrt.jar",
            "jregex1.2_01-bundle.jar",
            "js-14-bundle.jar",
            "json-bundle.jar",
            "picocontainer-1.3-bundle.jar",
            "snapirouter.jar",
            "snapirouter2.jar",
        ],
    },
    build: {
        nlrc: {
            path: "C:/Program Files (x86)/Common Files/AMXShare/COM/NLRC.exe",
            option: {
                cfg: "-CFG",
                includePath: "-I",
                modulePath: "-M",
                libraryPath: "-L",
            },
            includePath: ["C:/Program Files (x86)/Common Files/AMXShare/AXIs"],
            modulePath: [
                "C:/Program Files (x86)/Common Files/AMXShare/Duet/bundle",
                "C:/Program Files (x86)/Common Files/AMXShare/Duet/lib",
                "C:/Program Files (x86)/Common Files/AMXShare/Duet/module",
            ],
            libraryPath: ["C:/Program Files (x86)/Common Files/AMXShare/SYCs"],
        },
        shell: {
            path: "C:/Windows/System32/cmd.exe",
        },
        all: false,
        createCfg: true,
    },
};

export default config;
