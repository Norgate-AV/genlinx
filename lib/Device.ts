import { Date } from "./@types/index.js";

export class Device {
    public ip: string;
    public system: number;
    public date: Date;
    public time: string;
    public day: string;
    public mac: string;
    public hostname: string;
    public id: string;

    constructor() {
        this.ip = "";
        this.system = 0;
        this.date = { numeric: "", text: "" };
        this.time = "";
        this.day = "";
        this.mac = "";
        this.hostname = "";
        this.id = "";
    }
}
