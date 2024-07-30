import { execa } from "execa";
import chalk from "chalk";
import util from "util";
import { ConfigCliArgs } from "../../lib/@types/index.js";
import {
    getAppConfig,
    getGlobalConfig,
    getLocalConfig,
} from "../../lib/utils/index.js";

function getEditor() {
    return process.env.EDITOR || "nvim";
}

async function openInEditor({
    editor,
    filepath,
}: {
    editor: string;
    filepath: string;
}) {
    try {
        const childProcess = execa(editor, [filepath], {
            stdio: "inherit",
        });

        await childProcess;
    } catch (error: any) {
        console.error(error);
        process.exit(1);
    }
}

// async function set(key: string, value: string, cliOptions: CliOptions) {
//     try {
//         const { global, local } = cliOptions;

//         if (global) {
//             const globalConfig = getAppConfig();
//             globalConfig.set(key, value[0]);
//         }

//         if (local) {
//             const localConfig = await getLocalAppConfig();
//             localConfig.config.set(key, value[0]);
//         }
//     } catch (error) {
//         console.error(error);
//         process.exit(1);
//     }
// }

// async function get(key, cliOptions) {
//     try {
//         const { global, local } = cliOptions;

//         if (global) {
//             const globalConfig = getGlobalAppConfig();
//             console.log(globalConfig.get(key));
//         }

//         if (local) {
//             const localConfig = await getLocalAppConfig();
//             console.log(localConfig.config.get(key));
//         }
//     } catch (error) {
//         console.error(error);
//         process.exit(1);
//     }
// }

async function getConfig(args: ConfigCliArgs) {
    return args.global ? await getGlobalConfig() : await getLocalConfig({});
}

async function listConfig(args: ConfigCliArgs) {
    try {
        const { global, local } = args;

        if (!global && !local) {
            // Show the combined configuration
            const config = await getAppConfig({});

            console.log(
                util.inspect(config, {
                    depth: null,
                    colors: true,
                }),
            );

            return;
        }

        const result = await getConfig(args);

        if (!result) {
            console.log(
                `No ${global ? "global" : "local"} configuration found.`,
            );

            return;
        }

        console.log(
            util.inspect(result.config, {
                depth: null,
                colors: true,
            }),
        );
    } catch (error: any) {
        console.error(error);
        process.exit(1);
    }
}

async function editConfig(args: ConfigCliArgs) {
    try {
        const { global, local } = args;

        if (!global && !local) {
            console.log(chalk.red("Please specify --global or --local"));
            return;
        }

        const result = await getConfig(args);

        if (!result) {
            console.log(
                `No ${global ? "global" : "local"} configuration found.`,
            );

            return;
        }

        await openInEditor({
            editor: getEditor(),
            filepath: result.filepath,
        });
    } catch (error: any) {
        console.error(error);
        process.exit(1);
    }
}

export const config = {
    async process(args: ConfigCliArgs) {
        try {
            const { list, edit } = args;

            if (list) {
                listConfig(args);
                return;
            }

            if (edit) {
                editConfig(args);
                return;
            }
        } catch (error: any) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default config;
