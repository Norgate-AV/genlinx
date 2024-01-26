import fs from "fs-extra";
import path from "path";
import StringBuilder from "string-builder";
import figlet from "figlet";
import APW from "../APW";
import { CfgConfig, CliCfgOptions, FileReference } from "../@types";

export class CfgBuilder {
    private readonly apw: APW;
    private readonly options: CfgConfig & CliCfgOptions;
    private readonly builder = new StringBuilder();

    public constructor(apw: APW, options: CfgConfig & CliCfgOptions) {
        this.apw = apw;
        this.options = options;
    }

    private convertToComment(comment: string): string {
        const lines = comment.split("\n");
        return lines.map((line) => (line ? `; ${line}` : ";")).join("\n");
    }

    private addBanner(banner: string): this {
        const { builder } = this;

        const figletBanner = figlet.textSync(banner);
        builder.appendLine(this.convertToComment(figletBanner));

        return this;
    }

    private addLicense(): this {
        const { builder } = this;

        try {
            const license = fs.readFileSync(
                path.resolve(__dirname, "../../LICENSE"),
                {
                    encoding: "utf8",
                },
            );

            builder.appendLine(this.convertToComment(license));
        } catch (error: any) {
            console.error(error);
        }

        return this;
    }

    private addHeader(contentLines: Array<string>): this {
        const { builder } = this;

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

    private addMainRootDirectory(): this {
        const { builder, options } = this;

        // prettier-ignore
        this.addHeader([
            ";",
            ";  Main AXS Root Directory Reference",
            ";"
        ]);

        const rootDirectory = options.rootDirectory
            ? path.resolve(options.rootDirectory)
            : "-R";

        builder
            .appendLine(`MainAXSRootDirectory=${rootDirectory}`)
            .appendLine();

        return this;
    }

    private addAxsFiles(files: Array<FileReference>): this {
        const { builder } = this;

        for (const file of files) {
            builder.appendLine(`AXSFILE=${file.path}`);
        }

        builder.appendLine();
        return this;
    }

    private addSourceFiles(): this {
        const { apw } = this;

        this.addHeader([
            ";",
            "; AXS files when specifying the MainAXSRootDirectory key above. You can have more",
            "; than one, order of the compile is as written.",
            ";",
        ])
            .addAxsFiles(apw.moduleFiles)
            .addAxsFiles(apw.masterSrcFiles);

        return this;
    }

    private addLogFileOptions(): this {
        const { builder, apw, options } = this;

        this.addHeader([
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

        const outputLogFile = `${apw.id}.${options.outputLogFileSuffix}`;

        builder
            .appendLine(`OutputLogFile=${outputLogFile}`)
            .appendLine(`OutputLogFileOption=${options.outputLogFileOption}`)
            .appendLine(
                `OutputLogConsoleOption=${
                    options.outputLogConsoleOption ? "Y" : "N"
                }`,
            )
            .appendLine();

        return this;
    }

    private addCompilerOptionOverrides(): this {
        const { builder, options } = this;

        this.addHeader([
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

    private addAdditionalIncludePaths(): this {
        const { builder, apw, options } = this;

        options.includePath.push(...apw.includePath);

        if (options.includePath) {
            for (const includePath of options.includePath) {
                builder.appendLine(`AdditionalIncludePath=${includePath}`);
            }
        }

        builder.appendLine();
        return this;
    }

    private addAdditionalModulePaths(): this {
        const { builder, apw, options } = this;

        options.modulePath.push(...apw.modulePath);

        if (options.modulePath) {
            for (const modulePath of options.modulePath) {
                builder.appendLine(`AdditionalModulePath=${modulePath}`);
            }
        }

        builder.appendLine();
        return this;
    }

    private addAdditionalLibraryPaths(): this {
        const { builder, options } = this;

        if (options.libraryPath) {
            for (const libraryPath of options.libraryPath) {
                builder.appendLine(`AdditionalLibraryPath=${libraryPath}`);
            }
        }

        builder.appendLine();
        return this;
    }

    private addAdditionalPaths(): this {
        this.addHeader([
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
            .addAdditionalIncludePaths()
            .addAdditionalModulePaths()
            .addAdditionalLibraryPaths();

        return this;
    }

    public build(): string {
        const { builder } = this;

        this.addBanner("NorgateAV");
        builder.appendLine(";");

        this.addLicense();
        builder
            .appendLine(
                "; This file was generated by genlinx, https://github.com/Norgate-AV/genlinx",
            )
            .appendLine(";");

        this.addHeader([
            "; Used by the NetLinx Compiler Console program (NLRC.EXE) that specifies",
            "; how to invoke the the NetLinx Compiler with a configuration file via a",
            "; command console window.",
            ";",
            ';   > NLRC -C"C:\\AMX Projects\\build.cfg"',
            ";",
        ])
            .addMainRootDirectory()
            .addSourceFiles()
            .addLogFileOptions()
            .addCompilerOptionOverrides()
            .addAdditionalPaths();

        return builder.toString();
    }
}

export default CfgBuilder;
