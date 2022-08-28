export const defaultAppConfig = {
    cfg: {
        outputLogFile: "build.log",
        outputLogFileOption: "N",
        outputLogConsoleOption: true,
        buildWithDebug: false,
        buildWithSource: false,
        includePath: [
            "C:\\Program Files (x86)\\Common Files\\AMXShare\\AXIs\\",
            "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18\\includes\\",
        ],
        modulePath: [
            "C:\\Program Files (x86)\\Common Files\\AMXShare\\Duet\\bundle\\",
            "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18\\",
            "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18\\monitors\\",
            "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18\\monitors-duet\\",
            "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18\\monitors-netlinx\\",
        ],
        libraryPath: [
            "C:\\Program Files (x86)\\Common Files\\AMXShare\\SYCs\\",
        ],
    },
    archive: {},
    build: {
        nlrc: {
            path: "C:\\Program Files (x86)\\Common Files\\AMXShare\\COM\\NLRC.exe",
        },
    },
};

export default defaultAppConfig;
