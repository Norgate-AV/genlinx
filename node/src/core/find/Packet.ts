import { Device } from "./index.js";

export class Packet {
    private readonly DATA_LENGTH_BYTE = 2;
    private readonly payload: Buffer;
    private data: Buffer;
    private readonly type: string;
    private device = new Device();
    private readonly days = {
        "1": "Monday",
        "2": "Tuesday",
        "3": "Wednesday",
        "4": "Thursday",
        "5": "Friday",
        "6": "Saturday",
        "7": "Sunday",
    };

    constructor(payload: Buffer) {
        this.payload = payload;
        this.data = this.getData();
        this.type = this.data.toString("utf8", 1, 2);
    }

    private getData(): Buffer {
        const { payload } = this;

        const length = payload[this.DATA_LENGTH_BYTE]!;
        const start = this.DATA_LENGTH_BYTE + 1;
        const end = start + length;

        const data = payload.subarray(start, end);

        if (data.length !== length) {
            throw new Error("Invalid data length");
        }

        return data;
    }

    public getType(): string {
        return this.type;
    }

    public parse(): Packet {
        this.device.system = this.getSystem();
        this.device.date.numeric = this.getDateNumeric();
        this.device.time = this.getTime();
        this.device.day = this.days[this.getDay() as keyof typeof this.days];
        this.device.date.text = this.getDateText();
        this.device.mac = this.getMac();
        this.device.hostname = this.getHostname();
        this.device.id = this.getId();

        return this;
    }

    public getDevice(): Device {
        return this.device;
    }

    private getDay(): string {
        const day = this.data[0]!.toString();

        this.data = this.data.subarray(3);

        return day;
    }

    private getDateNumeric(): string {
        const month = this.data[0]!.toString().padStart(2, "0");
        const day = this.data[1]!.toString().padStart(2, "0");
        const year = this.data.subarray(2, 4).readInt16BE();

        this.data = this.data.subarray(4);

        return `${month}/${day}/${year}`;
    }

    private getEnd(start: number = 0, end: string = "\0"): number {
        const { data } = this;

        const length = data.indexOf(end, start) - start;

        return start + length;
    }

    private getDateText(): string {
        const end = this.getEnd(0, "\0\0\0");
        const date = this.data.subarray(0, end).toString();

        this.data = this.data.subarray(end + 3);

        return date;
    }

    private getSystem(): number {
        const system = this.data.subarray(2, 4).readInt16BE();
        this.data = this.data.subarray(21);

        return system;
    }

    private getHostname(): string {
        const end = this.getEnd(0);
        const hostname = this.data.subarray(0, end).toString();

        this.data = this.data.subarray(end + 1);

        return hostname;
    }

    private getId(): string {
        const end = this.getEnd(0);

        return this.data.subarray(0, end).toString();
    }

    private getTime(): string {
        const hours = this.data[0]!.toString().padStart(2, "0");
        const minutes = this.data[1]!.toString().padStart(2, "0");
        const seconds = this.data[3]!.toString().padStart(2, "0");

        this.data = this.data.subarray(3);

        return `${hours}:${minutes}:${seconds}`;
    }

    private getMac(): string {
        const mac: Array<string> = [];

        for (let i = 0; i < 6; i++) {
            mac.push(this.data[i]!.toString(16).padStart(2, "0"));
        }

        this.data = this.data.subarray(6);

        return mac.join(":");
    }
}
