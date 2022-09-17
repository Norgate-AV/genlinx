export const defaultGlobalAppConfig = {
    cfg: {
        outputFileSuffix: "build.cfg",
        outputLogFileSuffix: "build.log",
        outputLogFileOption: "N",
        outputLogConsoleOption: true,
        buildWithDebugInformation: false,
        buildWithSource: false,
        includePath: [
            "C:\\Program Files (x86)\\Common Files\\AMXShare\\AXIs",
            "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18\\includes",
        ],
        modulePath: [
            "C:\\Program Files (x86)\\Common Files\\AMXShare\\Duet\\bundle",
            "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18",
            "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18\\monitors",
            "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18\\monitors-duet",
            "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18\\monitors-netlinx",
        ],
        libraryPath: ["C:\\Program Files (x86)\\Common Files\\AMXShare\\SYCs"],
        all: false,
    },
    archive: {
        outputFileSuffix: "archive.zip",
        includeCompiledSourceFiles: true,
        includeCompiledModuleFiles: true,
        includeFilesNotInWorkspace: true,
        extraFileSearchLocations: [
            "C:\\Program Files (x86)\\Common Files\\AMXShare",
            "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18",
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
            path: "C:\\Program Files (x86)\\Common Files\\AMXShare\\COM\\NLRC.exe",
            option: {
                cfg: "-CFG",
                includePath: "-I",
                modulePath: "-M",
                libraryPath: "-L",
            },
            includePath: [
                "C:\\Program Files (x86)\\Common Files\\AMXShare\\AXIs",
                "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18\\includes",
            ],
            modulePath: [
                "C:\\Program Files (x86)\\Common Files\\AMXShare\\Duet\\bundle",
                "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18",
                "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18\\monitors",
                "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18\\monitors-duet",
                "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18\\monitors-netlinx",
            ],
            libraryPath: [
                "C:\\Program Files (x86)\\Common Files\\AMXShare\\SYCs",
            ],
        },
        shell: {
            path: "C:\\Windows\\System32\\cmd.exe",
        },
        all: false,
        createCfg: true,
    },
};

export default defaultGlobalAppConfig;
