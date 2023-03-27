export const build = {
    type: "object",
    properties: {
        nlrc: {
            type: "object",
            properties: {
                path: {
                    type: "string",
                },
                option: {
                    type: "object",
                    properties: {
                        cfg: {
                            type: "string",
                        },
                        includePath: {
                            type: "string",
                        },
                        modulePath: {
                            type: "string",
                        },
                        libraryPath: {
                            type: "string",
                        },
                    },
                },
                includePath: {
                    type: "array",
                    uniqueItems: true,
                },
                modulePath: {
                    type: "array",
                    uniqueItems: true,
                },
                libraryPath: {
                    type: "array",
                    uniqueItems: true,
                },
            },
            required: [
                "path",
                "option",
                "includePath",
                "modulePath",
                "libraryPath",
            ],
        },
        shell: {
            type: "object",
            properties: {
                path: {
                    type: "string",
                },
            },
            required: ["path"],
        },
        all: {
            type: "boolean",
        },
        createCfg: {
            type: "boolean",
        },
    },
    required: ["nlrc", "shell", "all", "createCfg"],
    additionalProperties: false,
};
