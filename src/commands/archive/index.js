import { APW, ArchiveBuilder, Options } from "../../../lib";
import { getAppConfig } from "../../../lib/utils";

export const archive = {
    async create(filePath, cliOptions) {
        try {
            const config = getAppConfig();

            const apw = new APW(filePath);

            const options = Options.getArchiveOptions(
                apw,
                cliOptions,
                config.archive,
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
