# Changelog

All notable changes to this project will be documented in this file.

## [3.0.0](https://github.com/Norgate-AV/genlinx/compare/v2.7.0..v3.0.0) - 2026-10-12

### 🚀 Features

- [**breaking**] Rewrite in Go - ([af150fe](https://github.com/Norgate-AV/genlinx/commit/af150feea22dc2b02e585cc4af8d9747e4f707fe))

## [2.7.0](https://github.com/Norgate-AV/genlinx/compare/v2.6.1..v2.7.0) - 2025-04-22

### 🚀 Features

- *(build)* Add argument for source files so no flag is required - ([e2b2c5d](https://github.com/Norgate-AV/genlinx/commit/e2b2c5d8ad9cf8f2c56d34a80df32cc4a4acbf1e))

## [2.6.1](https://github.com/Norgate-AV/genlinx/compare/v2.6.0..v2.6.1) - 2024-11-22

### 🐛 Bug Fixes

- *(find)* Correctly display json output in pretty format - ([96a8937](https://github.com/Norgate-AV/genlinx/commit/96a89376a23b4a27afc023bb8fa98c0868505b36))

## [2.6.0](https://github.com/Norgate-AV/genlinx/compare/v2.5.0..v2.6.0) - 2024-11-22

### 🚀 Features

- *(find)* Add support for finding muse controllers - ([72f7d55](https://github.com/Norgate-AV/genlinx/commit/72f7d559396519511194e058f4f3c2ad1cd07de1))

### 🐛 Bug Fixes

- *(config)* Handle exception if no global config is found - ([3f232ce](https://github.com/Norgate-AV/genlinx/commit/3f232ce9c36c9a6ebeb66a165afdd25132933687))

## [2.5.0](https://github.com/Norgate-AV/genlinx/compare/v2.4.0..v2.5.0) - 2024-08-14

### 🚀 Features

- *(find)* Add cli option for timeout with 6000ms as default - ([772988e](https://github.com/Norgate-AV/genlinx/commit/772988e4ee77f3f8f655012c4f8ca871a95247e3))

## [2.4.0](https://github.com/Norgate-AV/genlinx/compare/v2.3.0..v2.4.0) - 2024-08-02

### 🚀 Features

- Integrate findlinx with the find command - ([61770f2](https://github.com/Norgate-AV/genlinx/commit/61770f2a0e43b51ac34ceb53403f382e120040e2))

## [2.3.0](https://github.com/Norgate-AV/genlinx/compare/v2.2.3..v2.3.0) - 2024-08-02

### 🚀 Features

- *(archive-cmd)* Add negated boolean options - ([56c414e](https://github.com/Norgate-AV/genlinx/commit/56c414e78c751b9c58e8c1aaa11ca8f241129844))

- *(config)* Implement config get value - ([f03f136](https://github.com/Norgate-AV/genlinx/commit/f03f13695d7647af69d8d32700ea54655cdddcaf))

- *(config)* Add config schema for validation - ([f1c299e](https://github.com/Norgate-AV/genlinx/commit/f1c299e88a46c2a244ab2e9cdbae67c4c71b418d))

- *(config)* Add core section - ([1ea0109](https://github.com/Norgate-AV/genlinx/commit/1ea01096f368a8a3f55f08f57bb0fd958a82cd57))

- *(config)* Add --edit option to open config in default editor - ([4bc70ac](https://github.com/Norgate-AV/genlinx/commit/4bc70ac19c4c5e1229cd5b86395ca1a735c06320))

- *(config)* Add --delete option for deleting a key - ([0c35ed1](https://github.com/Norgate-AV/genlinx/commit/0c35ed176b8299e4c4930353a1a1b297a57473d4))

- *(config)* Add --add and --remove options for adding/removing from an array - ([a0ee5e6](https://github.com/Norgate-AV/genlinx/commit/a0ee5e640437da10f59645dac7162981bceb31f6))

- *(config)* Add schema for config - ([beb53f3](https://github.com/Norgate-AV/genlinx/commit/beb53f3f12d44b1bf50c1c06c54725f125e84895))

- *(config-cmd)* Add basic implementation of cmd functionality - ([1e9ac7a](https://github.com/Norgate-AV/genlinx/commit/1e9ac7a192e77873f13d1516f5ada46e66d71603))

- Run paths through normalizer so / or \\ can be used - ([1071b26](https://github.com/Norgate-AV/genlinx/commit/1071b260ebab2a0f280c668e9320e8ee4e5a6379))

- Implement config listing and editing in editor - ([cc144c3](https://github.com/Norgate-AV/genlinx/commit/cc144c3230fc957f348143f7731b570d5347070c))

- Add config command to main cli - ([4ab0192](https://github.com/Norgate-AV/genlinx/commit/4ab0192caffd20aadb1127dfab098cd53611069a))

- Export new util functions - ([011aff7](https://github.com/Norgate-AV/genlinx/commit/011aff7d84074093f6c4d551c78e5e7e1693feb0))

- Add getLocalAppConfigFilePath function - ([96c03e0](https://github.com/Norgate-AV/genlinx/commit/96c03e0e92b6b23eff1a0baffb39bb93c3c6bf11))

- Add resolvePaths function - ([e785cec](https://github.com/Norgate-AV/genlinx/commit/e785cec75760b49b946608fec69eebe33ed579be))

- Add selectWorkspaceFiles function - ([688053a](https://github.com/Norgate-AV/genlinx/commit/688053ac92237e71914437d563be58485e0a8e9b))

### 🐛 Bug Fixes

- *(config)* Return after getting property - ([7d24ee9](https://github.com/Norgate-AV/genlinx/commit/7d24ee944e090bc92d66c3417b27f98ad99d088e))

## [2.2.3](https://github.com/Norgate-AV/genlinx/compare/v2.2.2..v2.2.3) - 2024-07-27

### 🐛 Bug Fixes

- Prepend arrays in config and dont sort - ([86c8887](https://github.com/Norgate-AV/genlinx/commit/86c88872156b5a003562248a220ac897236e1476))

## [2.2.2](https://github.com/Norgate-AV/genlinx/compare/v2.2.1..v2.2.2) - 2024-07-24

### 🐛 Bug Fixes

- Use @inquirer/prompts for selecting files - ([4ad90cd](https://github.com/Norgate-AV/genlinx/commit/4ad90cd77da30a73e0b6dfe3dc4df5060e7f4e5d))

### 🚜 Refactor

- *(config)* [**breaking**] Remove RMS locations from the default config - ([7eb65e6](https://github.com/Norgate-AV/genlinx/commit/7eb65e6eeab8f92285b726be81dbd31e26cc4680))

## [2.2.1](https://github.com/Norgate-AV/genlinx/compare/v2.2.0..v2.2.1) - 2024-02-01

### 🐛 Bug Fixes

- *(archive)* Correctly determine file type or return "Other" - ([66e68df](https://github.com/Norgate-AV/genlinx/commit/66e68dfcb7202fbf5ba895e068acf0b09f4a8eee))

## [2.2.0](https://github.com/Norgate-AV/genlinx/compare/v2.1.1..v2.2.0) - 2024-02-01

### 🚀 Features

- *(build)* Allow for building multiple source files - ([f2f8243](https://github.com/Norgate-AV/genlinx/commit/f2f8243928ecdd5733ab0ea0b6924467c98693bc))

## [2.1.1](https://github.com/Norgate-AV/genlinx/compare/v2.1.0..v2.1.1) - 2024-02-01

### 🐛 Bug Fixes

- Only log using local config if verbose - ([f72744d](https://github.com/Norgate-AV/genlinx/commit/f72744de9cd7a7eca207caf37cc57cd63408965c))

## [2.1.0](https://github.com/Norgate-AV/genlinx/compare/v2.0.3..v2.1.0) - 2024-02-01

### 🚀 Features

- Add --verbose flag to commands - ([e822091](https://github.com/Norgate-AV/genlinx/commit/e82209108b50aa1d64aa4807a2f349ebc5d64df6))

## [2.0.3](https://github.com/Norgate-AV/genlinx/compare/v2.0.2..v2.0.3) - 2024-02-01

### 🐛 Bug Fixes

- Open correct file to obtain latest version - ([34617d2](https://github.com/Norgate-AV/genlinx/commit/34617d2167274b3ec84e2bee27b68e531be329d1))

- Add year to banner dynamically - ([0425c24](https://github.com/Norgate-AV/genlinx/commit/0425c242d0e35336b50730ff44e9d97adeebb22d))

## [2.0.2](https://github.com/Norgate-AV/genlinx/compare/v2.0.1..v2.0.2) - 2024-01-31

### 🐛 Bug Fixes

- Get most upto date app version - ([9c7c784](https://github.com/Norgate-AV/genlinx/commit/9c7c784555c799db47987aee1621f90faad29b1d))

## [2.0.1](https://github.com/Norgate-AV/genlinx/compare/v2.0.0..v2.0.1) - 2024-01-31

### 🐛 Bug Fixes

- Replace __dirname in code - ([9432445](https://github.com/Norgate-AV/genlinx/commit/9432445ad4d26ab3e19212f444a9fdf4388e0663))

## [2.0.0](https://github.com/Norgate-AV/genlinx/compare/v1.1.4..v2.0.0) - 2024-01-31

### 🚀 Features

- [**breaking**] Convert codebase to typescript - ([937de15](https://github.com/Norgate-AV/genlinx/commit/937de15f1a30b354d7f104785218865c2f7f58bf))

## [1.1.4](https://github.com/Norgate-AV/genlinx/compare/v1.1.3..v1.1.4) - 2024-01-04

### 🐛 Bug Fixes

- *(build)* Remove accidental console.log - ([a71b27e](https://github.com/Norgate-AV/genlinx/commit/a71b27e628f95d10d2fb1d8501a25fa71b3e313c))

## [1.1.3](https://github.com/Norgate-AV/genlinx/compare/v1.1.2..v1.1.3) - 2024-01-04

### 🐛 Bug Fixes

- *(build)* Filter out falsy props from options - ([e105abe](https://github.com/Norgate-AV/genlinx/commit/e105abe536ffe37c74ea299124835113ec11d067))

## [1.1.2](https://github.com/Norgate-AV/genlinx/compare/v1.1.1..v1.1.2) - 2024-01-04

### 🐛 Bug Fixes

- *(build)* Add empty arrays to localconfig - ([f54d80b](https://github.com/Norgate-AV/genlinx/commit/f54d80b846156397e9113b44dce998d992cb7b73))

## [1.1.1](https://github.com/Norgate-AV/genlinx/compare/v1.1.0..v1.1.1) - 2024-01-04

### 🐛 Bug Fixes

- *(build)* Quick and dirty fix for local config issue - ([d7cbf0f](https://github.com/Norgate-AV/genlinx/commit/d7cbf0fdb1f04db536fdcd167eff0384b4a25782))

## [1.1.0](https://github.com/Norgate-AV/genlinx/compare/v1.0.0..v1.1.0) - 2023-09-20

### 🚀 Features

- Add man page file - ([957da87](https://github.com/Norgate-AV/genlinx/commit/957da87737d71f487c54d0e93e791c47d385e410))

## [1.0.0](https://github.com/Norgate-AV/genlinx/tree/v1.0.0) - 2023-09-18

### 🎉🎉 Initial Release

First release of genlinx.
