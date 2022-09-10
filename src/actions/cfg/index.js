import fs from "fs-extra";
import { APW, CfgBuilder, Options } from "../../../lib";
import { getGlobalAppConfig, getLocalAppConfig } from "../../../lib/utils";

export const cfg = {
    async create(filePath, cliOptions) {
        try {
            const globalConfig = getGlobalAppConfig();
            const localConfig = getLocalAppConfig(filePath);

            const apw = new APW(filePath);

            const options = Options.getCfgOptions(
                apw,
                cliOptions,
                localConfig.cfg,
                globalConfig.cfg,
            );

            const cfgBuilder = new CfgBuilder(apw, options);
            const cfg = cfgBuilder.build();

            await fs.writeFile(options.outputFile, cfg);
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default cfg;
