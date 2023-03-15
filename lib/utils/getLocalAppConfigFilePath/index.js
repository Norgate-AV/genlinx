import findUp from "find-up";

export async function getLocalAppConfigFilePath(cwd) {
    const filePath = await findUp(".genlinxrc.json", {
        cwd,
    });

    return filePath;
}

export default getLocalAppConfigFilePath;
