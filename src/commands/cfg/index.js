import fs from "fs-extra";
import { APW, CfgBuilder } from "../../../lib";
import { getAppConfig } from "../../../lib/utils";

export const cfg = {
    async create(filePath, options) {
        try {
            const config = getAppConfig();
            const apw = new APW(filePath);
            const cfgBuilder = new CfgBuilder(apw, options, config);
            const cfg = cfgBuilder.build();

            await fs.writeFile(`${apw.id}.build.cfg`, cfg);
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default cfg;
