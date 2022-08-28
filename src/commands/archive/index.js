export const archive = {
    create(filePath, options) {
        try {
            console.log("archive command");
            console.log(filePath);
            console.log(options);
        } catch (error) {
            console.error(error);
            process.exit(1);
        }
    },
};

export default archive;
