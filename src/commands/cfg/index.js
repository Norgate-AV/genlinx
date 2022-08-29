import fs from "fs-extra";
import { APW, CfgBuilder, Options } from "../../../lib";
import { getAppConfig } from "../../../lib/utils";

export const cfg = {
    async create(filePath, cliOptions) {
        try {
            const config = getAppConfig();

            const apw = new APW(filePath);

            const options = Options.getCfgOptions(cliOptions, config.cfg);

            const cfgBuilder = new CfgBuilder(apw, options);
            const cfg = cfgBuilder.build();

            await fs.writeFile(`${apw.id}.build.cfg`, cfg);
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default cfg;
