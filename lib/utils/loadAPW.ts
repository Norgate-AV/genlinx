import chalk from "chalk";
import { APW } from "../index.js";

export async function loadAPW(filePath: string): Promise<APW> {
    const apw = new APW(filePath);

    await apw.load();

    if (apw.totalFileCount === 0) {
        console.log(
            chalk.red("No files found in workspace. There is nothing to do."),
        );

        process.exit(1);
    }

    return apw;
}
export default loadAPW;
