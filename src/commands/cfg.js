import fs from "fs-extra";
import { APW, CfgBuilder } from "../../lib";
import { getAppConfig } from "../../lib/utils";

export default {
    create(filePath, options) {
        try {
            const config = getAppConfig();
            const apw = new APW(filePath);
            const cfgBuilder = new CfgBuilder(apw, options, config);
            const cfg = cfgBuilder.build();

            fs.writeFileSync(options.output, cfg);
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
};
