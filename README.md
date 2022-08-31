# genlinx

<div align="center">
    <img src="./assets/img/AMX_NS_03.png" alt="netlinx-studio-logo" width="150" />
</div>

---

[![CI](https://github.com/Norgate-AV-Solutions-Ltd/genlinx/actions/workflows/main.yml/badge.svg)](https://github.com/Norgate-AV-Solutions-Ltd/genlinx/actions)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)
[![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/)
[![GitHub contributors](https://img.shields.io/github/contributors/Norgate-AV-Solutions-Ltd/genlinx)](#contributors)
[![NPM](https://img.shields.io/npm/v/genlinx.svg)](https://www.npmjs.com/package/genlinx)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

---

A CLI utility for Netlinx projects 🚀🚀🚀

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

## Contents 📖

-   [Installation :zap:](#installation-zap)
-   [Usage :rocket:](#usage-rocket)
-   [Command Line :man_technologist:](#command-line-man_technologist)
    -   [cfg](#cfg)
    -   [archive](#archive)
    -   [build](#build)
-   [Configuration :gear:](#configuration-gear)
-   [Team :soccer:](#team-soccer)
-   [Contributors :sparkles:](#contributors-sparkles)
-   [LICENSE :balance_scale:](#license-balance_scale)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Installation :zap:

Install genlinx globally with:

```bash
npm install -g genlinx

# or

yarn global add genlinx
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
Usage: genlinx [options] [command]

cli helper utility for NetLinx projects 🚀🚀🚀

Options:
  -v, --version                          output the version number
  -h, --help                             display help for command

Commands:
  cfg [options] <apw file <string>>      generate a NetLinx build CFG file
  archive [options] <apw file <string>>  generate a NetLinx workspace zip archive
  build <cfg file <string>>              build a NetLinx project from a CFG file
  help [command]                         display help for command
```

### cfg

```bash
Usage: genlinx cfg [options] <apw file <string>>

generate a NetLinx build CFG file

Arguments:
  apw file <string>                      apw file to generate the CFG from

Options:
  -d, --root-directory <string>          root directory reference (default: ".")
  -o, --output-file <string>             output file name
  -L, --output-log-file <string>         output log file name
  -O, --output-log-file-option <string>  output log file option (choices: "A", "N")
  -C, --output-log-console-option        output log to console
  -D, --build-with-debug-information     build with debug information
  -S, --build-with-source                build with source
  -i, --include-path <string...>         additional include paths
  -m, --module-path <string...>          additional module paths
  -l, --library-path <string...>         additional library paths
  -h, --help                             display help for command
```

### archive

```bash
Usage: genlinx archive [options] <apw file <string>>

generate a NetLinx workspace zip archive

Arguments:
  apw file <string>           apw file to generate the archive from

Options:
  -o, --output-file <string>  output file name
  -h, --help                  display help for command
```

### build

```bash
Usage: genlinx build [options] <cfg file <string>>

build a NetLinx project from a CFG file

Arguments:
  cfg file <string>  cfg file to build from

Options:
  -h, --help         display help for command
```

## Configuration :gear:

A configuration file will be installed to `%USERPROFILE%\.config\genlinx\config.json` to assist in using genlinx.

Options passed by via the CLI will override the configuration file where applicable. However, additional include, module and library file paths passed via the CLI in relation to the `cfg` command will be appended to the corresponding path lists from the configuration file.

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
        "libraryPath": ["C:\\Program Files (x86)\\Common Files\\AMXShare\\SYCs"]
    },
    "archive": {
        "outputFile": "archive.zip"
    },
    "build": {
        "nlrc": {
            "path": "C:\\Program Files (x86)\\Common Files\\AMXShare\\COM\\NLRC.exe",
            "option": {
                "cfg": "-CFG"
            }
        },
        "shell": {
            "path": "C:\\Windows\\System32\\cmd.exe"
        }
    }
}
```

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

This project is maintained by the following person(s) and a bunch of [awesome contributors](https://github.com/Norgate-AV-Solutions-Ltd/genlinx/graphs/contributors).

<table>
  <tr>
    <td align="center"><a href="https://github.com/damienbutt"><img src="https://avatars.githubusercontent.com/damienbutt?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Damien Butt</b></sub></a><br /></td>
  </tr>
</table>

## Contributors :sparkles:

<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->

[![All Contributors](https://img.shields.io/badge/all_contributors-1-orange.svg?style=flat-square)](#contributors-)

<!-- ALL-CONTRIBUTORS-BADGE:END -->

Thanks go to these awesome people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://allcontributors.org) specification.
Contributions of any kind are welcome!

Check out the [contributing guide](CONTRIBUTING.md) for more information.

## LICENSE :balance_scale:

[MIT](LICENSE)
