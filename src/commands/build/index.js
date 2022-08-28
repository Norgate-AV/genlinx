import { spawn } from "child_process";
import { NLRC } from "../../../lib";
import { getAppConfig } from "../../../lib/utils";

export const build = {
    build(filePath, options) {
        try {
            const config = getAppConfig();
            const command = NLRC.getCfgBuildCommand(filePath, options, config);

            console.log(command);
            console.log(spawn(command, { shell: true }));
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default build;
