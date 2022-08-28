import fs from "fs-extra";
import path from "path";
import StringBuilder from "string-builder";
import figlet from "figlet";

export class CfgBuilder {
    constructor(apw, options, config) {
        this.apw = apw;
        this.options = options;
        this.builder = new StringBuilder();
        this.config = config;
    }

    #convertToComment(comment) {
        const lines = comment.split("\n");
        return lines.map((line) => (line ? `; ${line}` : ";")).join("\n");
    }

    #addBanner(builder, banner) {
        const figletBanner = figlet.textSync(banner);
        builder.appendLine(this.#convertToComment(figletBanner));
    }

    #addLicense(builder) {
        try {
            const license = fs.readFileSync(
                path.resolve(__dirname, "../../LICENSE"),
                { encoding: "utf8" },
            );

            builder.appendLine(this.#convertToComment(license));
        } catch (error) {
            console.error(error);
        }
    }

    #addHeader(builder, contentLines) {
        builder.appendLine(
            ";------------------------------------------------------------------------------",
        );

        contentLines.forEach((line) => {
            builder.appendLine(line);
        });

        builder.appendLine(
            ";------------------------------------------------------------------------------",
        );
    }

    #addModuleFiles(builder) {
        for (const file of this.apw.moduleFiles) {
            builder.appendLine(`AXSFILE=${file.path}`);
        }
    }

    #addMasterSrcFiles(builder) {
        for (const file of this.apw.masterSrcFiles) {
            builder.appendLine(`AXSFILE=${file.path}`);
        }
    }

    #addAdditionalIncludePaths(builder, options, config) {
        if (config.cfg.includePath) {
            for (const includePath of config.cfg.includePath) {
                builder.appendLine(`AdditionalIncludePath=${includePath}`);
            }
        }

        if (options.includePath) {
            for (const includePath of options.includePath) {
                builder.appendLine(`AdditionalIncludePath=${includePath}`);
            }
        }
    }

    #addAdditionalModulePaths(builder, options, config) {
        if (config.cfg.modulePath) {
            for (const modulePath of config.cfg.modulePath) {
                builder.appendLine(`AdditionalModulePath=${modulePath}`);
            }
        }

        if (options.modulePath) {
            for (const modulePath of options.modulePath) {
                builder.appendLine(`AdditionalModulePath=${modulePath}`);
            }
        }
    }

    #addAdditionalLibraryPaths(builder, options, config) {
        if (config.cfg.libraryPath) {
            for (const libraryPath of config.cfg.libraryPath) {
                builder.appendLine(`AdditionalLibraryPath=${libraryPath}`);
            }
        }

        if (options.libraryPath) {
            for (const libraryPath of options.libraryPath) {
                builder.appendLine(`AdditionalLibraryPath=${libraryPath}`);
            }
        }
    }

    build() {
        const { builder, options, config } = this;

        this.#addBanner(builder, "NorgateAV");
        builder.appendLine(";");
        this.#addLicense(builder);
        builder.appendLine(
            "; This file was generated by genlinx, https://github.com/Norgate-AV-Solutions-Ltd/genlinx",
        );

        builder.appendLine(";");
        this.#addHeader(builder, [
            "; Used by the NetLinx Compiler Console program (NLRC.EXE) that specifies",
            "; how to invoke the the NetLinx Compiler with a configuration file via a",
            "; command console window.",
            ";",
            ';   > NLRC -C"C:\\AMX Projects\\build.cfg"',
            ";",
        ]);

        // prettier-ignore
        this.#addHeader(builder, [
            ";",
            ";  Main AXS Root Directory Reference",
            ";"
        ]);

        builder.appendLine(
            `MainAXSRootDirectory=${options.dir === "." ? "-R" : options.dir}`,
        );

        builder.appendLine();
        this.#addHeader(builder, [
            ";",
            "; AXS files when specifying the MainAXSRootDirectory key above. You can have more",
            "; than one, order of the compile is as written.",
            ";",
        ]);

        this.#addModuleFiles(builder);
        builder.appendLine();
        this.#addMasterSrcFiles(builder);

        builder.appendLine();
        this.#addHeader(builder, [
            ";",
            "; Output Log File and Log File Options.",
            ";",
            "; OutputLogFile=        <--: Output log file name",
            ";",
            ";    Fully qualified file name (no quotes are needed)",
            ";    If no OutputLogFile key present, then by default, log to the console",
            ";    window.  Unless the OptionLogConsoleOptions= is specified (see below).",
            ";",
            "; OutputLogFileOption=  <--: Output log file option",
            ";",
            ";    A = Append status to the output file. If file does not exist,",
            ";        then the program will create a new one.",
            ";    N = Create a new output file. Overwrites if the file already exists.",
            ";",
            ";  If no OutputLogFileOption key present, then the default is N.",
            ";",
            "; OutputLogConsoleOption= <--: Output Log to the Console",
            ";",
            ";    Y = Send log info to the console.",
            ";    N = Do no send log info to the console.",
        ]);

        builder.appendLine(`OutputLogFile=${options.logFile}`);
        builder.appendLine(`OutputLogFileOption=${options.logFileOption}`);
        builder.appendLine(
            `OutputLogConsoleOption=${options.logConsoleOption ? "Y" : "N"}`,
        );

        builder.appendLine();
        this.#addHeader(builder, [
            ";",
            "; NetLinx Compiler Option Overrides",
            ";",
            ";   Ability to override the NetLinx Studio Compiler options that are defined",
            ";   within NetLinx Studio.",
            ";",
            ";   Y = Yes   N = No",
            ";",
            "; Comment these options out if you want to use the NetLinx Studio options.",
        ]);

        builder.appendLine(
            `BuildWithDebugInformation=${options.buildWithDebug ? "Y" : "N"}`,
        );
        builder.appendLine(
            `BuildWithSource=${options.buildWithSource ? "Y" : "N"}`,
        );

        builder.appendLine();
        this.#addHeader(builder, [
            "; Additional Paths",
            ";",
            "; If you need to specify additional paths for the NetLinx compiler, you can add",
            "; the following keys:",
            ";",
            ";    AdditionalIncludePath=",
            ";    AdditionalLibraryPath=",
            ";    AdditionalModulePath=",
            ";",
            "; You can specify upto 50 additional paths for each type (one directory per",
            "; key upto 50 keys per type).  No quotes are needed for the directory names.",
        ]);

        this.#addAdditionalIncludePaths(builder, options, config);
        builder.appendLine();
        this.#addAdditionalModulePaths(builder, options, config);
        builder.appendLine();
        this.#addAdditionalLibraryPaths(builder, options, config);
        builder.appendLine();

        return builder.toString();
    }
}

export default CfgBuilder;