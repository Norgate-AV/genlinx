// import { execa } from "execa";
import util from "util";
import { ConfigCliArgs } from "../../lib/@types/index.js";
import { getAppConfig } from "../../lib/utils/index.js";

// async function openInEditor(config: Config) {
//     const editor = config.get("core.editor") || process.env.EDITOR || "code";

//     const childProcess = execa(editor, [config.path], {
//         stdio: "inherit",
//     });

//     await childProcess;
// }

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

async function listConfig(args: ConfigCliArgs) {
    try {
        const config = await getAppConfig({});

        console.log(
            util.inspect(config, {
                depth: null,
                colors: true,
            }),
        );
        // const { local } = cliOptions;
        // if (local) {
        //     const localConfig = await getLocalAppConfig();
        //     console.log(
        //         util.inspect(localConfig.config.store, {
        //             depth: null,
        //             colors: true,
        //         }),
        //     );
        //     return;
        // }
        // const globalConfig = getGlobalAppConfig();
        // console.log(
        //     util.inspect(globalConfig.store, { depth: null, colors: true }),
        // );
    } catch (error) {
        console.error(error);
        process.exit(1);
    }
}

// async function editConfig(args: ConfigCliArgs) {
//     try {
//         // const { local } = cliOptions;
//         // if (local) {
//         //     const localConfig = await getLocalAppConfig();
//         //     await openInEditor(localConfig);
//         //     return;
//         // }
//         // const globalConfig = getGlobalAppConfig();
//         // await openInEditor(globalConfig);
//     } catch (error) {
//         console.error(error);
//         process.exit(1);
//     }
// }

export const config = {
    async process(args: ConfigCliArgs) {
        try {
            const { list } = args;

            if (list) {
                listConfig(args);
                return;
            }

            // console.log(key, value, args);

            // if (!key && value.length === 0) {
            //     // console.log("No key or value provided");

            //     // if (list) {
            //     //     listConfig(args);
            //     // }

            //     // if (edit) {
            //     //     await editConfig(cliOptions);
            //     // }

            //     return;
            // }

            // if (value.length === 0) {
            //     get(key, cliOptions);
            //     return;
            // }

            // await set(key, value, cliOptions);
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default config;
