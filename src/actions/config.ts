export const config = {
    async set(): Promise<void> {
        try {
            console.log("config set");
        } catch (error: any) {
            console.error(error);
            process.exit(1);
        }
    },
    async get(): Promise<void> {
        try {
            console.log("config get");
        } catch (error: any) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default config;
