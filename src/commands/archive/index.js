import { APW, ArchiveBuilder, Options } from "../../../lib";
import { getGlobalAppConfig, getLocalAppConfig } from "../../../lib/utils";

export const archive = {
    async create(filePath, cliOptions) {
        try {
            const globalConfig = getGlobalAppConfig();
            const localConfig = getLocalAppConfig(filePath);

            const apw = new APW(filePath);

            const options = Options.getArchiveOptions(
                apw,
                cliOptions,
                localConfig.archive ? localConfig.archive : {},
                globalConfig.archive,
            );

            const builder = new ArchiveBuilder(apw, options);
            builder.build();
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default archive;
