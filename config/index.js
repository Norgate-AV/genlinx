export const defaultAppConfig = {
    cfg: {
        outputFile: "build.cfg",
        outputLogFile: "build.log",
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
    },
    archive: {
        outputFile: "archive.zip",
        includeCompiledSourceFiles: true,
        includeCompiledModuleFiles: true,
        includeFilesNotInWorkspace: true,
        extraFileLocations: [
            "C:\\Program Files (x86)\\Common Files\\AMXShare",
            "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18",
        ],
    },
    build: {
        nlrc: {
            path: "C:\\Program Files (x86)\\Common Files\\AMXShare\\COM\\NLRC.exe",
            option: {
                cfg: "-CFG",
            },
        },
        shell: {
            path: "C:\\Windows\\System32\\cmd.exe",
        },
    },
};

export default defaultAppConfig;
