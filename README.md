# genlinx

<div align="center">
    <img src="./assets/img/AMX_NS_03.png" alt="netlinx-studio-logo" width="150" />
</div>

---

[![CI][ci]](https://github.com/Norgate-AV/genlinx/actions/workflows/ci.yml)
[![GitHub Release][release]](https://github.com/Norgate-AV/genlinx/releases)
[![npm][npm]](https://www.npmjs.com/package/@norgate-av/genlinx)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)
[![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/)
[![GitHub contributors](https://img.shields.io/github/contributors/Norgate-AV/genlinx)](#contributors-sparkles)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

---

A CLI utility for Netlinx projects 🚀🚀🚀

[ci]: https://img.shields.io/github/actions/workflow/status/Norgate-AV/genlinx/ci.yml?logo=github&label=CI
[release]: https://img.shields.io/github/v/release/Norgate-AV/genlinx?logo=github&label=Release
[npm]: https://img.shields.io/npm/v/@norgate-av/genlinx?logo=npm

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

## Contents 📖

- [Installation :zap:](#installation-zap)
- [Usage :rocket:](#usage-rocket)
- [Command Line :man_technologist:](#command-line-man_technologist)
    - [archive](#archive)
    - [build](#build)
    - [cfg](#cfg)
    - [config](#config)
    - [find](#find)
    - [ftl](#ftl)
    - [studio](#studio)
- [Configuration :gear:](#configuration-gear)
    - [Default](#default)
    - [Global](#global)
    - [Local](#local)
        - [Example](#example)
    - [Command Line](#command-line)
    - [Precedence](#precedence)
- [Migrating from v2 :arrows_counterclockwise:](#migrating-from-v2-arrows_counterclockwise)
- [Team :soccer:](#team-soccer)
- [Contributors :sparkles:](#contributors-sparkles)
- [LICENSE :balance_scale:](#license-balance_scale)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Installation :zap:

### Go install

```bash
go install github.com/Norgate-AV/genlinx@latest
```

### npm

```bash
npm install -g @norgate-av/genlinx
```

### Direct download

Pre-built binaries for Linux, macOS, and Windows are available on the [GitHub Releases](https://github.com/Norgate-AV/genlinx/releases) page.

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
  / _` |/ _ \ '_ \| | | '_ \/ /
 | (_| |  __/ | | | | | | | |>  <
  \__, |\___|_| |_|_|_|_| |_/_/\_\
  |___/

Open source CLI tool for NetLinx projects
Copyright (c) 2026, Norgate AV
https://github.com/Norgate-AV/genlinx

===================================================

Usage:
  genlinx [command]

Available Commands:
  archive     generate a NetLinx workspace zip archive
  build       build NetLinx source file(s) or CFG file(s)
  cfg         generate NetLinx build CFG files
  completion  Generate the autocompletion script for the specified shell
  config      view/edit configuration properties for genlinx
  find        find NetLinx devices on a local broadcast subnet
  ftl         generate a File Transfer List (.ftl) from a NetLinx workspace
  help        Help about any command
  studio      NetLinx Studio utilities

Flags:
  -h, --help      help for genlinx
      --verbose   verbose output
  -v, --version   version for genlinx

Use "genlinx [command] --help" for more information about a command.

===================================================

For more help, make sure to check out the man page:
    $ man genlinx
```

### archive

Generate a NetLinx workspace zip archive.

```bash
Usage:
  genlinx archive [flags]

Flags:
  -a, --all                                   process all found workspace files without prompting
  -l, --extra-file-search-locations strings   extra file locations to search
  -g, --extra-glob-patterns strings           glob patterns for extra files/directories to include in the archive (supports **)
  -h, --help                                  help for archive
  -m, --include-compiled-module-files         include compiled module files
  -s, --include-compiled-source-files         include compiled source files
  -n, --include-files-not-in-workspace        include files not in workspace
  -M, --no-include-compiled-module-files      do not include compiled module files
  -S, --no-include-compiled-source-files      do not include compiled source files
  -N, --no-include-files-not-in-workspace     do not include files not in workspace
  -o, --output-file-suffix string             output file suffix
  -p, --project string                        restrict archive to the named project within the workspace
  -y, --system string                         restrict archive to a single named system
  -w, --workspace-files strings               workspace file(s) to generate archive(s) for

Global Flags:
      --verbose   verbose output
```

### build

Build NetLinx source file(s) or CFG file(s).

```bash
Usage:
  genlinx build [sourceFiles...] [flags]

Flags:
  -a, --all                    select all cfg files without prompting
  -c, --cfg-files strings      cfg file(s) to build from
  -h, --help                   help for build
  -i, --include-path strings   add additional include paths
  -l, --library-path strings   add additional library paths
  -m, --module-path strings    add additional module paths
  -o, --output-path string     set the output path for the compiled files
  -s, --source-files strings   source file(s) to build (*.axs, *.axi)

Global Flags:
      --verbose   verbose output
```

### cfg

Generate a NetLinx build CFG file.

```bash
Usage:
  genlinx cfg [flags]

Flags:
  -a, --all                               process all found workspace files without prompting
  -d, --build-with-debug-information      build with debug information
  -s, --build-with-source                 build with source
  -h, --help                              help for cfg
  -i, --include-path strings              add additional include paths
  -l, --library-path strings              add additional library paths
  -m, --module-path strings               add additional module paths
  -D, --no-build-with-debug-information   do not build with debug information
  -S, --no-build-with-source              do not build with source
  -C, --no-output-log-console-option      do not output log to console
  -o, --output-file-suffix string         output file suffix
  -c, --output-log-console-option         output log to console
  -k, --output-log-file-option string     output log file option (A = append, N = overwrite)
  -f, --output-log-file-suffix string     output log file suffix
  -r, --root-directory string             root directory reference (default: parent directory of workspace file)
  -w, --workspace-files strings           workspace file(s) to generate a CFG for (default: search current directory)

Global Flags:
      --verbose   verbose output
```

### config

View/edit configuration properties for genlinx.

```bash
Usage:
  genlinx config [key] [flags]

Flags:
  -e, --edit     edit the configuration with default text editor
      --global   use global configuration
  -h, --help     help for config
  -l, --list     display the configuration in stdout
      --local    use local configuration

Global Flags:
      --verbose   verbose output
```

### find

Find NetLinx devices on a local broadcast subnet.

```bash
Usage:
  genlinx find [flags]

Flags:
  -h, --help   help for find
  -j, --json   output as JSON

Global Flags:
      --verbose   verbose output
```

### ftl

Generate a File Transfer List (`.ftl`) from a NetLinx workspace. The output can be imported directly into NetLinx Studio or the AMX File Transfer Tool.

```bash
Usage:
  genlinx ftl [workspaceFiles...] [flags]

Flags:
  -a, --all                       process all discovered workspace files without prompting
  -h, --help                      help for ftl
  -o, --output string             output filename (default: <workspace-name>.ftl)
  -p, --project string            include only systems from this project
  -s, --system string             include only this system
  -w, --workspace-files strings   workspace file(s) to process

Global Flags:
      --verbose   verbose output
```

### studio

NetLinx Studio utilities (Windows only). Reads and writes settings directly via the Windows registry.

```bash
Usage:
  genlinx studio [command]

Available Commands:
  backup    backup NetLinx Studio settings to an .epx file
  restore   restore NetLinx Studio settings from an .epx file

Flags:
  -h, --help   help for studio

Global Flags:
      --verbose   verbose output
```

## Configuration :gear:

### Default

```json
{
    "cfg": {
        "outputFile": "build.cfg",
        "outputLogFile": "build.log",
        "outputLogFileOption": "N",
        "outputLogConsoleOption": true,
        "buildWithDebugInformation": false,
        "buildWithSource": false
    },
    "archive": {
        "outputFile": "archive.zip",
        "includeCompiledSourceFiles": true,
        "includeCompiledModuleFiles": true,
        "includeFilesNotInWorkspace": true,
        "extraFileSearchLocations": [
            "C:/Program Files (x86)/Common Files/AMXShare"
        ],
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
    "compiler": {
        "path": "C:/Program Files (x86)/Common Files/AMXShare/COM/NLRC.exe",
        "includePath": ["C:/Program Files (x86)/Common Files/AMXShare/AXIs"],
        "modulePath": [
            "C:/Program Files (x86)/Common Files/AMXShare/Duet/bundle",
            "C:/Program Files (x86)/Common Files/AMXShare/Duet/lib",
            "C:/Program Files (x86)/Common Files/AMXShare/Duet/module"
        ],
        "libraryPath": ["C:/Program Files (x86)/Common Files/AMXShare/SYCs"]
    }
}
```

### Global

A global configuration file can be installed to `%USERPROFILE%\.config\genlinx\config.json` to assist in using genlinx. This file can be edited to change the default values for the commands options.

If you would like to use a different location to store the global configuration file, you can set the `GENLINX_CONFIG_DIR` environment variable to the path of the directory you would like to use. This directory must contain the `config.json` file.

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
4. Default Configuration

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

## Migrating from v2 :arrows_counterclockwise:

The v3 (Go) rewrite reorganises the configuration file. The config file name (`.genlinxrc.json`) is unchanged.

### Renamed keys

| Section   | v2 key                | v3 key          |
| --------- | --------------------- | --------------- |
| `cfg`     | `outputFileSuffix`    | `outputFile`    |
| `cfg`     | `outputLogFileSuffix` | `outputLogFile` |
| `archive` | `outputFileSuffix`    | `outputFile`    |

### New `compiler` section

The `build.nlrc` block has been promoted to a top-level `compiler` section, and the per-command path overrides inside `cfg` have moved there too:

```diff
-  "build": {
-    "nlrc": {
-      "path": "C:/Program Files (x86)/Common Files/AMXShare/COM/NLRC.exe",
-      "includePath": [...],
-      "modulePath": [...],
-      "libraryPath": [...]
-    }
-  },
-  "cfg": {
-    "includePath": [...],
-    "modulePath": [...],
-    "libraryPath": [...]
-  }
+  "compiler": {
+    "path": "C:/Program Files (x86)/Common Files/AMXShare/COM/NLRC.exe",
+    "includePath": [...],
+    "modulePath": [...],
+    "libraryPath": [...]
+  }
```

### Removed keys

The following keys are no longer valid in the config file. Where applicable, the equivalent behaviour is now a CLI flag.

| Section   | v2 key                     | Replacement      |
| --------- | -------------------------- | ---------------- |
| `build`   | `shell.path`               | removed          |
| `build`   | `all`                      | `--all` CLI flag |
| `build`   | `createCfg`                | `--cfg` CLI flag |
| `cfg`     | `all`                      | `--all` CLI flag |
| `archive` | `all`                      | `--all` CLI flag |
| `archive` | `extraFileArchiveLocation` | removed          |

### New keys

| Section   | v3 key              | Description                                             |
| --------- | ------------------- | ------------------------------------------------------- |
| `archive` | `extraGlobPatterns` | Glob patterns for extra files to include in the archive |

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
