export const config = {
    async set() {
        try {
            console.log("config set");
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
    async get() {
        try {
            console.log("config get");
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default config;
