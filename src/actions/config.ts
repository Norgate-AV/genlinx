import util from "node:util";
import { execa } from "execa";
import chalk from "chalk";
import { default as getConfigVal } from "lodash/get.js";
import { Config, ConfigCliArgs } from "../../lib/@types/index.js";
import {
    getAppConfig,
    getGlobalConfig,
    getLocalConfig,
} from "../../lib/utils/index.js";

function getEditor(config: Config) {
    return getConfigVal(config, "core.editor") || process.env.EDITOR || "code";
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

async function get(key: string, args: ConfigCliArgs) {
    try {
        const { global, local } = args;

        if (!key) {
            console.log("No key specified");
            return;
        }

        if (!global && !local) {
            // Get from combined config
            const config = await getAppConfig({});

            const value = getConfigVal(config, key);

            if (!value) {
                console.log(`No configuration found for key: ${key}`);
                return;
            }

            console.log(value);
            return;
        }

        const result = await getConfig(args);

        if (!result) {
            console.log(
                `No ${global ? "global" : "local"} configuration found.`,
            );

            return;
        }

        const config = result.config as Partial<Config>;

        const value = getConfigVal(config, key);

        if (!value) {
            console.log(`No configuration found for key: ${key}`);
            return;
        }

        console.log(value);
    } catch (error) {
        console.error(error);
        process.exit(1);
    }
}

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
            editor: getEditor(result.config),
            filepath: result.filepath,
        });
    } catch (error: any) {
        console.error(error);
        process.exit(1);
    }
}

export const config = {
    async process(key: string, args: ConfigCliArgs) {
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

            if (key) {
                get(key, args);
            }
        } catch (error: any) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default config;
