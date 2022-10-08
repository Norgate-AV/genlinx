import execa from "execa";
import { getGlobalAppConfig, getLocalAppConfig } from "../../../lib/utils";

async function openInEditor(config) {
    const editor = config.get("core.editor") || process.env.EDITOR || "code";

    const childProcess = execa(editor, [config.path], {
        stdio: "inherit",
    });

    await childProcess;
}

async function set(key, value, cliOptions) {
    try {
        const { global, local } = cliOptions;

        if (global) {
            const globalConfig = getGlobalAppConfig();
            globalConfig.set(key, value);
        }

        if (local) {
            const localConfig = await getLocalAppConfig();
            localConfig.set(key, value);
        }
    } catch (error) {
        console.error(error);
        process.exit(1);
    }
}

async function get(key, cliOptions) {
    try {
        const { global, local } = cliOptions;

        if (global) {
            const globalConfig = getGlobalAppConfig();
            console.log(globalConfig.get(key));
        }

        if (local) {
            const localConfig = await getLocalAppConfig();
            console.log(localConfig.get(key));
        }
    } catch (error) {
        console.error(error);
        process.exit(1);
    }
}

async function listConfig(cliOptions) {
    try {
        const { local } = cliOptions;

        if (local) {
            const localConfig = await getLocalAppConfig();
            console.log(JSON.stringify(localConfig.store, null, 4));
            return;
        }

        const globalConfig = getGlobalAppConfig();
        console.log(JSON.stringify(globalConfig.store, null, 4));
    } catch (error) {
        console.error(error);
        process.exit(1);
    }
}

async function editConfig(cliOptions) {
    try {
        const { local } = cliOptions;

        if (local) {
            const localConfig = await getLocalAppConfig();
            await openInEditor(localConfig);
            return;
        }

        const globalConfig = getGlobalAppConfig();
        await openInEditor(globalConfig);
    } catch (error) {
        console.error(error);
        process.exit(1);
    }
}

export const config = {
    async process(key, value, cliOptions) {
        try {
            const { global, local, list, edit } = cliOptions;

            if (global) {
                const globalConfig = getGlobalAppConfig();
                globalConfig.set(key, value);
            }

            if (local) {
                const localConfig = await getLocalAppConfig();
                localConfig.set(key, value);
            }

            if (!key && !value) {
                if (list) {
                    listConfig(cliOptions);
                }

                if (edit) {
                    await editConfig(cliOptions);
                }

                return;
            }

            if (!value) {
                get(key, cliOptions);
            }

            await set(key, value, cliOptions);
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default config;
