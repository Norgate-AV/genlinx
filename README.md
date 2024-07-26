# genlinx

<div align="center">
    <img src="./assets/img/AMX_NS_03.png" alt="netlinx-studio-logo" width="150" />
</div>

---

[![CI](https://github.com/Norgate-AV/genlinx/actions/workflows/main.yml/badge.svg)](https://github.com/Norgate-AV/genlinx/actions)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)
[![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/)
[![GitHub contributors](https://img.shields.io/github/contributors/Norgate-AV/genlinx)](#contributors-sparkles)
[![NPM](https://img.shields.io/npm/v/%40norgate-av%2Fgenlinx)](https://www.npmjs.com/package/@norgate-av/genlinx)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

---

A CLI utility for Netlinx projects 🚀🚀🚀

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

## Contents 📖

-   [Installation :zap:](#installation-zap)
-   [Usage :rocket:](#usage-rocket)
-   [Command Line :man_technologist:](#command-line-man_technologist)
    -   [archive](#archive)
    -   [build](#build)
    -   [cfg](#cfg)
-   [Configuration :gear:](#configuration-gear)
    -   [Global](#global)
    -   [Local](#local)
        -   [Example](#example)
    -   [Command Line](#command-line)
    -   [Precedence](#precedence)
-   [Team :soccer:](#team-soccer)
-   [Contributors :sparkles:](#contributors-sparkles)
-   [LICENSE :balance_scale:](#license-balance_scale)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Installation :zap:

Install genlinx globally with:

```bash
pnpm add -g @norgate-av/genlinx

# or

yarn global add @norgate-av/genlinx

# or

npm install -g @norgate-av/genlinx
```

## Usage :rocket:

To run genlinx simply type:

```bash
genlinx
```

<div align="center">

</div>

<!-- ## Output :package:

genlinx will

-->

## Command Line :man_technologist:

```bash
                   _ _
   __ _  ___ _ __ | (_)_ __ __  __
  / _` |/ _ \ '_ \| | | '_ \\ \/ /
 | (_| |  __/ | | | | | | | |>  <
  \__, |\___|_| |_|_|_|_| |_/_/\_\
  |___/

Open source CLI tool for NetLinx projects
Copyright (c) 2022, Norgate AV Services Limited
https://github.com/Norgate-AV/genlinx

===================================================

Usage: genlinx [options] [command]

cli helper utility for NetLinx projects 🚀🚀🚀

Options:
  -v, --version      output the version number
  -h, --help         display help for command

Commands:
  archive [options]  generate a NetLinx workspace zip archive
  build [options]    build a NetLinx project
  cfg [options]      generate NetLinx build CFG files
  help [command]     display help for command

===================================================

For more help, make sure to check out the man page:
    $ man genlinx
```

### archive

Generate a NetLinx workspace zip archive.

```bash
Usage: genlinx archive [options]

generate a NetLinx workspace zip archive

Options:
  -w, --workspace-files <string...>              workspace file(s) to generate archive(s) for (default: search for workspace files in current directory)
  -o, --output-file-suffix <string>              output file suffix
  -S, --include-compiled-source-files            include compiled source files
  -M, --include-compiled-module-files            include compiled module files
  -N, --include-files-not-in-workspace           include files not in workspace
  -l, --extra-file-search-locations <string...>  extra file locations to search
  -p, --extra-file-archive-location <string>     location to place extra files in the archive
  --verbose                                      verbose output
  -h, --help                                     display help for command
```

### build

Build a NetLinx project from a CFG file

```bash
Usage: genlinx build [options]

build a NetLinx project

Options:
  -c, --cfg-files <string...>     cfg file(s) to build from (default: search for CFG files in current directory)
  -s, --source-files <string...>  axs source file(s) to build
  -i, --include-path <string...>  add additional include paths
  -m, --module-path <string...>   add additional module paths
  -l, --library-path <string...>  add additional library paths
  -a, --all                       select all cfg files without prompting
  -A, --no-all                    select multiple cfg files with a prompt
  --verbose                       verbose output
  -h, --help                      display help for command
```

### cfg

Generate a NetLinx build CFG file

```bash
Usage: genlinx cfg [options]

generate NetLinx build CFG files

Options:
  -w, --workspace-files <string...>      workspace file(s) to generate a CFG for (default: search for workspace files in current directory)
  -r, --root-directory <string>          root directory reference (default: use current directory as root)
  -o, --output-file-suffix <string>      output file suffix
  -f, --output-log-file-suffix <string>  output log file suffix
  -k, --output-log-file-option <string>  output log file option
                                         A - append
                                         N - overwrite (choices: "A", "N")
  -c, --output-log-console-option        output log to console
  -C, --no-output-log-console-option     do not output log to console
  -d, --build-with-debug-information     build with debug information
  -D, --no-build-with-debug-information  do not build with debug information
  -s, --build-with-source                build with source
  -S, --no-build-with-source             do not build with source
  -i, --include-path <string...>         add additional include paths
  -m, --module-path <string...>          add additional module paths
  -l, --library-path <string...>         add additional library paths
  -a, --all                              if no workspace files are specified with the -w option and more than one
                                         workspace file is found in the current directory, select all of them
                                         without prompting
  -A, --no-all                           if no workspace files are specified with the -w option and more than one
                                         workspace file is found in the current directory, prompt to select which
                                         workspace files to use
  --verbose                              verbose output
  -h, --help                             display help for command
```

<!-- ### config

Edit configuration properties for genlinx

```bash
Usage: genlinx config [options] [command]

edit configuration properties for genlinx

Options:
  -g, --global                               edit the global configuration
  -l, --local                                edit the local configuration
  -h, --help                                 display help for command

Commands:
  set <key <string>> <value(s) <string...>>  set configuration properties for genlinx
  get <key <string>>                         get configuration properties for genlinx
  help [command]                             display help for command
``` -->

## Configuration :gear:

### Global

A global configuration file can be installed to `%USERPROFILE%\.config\genlinx\config.json` to assist in using genlinx. This file can be edited to change the default values for the commands options.

If you would like to use a different location to store the global configuration file, you can set the `GENLINX_CONFIG_DIR` environment variable to the path of the directory you would like to use. This directory must contain the `config.json` file.

```json
{
    "cfg": {
        "outputFile": "build.cfg",
        "outputLogFile": "build.log",
        "outputLogFileOption": "N",
        "outputLogConsoleOption": true,
        "buildWithDebugInformation": false,
        "buildWithSource": false,
        "includePath": [
            "C:\\Program Files (x86)\\Common Files\\AMXShare\\AXIs",
            "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18\\includes"
        ],
        "modulePath": [
            "C:\\Program Files (x86)\\Common Files\\AMXShare\\Duet\\bundle",
            "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18",
            "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18\\monitors",
            "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18\\monitors-duet",
            "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18\\monitors-netlinx"
        ],
        "libraryPath": [
            "C:\\Program Files (x86)\\Common Files\\AMXShare\\SYCs"
        ],
        "all": false
    },
    "archive": {
        "outputFile": "archive.zip",
        "includeCompiledSourceFiles": true,
        "includeCompiledModuleFiles": true,
        "includeFilesNotInWorkspace": true,
        "extraFileSearchLocations": [
            "C:\\Program Files (x86)\\Common Files\\AMXShare",
            "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18"
        ],
        "extraFileArchiveLocation": ".genlinx",
        "all": false,
        "ignoredFiles": [
            "G4API.axi",
            "NetLinx.axi",
            "SNAPI.axi",
            "UnicodeLib.axi",
            "componentssdk.jar",
            "componentssdkrt.jar",
            "DeviceDriverEngine.jar",
            "devicesdkrt.jar",
            "jregex1.2_01-bundle.jar",
            "js-14-bundle.jar",
            "json-bundle.jar",
            "picocontainer-1.3-bundle.jar",
            "snapirouter.jar",
            "snapirouter2.jar"
        ]
    },
    "build": {
        "nlrc": {
            "path": "C:\\Program Files (x86)\\Common Files\\AMXShare\\COM\\NLRC.exe",
            "option": {
                "cfg": "-CFG",
                "includePath": "-I",
                "modulePath": "-M",
                "libraryPath": "-L"
            },
            "includePath": [
                "C:\\Program Files (x86)\\Common Files\\AMXShare\\AXIs",
                "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18\\includes"
            ],
            "modulePath": [
                "C:\\Program Files (x86)\\Common Files\\AMXShare\\Duet\\bundle",
                "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18",
                "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18\\monitors",
                "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18\\monitors-duet",
                "C:\\Program Files\\AMX\\Resource Management Suite\\SDK\\NetLinx\\4.7.18\\monitors-netlinx"
            ],
            "libraryPath": [
                "C:\\Program Files (x86)\\Common Files\\AMXShare\\SYCs"
            ]
        },
        "shell": {
            "path": "C:\\Windows\\System32\\cmd.exe"
        },
        "all": false,
        "createCfg": true
    }
}
```

### Local

A local configuration file can be used to override parts of global configuration file. This file should be placed in the same directory as the `apw` file and should be named `.genlinxrc.json`. The local configuration file will override only the parts you define. You don't have to redefine the entire configuration.

#### Example

To override the global configuration and not include compiled source or module files in the archive for a particular project, create a `.genlinxrc.json` file in the same directory as the `apw` file and add the following:

```json
{
    "archive": {
        "includeCompiledSourceFiles": false,
        "includeCompiledModuleFiles": false
    }
}
```

### Command Line

Options passed by via the CLI have the highest precedence and will override any configuration file options. However, additional paths passed via the CLI will be appended to the configuration file paths.

### Precedence

The precedence of the configuration options is as follows:

1. Command Line
2. Local Configuration File
3. Global Configuration File

<!-- ## Run with Docker :whale:

If you don't want to install nodejs or any node packages, use this method to run the genlinx from within a Docker container.

```bash
docker run -it -rm -v $(pwd):/usr/src/app genlinx:latest
````

> or

You can download this bash script from [here](./bin/genlinx) which wraps the above command into a simple command.

````bash
genlinx
``` -->

## Team :soccer:

This project is maintained by the following person(s) and a bunch of [awesome contributors](https://github.com/Norgate-AV/genlinx/graphs/contributors).

<table>
  <tr>
    <td align="center"><a href="https://github.com/damienbutt"><img src="https://avatars.githubusercontent.com/damienbutt?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Damien Butt</b></sub></a><br /></td>
  </tr>
</table>

## Contributors :sparkles:

<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->

[![All Contributors](https://img.shields.io/badge/all_contributors-2-orange.svg?style=flat-square)]

<!-- ALL-CONTRIBUTORS-BADGE:END -->

Thanks go to these awesome people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="https://github.com/features/security"><img src="https://avatars.githubusercontent.com/u/27347476?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Dependabot</b></sub></a><br /><a href="#maintenance-dependabot" title="Maintenance">🚧</a></td>
    <td align="center"><a href="https://imgbot.net/"><img src="https://avatars.githubusercontent.com/u/80986210?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Imgbot</b></sub></a><br /><a href="#maintenance-imgbot" title="Maintenance">🚧</a></td>
  </tr>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://allcontributors.org) specification.
Contributions of any kind are welcome!

Check out the [contributing guide](CONTRIBUTING.md) for more information.

## LICENSE :balance_scale:

[MIT](LICENSE)
