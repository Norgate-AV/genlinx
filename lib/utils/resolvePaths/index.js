import { resolve, isAbsolute } from "path";

export function resolvePaths(cwd, paths) {
    if (!paths) {
        return paths;
    }

    return paths.map((path) => {
        if (isAbsolute(path)) {
            return path;
        }

        return resolve(cwd, path);
    });
}

export default resolvePaths;
