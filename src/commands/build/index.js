import { execa } from "execa";
import { NLRC } from "../../../lib";
import { getAppConfig } from "../../../lib/utils";

export const build = {
    async build(filePath, options) {
        try {
            const config = getAppConfig();
            const command = NLRC.getCfgBuildCommand(filePath, options, config);

            const { executable } = config.build;
            const result = await execa(executable.path, [
                "/c",
                command.split(" "),
            ]);

            if (result.failed) {
                throw new Error(result.stderr);
            }
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default build;
