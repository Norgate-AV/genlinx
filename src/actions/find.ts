import dgram, { RemoteInfo } from "node:dgram";
import { Listr } from "listr2";
import * as jq from "node-jq";
import { Table } from "console-table-printer";
import { Device, Packet } from "../../lib/index.js";
import { FindCliArgs } from "../../lib/@types/index.js";

const ICSP_PORT = 1319;

export async function discover(): Promise<Array<Device>> {
    const devices: Array<Device> = [];

    const socket = dgram.createSocket("udp4");
    socket.bind(ICSP_PORT);

    socket.on("message", (message: Buffer, remoteInfo: RemoteInfo) => {
        let packet: Packet;

        try {
            packet = new Packet(message);
        } catch (error: any) {
            return;
        }

        if (packet.getType() !== "A") {
            return;
        }

        devices.push({
            ...packet.parse().getDevice(),
            ip: remoteInfo.address,
        });
    });

    return new Promise((resolve) => {
        setTimeout(() => {
            socket.close();

            resolve([
                ...new Map(
                    devices.map((device) => [device.mac, device]),
                ).values(),
            ]);
        }, 6000);
    });
}

interface Context {
    devices: Array<Device>;
}

export const find = {
    async execute(args: FindCliArgs) {
        try {
            const tasks = new Listr<Context>([
                {
                    title: "Listening for NetLinx devices...",
                    task: async (context): Promise<void> => {
                        context.devices = await discover();

                        const { devices } = context;

                        if (devices.length === 0) {
                            console.log("No devices found");
                            return;
                        }

                        if (args.json) {
                            console.log(
                                await jq.run(".", JSON.stringify(devices), {
                                    color: true,
                                    input: "json",
                                }),
                            );

                            return;
                        }

                        const table = new Table({
                            title: "Discovered Devices",
                            enabledColumns: [
                                "ip",
                                "mac",
                                "date",
                                "time",
                                "system",
                                "hostname",
                                "id",
                            ],
                            columns: [
                                { name: "ip", title: "IP Address" },
                                { name: "system", title: "System" },
                                { name: "date", title: "Date" },
                                { name: "time", title: "Time" },
                                { name: "day", title: "Day" },
                                { name: "mac", title: "MAC Address" },
                                { name: "hostname", title: "Hostname" },
                                { name: "id", title: "ID" },
                            ],
                            sort: (a, b) => a.ip.localeCompare(b.ip),
                        });

                        for (const device of devices) {
                            table.addRow(
                                {
                                    ip: device.ip,
                                    system: device.system,
                                    date: device.date.text,
                                    time: device.time,
                                    day: device.day,
                                    mac: device.mac,
                                    hostname: device.hostname,
                                    id: device.id,
                                },
                                { color: "green" },
                            );
                        }

                        table.printTable();
                    },
                },
            ]);

            await tasks.run();
        } catch (error: any) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default find;
