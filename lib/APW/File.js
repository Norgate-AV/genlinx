import fs from "fs-extra";

export class File {
    constructor(type, id, path) {
        this.type = type;
        this.id = id;
        this.path = path;
        this.exists = fs.existsSync(this.path);
    }
}

export default File;
