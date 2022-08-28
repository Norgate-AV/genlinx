import execa from "execa";
import { NLRC } from "../../../lib";
import { getAppConfig } from "../../../lib/utils";

export const build = {
    async build(filePath, options) {
        try {
            const config = getAppConfig();
            const command = NLRC.getCfgBuildCommand(filePath, options, config);

            const { executable } = config.build;
            await execa(
                executable.path,
                [executable.args, command.path, command.args],
                {
                    cwd: process.cwd(),
                },
            );
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default build;
