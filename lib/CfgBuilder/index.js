import fs from "fs-extra";
import path from "path";
import StringBuilder from "string-builder";
import figlet from "figlet";

export class CfgBuilder {
    constructor(apw, cliOptions, config) {
        this.apw = apw;
        this.config = config;
        this.builder = new StringBuilder();
        this.#setupOptions(cliOptions, config);
    }

    #setupOptions(cliOptions, config) {
        this.options = {
            rootDirectory: path.resolve(cliOptions.rootDirectory),

            outputLogFile: cliOptions.outputLogFile || config.cfg.outputLogFile,
            outputLogFileOption:
                cliOptions.outputLogFileOption ||
                config.cfg.outputLogFileOption,
            outputLogConsoleOption:
                cliOptions.outputLogConsoleOption ||
                config.cfg.outputLogConsoleOption,

            buildWithDebugInformation:
                cliOptions.buildWithDebugInformation ||
                config.cfg.buildWithDebugInformation,
            buildWithSource:
                cliOptions.buildWithSource || config.cfg.buildWithSource,

            includePath: cliOptions.includePath
                ? [...config.cfg.includePath, ...cliOptions.includePath]
                : config.cfg.includePath,

            modulePath: cliOptions.modulePath
                ? [...config.cfg.modulePath, ...cliOptions.modulePath]
                : config.cfg.modulePath,

            libraryPath: cliOptions.libraryPath
                ? [...config.cfg.libraryPath, ...cliOptions.libraryPath]
                : config.cfg.libraryPath,
        };

        return this;
    }

    #convertToComment(comment) {
        const lines = comment.split("\n");
        return lines.map((line) => (line ? `; ${line}` : ";")).join("\n");
    }

    #addBanner(builder, banner) {
        const figletBanner = figlet.textSync(banner);
        builder.appendLine(this.#convertToComment(figletBanner));
        return this;
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

        return this;
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

        return this;
    }

    #addMainRootDirectory(builder, options) {
        // prettier-ignore
        this.#addHeader(builder, [
            ";",
            ";  Main AXS Root Directory Reference",
            ";"
        ]);

        builder
            .appendLine(
                `MainAXSRootDirectory=${
                    options.rootDirectory === "." ? "-R" : options.rootDirectory
                }`,
            )
            .appendLine();

        return this;
    }

    #addAxsFiles(builder, files) {
        for (const file of files) {
            builder.appendLine(`AXSFILE=${file.path}`);
        }

        builder.appendLine();
        return this;
    }

    #addSourceFiles(builder) {
        this.#addHeader(builder, [
            ";",
            "; AXS files when specifying the MainAXSRootDirectory key above. You can have more",
            "; than one, order of the compile is as written.",
            ";",
        ])
            .#addAxsFiles(builder, this.apw.moduleFiles)
            .#addAxsFiles(builder, this.apw.masterSrcFiles);

        return this;
    }

    #addLogFileOptions(builder, options) {
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

        builder
            .appendLine(`OutputLogFile=${options.outputLogFile}`)
            .appendLine(`OutputLogFileOption=${options.outputLogFileOption}`)
            .appendLine(
                `OutputLogConsoleOption=${
                    options.outputLogConsoleOption ? "Y" : "N"
                }`,
            )
            .appendLine();

        return this;
    }

    #addCompilerOptionOverrides(builder, options) {
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
            `BuildWithDebugInformation=${
                options.buildWithDebugInformation ? "Y" : "N"
            }`,
        );
        builder.appendLine(
            `BuildWithSource=${options.buildWithSource ? "Y" : "N"}`,
        );

        builder.appendLine();
        return this;
    }

    #addAdditionalIncludePaths(builder, options) {
        if (options.includePath) {
            for (const includePath of options.includePath) {
                builder.appendLine(`AdditionalIncludePath=${includePath}`);
            }
        }

        if (this.apw.includePath) {
            for (const includePath of this.apw.includePath) {
                builder.appendLine(`AdditionalIncludePath=${includePath}`);
            }
        }

        builder.appendLine();
        return this;
    }

    #addAdditionalModulePaths(builder, options) {
        if (options.modulePath) {
            for (const modulePath of options.modulePath) {
                builder.appendLine(`AdditionalModulePath=${modulePath}`);
            }
        }

        builder.appendLine();
        return this;
    }

    #addAdditionalLibraryPaths(builder, options) {
        if (options.libraryPath) {
            for (const libraryPath of options.libraryPath) {
                builder.appendLine(`AdditionalLibraryPath=${libraryPath}`);
            }
        }

        builder.appendLine();
        return this;
    }

    #addAdditionalPaths(builder, options) {
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
        ])
            .#addAdditionalIncludePaths(builder, options)
            .#addAdditionalModulePaths(builder, options)
            .#addAdditionalLibraryPaths(builder, options);

        return this;
    }

    build() {
        const { builder, options } = this;

        this.#addBanner(builder, "NorgateAV");
        builder.appendLine(";");

        this.#addLicense(builder);
        builder
            .appendLine(
                "; This file was generated by genlinx, https://github.com/Norgate-AV-Solutions-Ltd/genlinx",
            )
            .appendLine(";");

        this.#addHeader(builder, [
            "; Used by the NetLinx Compiler Console program (NLRC.EXE) that specifies",
            "; how to invoke the the NetLinx Compiler with a configuration file via a",
            "; command console window.",
            ";",
            ';   > NLRC -C"C:\\AMX Projects\\build.cfg"',
            ";",
        ])
            .#addMainRootDirectory(builder, options)
            .#addSourceFiles(builder)
            .#addLogFileOptions(builder, options)
            .#addCompilerOptionOverrides(builder, options)
            .#addAdditionalPaths(builder, options);

        return builder.toString();
    }
}

export default CfgBuilder;
