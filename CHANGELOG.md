# Changelog

## [2.7.0](https://github.com/Norgate-AV/genlinx/compare/v2.6.1...v2.7.0) (2025-04-22)

### 🌟 Features

- **build:** add argument for source files so no flag is required ([e2b2c5d](https://github.com/Norgate-AV/genlinx/commit/e2b2c5d8ad9cf8f2c56d34a80df32cc4a4acbf1e))

### 📖 Documentation

- add RunNetLinxCompiler doc ([e4337e8](https://github.com/Norgate-AV/genlinx/commit/e4337e86b9576c24641d377c3d9ee10cc7cf40f8))
- **readme:** update readme ([4ecc6c8](https://github.com/Norgate-AV/genlinx/commit/4ecc6c86d317a621f50a4e50d941d33f3c03bdc8))

### ✨ Refactor

- **build:** make cfgFiles optional ([bc2434a](https://github.com/Norgate-AV/genlinx/commit/bc2434a53808b3ac70606e9d3b2a095b56c7cbbb))

## [2.6.1](https://github.com/Norgate-AV/genlinx/compare/v2.6.0...v2.6.1) (2024-11-22)

### 🐛 Bug Fixes

- **find:** correctly display json output in pretty format ([96a8937](https://github.com/Norgate-AV/genlinx/commit/96a89376a23b4a27afc023bb8fa98c0868505b36))

## [2.6.0](https://github.com/Norgate-AV/genlinx/compare/v2.5.0...v2.6.0) (2024-11-22)

### 🌟 Features

- **find:** add support for finding muse controllers ([72f7d55](https://github.com/Norgate-AV/genlinx/commit/72f7d559396519511194e058f4f3c2ad1cd07de1))

### 🐛 Bug Fixes

- **config:** handle exception if no global config is found ([3f232ce](https://github.com/Norgate-AV/genlinx/commit/3f232ce9c36c9a6ebeb66a165afdd25132933687))

## [2.5.0](https://github.com/Norgate-AV/genlinx/compare/v2.4.0...v2.5.0) (2024-08-14)

### 🌟 Features

- **find:** add cli option for timeout with 6000ms as default ([772988e](https://github.com/Norgate-AV/genlinx/commit/772988e4ee77f3f8f655012c4f8ca871a95247e3))

## [2.4.0](https://github.com/Norgate-AV/genlinx/compare/v2.3.0...v2.4.0) (2024-08-02)

### 🌟 Features

- integrate findlinx with the find command ([61770f2](https://github.com/Norgate-AV/genlinx/commit/61770f2a0e43b51ac34ceb53403f382e120040e2))

### 📖 Documentation

- update contents in readme to link to default config ([faa9362](https://github.com/Norgate-AV/genlinx/commit/faa936255aabeb1a3c589121764ebec56820656c))
- update readme and man page with find command ([dfebf5b](https://github.com/Norgate-AV/genlinx/commit/dfebf5b6dee0f56541b7915d786b3beed0bf22e3))

## [2.3.0](https://github.com/Norgate-AV/genlinx/compare/v2.2.3...v2.3.0) (2024-08-02)

### 🌟 Features

- **config:** add --add and --remove options for adding/removing from an array ([a0ee5e6](https://github.com/Norgate-AV/genlinx/commit/a0ee5e640437da10f59645dac7162981bceb31f6))
- **config:** add --delete option for deleting a key ([0c35ed1](https://github.com/Norgate-AV/genlinx/commit/0c35ed176b8299e4c4930353a1a1b297a57473d4))
- **config:** add --edit option to open config in default editor ([4bc70ac](https://github.com/Norgate-AV/genlinx/commit/4bc70ac19c4c5e1229cd5b86395ca1a735c06320))
- **config-cmd:** add basic implementation of cmd functionality ([1e9ac7a](https://github.com/Norgate-AV/genlinx/commit/1e9ac7a192e77873f13d1516f5ada46e66d71603))
- add config command to main cli ([4ab0192](https://github.com/Norgate-AV/genlinx/commit/4ab0192caffd20aadb1127dfab098cd53611069a))
- **config:** add config schema for validation ([f1c299e](https://github.com/Norgate-AV/genlinx/commit/f1c299e88a46c2a244ab2e9cdbae67c4c71b418d))
- **config:** add core section ([1ea0109](https://github.com/Norgate-AV/genlinx/commit/1ea01096f368a8a3f55f08f57bb0fd958a82cd57))
- add getLocalAppConfigFilePath function ([96c03e0](https://github.com/Norgate-AV/genlinx/commit/96c03e0e92b6b23eff1a0baffb39bb93c3c6bf11))
- **archive-cmd:** add negated boolean options ([56c414e](https://github.com/Norgate-AV/genlinx/commit/56c414e78c751b9c58e8c1aaa11ca8f241129844))
- add resolvePaths function ([e785cec](https://github.com/Norgate-AV/genlinx/commit/e785cec75760b49b946608fec69eebe33ed579be))
- **config:** add schema for config ([beb53f3](https://github.com/Norgate-AV/genlinx/commit/beb53f3f12d44b1bf50c1c06c54725f125e84895))
- add selectWorkspaceFiles function ([688053a](https://github.com/Norgate-AV/genlinx/commit/688053ac92237e71914437d563be58485e0a8e9b))
- export new util functions ([011aff7](https://github.com/Norgate-AV/genlinx/commit/011aff7d84074093f6c4d551c78e5e7e1693feb0))
- **config:** implement config get value ([f03f136](https://github.com/Norgate-AV/genlinx/commit/f03f13695d7647af69d8d32700ea54655cdddcaf))
- implement config listing and editing in editor ([cc144c3](https://github.com/Norgate-AV/genlinx/commit/cc144c3230fc957f348143f7731b570d5347070c))
- run paths through normalizer so / or \\ can be used ([1071b26](https://github.com/Norgate-AV/genlinx/commit/1071b260ebab2a0f280c668e9320e8ee4e5a6379))

### 🐛 Bug Fixes

- **config:** return after getting property ([7d24ee9](https://github.com/Norgate-AV/genlinx/commit/7d24ee944e090bc92d66c3417b27f98ad99d088e))

### 📖 Documentation

- update docs ([427619e](https://github.com/Norgate-AV/genlinx/commit/427619e92f10da2243994bed242ae67c951c01f1))
- update readme and man page ([c529be8](https://github.com/Norgate-AV/genlinx/commit/c529be8e35e68174bc4d42465339d5ff3e1d2c3c))

### ✨ Refactor

- **config:** add line break. remove console.log ([16b75de](https://github.com/Norgate-AV/genlinx/commit/16b75dea75152e2c4157048f22d45248e30ba37f))
- log to console after evaluating options ([f9431ad](https://github.com/Norgate-AV/genlinx/commit/f9431ad3fa9c40e8ed6b6835ae2425d6b131cb64))
- **config:** make list option conflict with edit ([1fa14de](https://github.com/Norgate-AV/genlinx/commit/1fa14de044a71b88de3170b5985b3410eeda9baa))
- **config:** make value arg an array ([a34ebdb](https://github.com/Norgate-AV/genlinx/commit/a34ebdbf48073dbae8da086e3ce3fea409610fbd))
- **getOptions:** merge default global config ([b464d68](https://github.com/Norgate-AV/genlinx/commit/b464d68df292b8e5d512c36bbb85b1206dc1038b))
- **config:** move file patsh to core section of config to be shared ([a398887](https://github.com/Norgate-AV/genlinx/commit/a3988874b66d3fae54228176684173f6242f912c))
- **config:** only take first value when setting property ([4a0a0f9](https://github.com/Norgate-AV/genlinx/commit/4a0a0f91c7c920388cd3b238c08b72f1306df6bd))
- **config:** outout lists in green ([5be2034](https://github.com/Norgate-AV/genlinx/commit/5be20347acd4f2b4cc16e7e53e39f182786e4970))
- **options:** pass config stores into mergician ([097f951](https://github.com/Norgate-AV/genlinx/commit/097f951ec387d98669b071ccbfb85147e6b184ef))
- prefix native imports with node: ([0af7f2e](https://github.com/Norgate-AV/genlinx/commit/0af7f2e49b764872e9bf50e82000f77fd51930b1))
- **config:** process options within the action ([483cf4c](https://github.com/Norgate-AV/genlinx/commit/483cf4cff727f4dc32d2562f4ab3ba9f905c1841))
- **global-config:** refactor to use conf package ([0b205f1](https://github.com/Norgate-AV/genlinx/commit/0b205f151f33c854cba3a1bb26e4cb8af0234c28))
- **local-config:** refactor to use conf package ([f319815](https://github.com/Norgate-AV/genlinx/commit/f31981548bd9a0e5ce1373a633f33aa2ba2f7875))
- **config:** remove delete key option ([a0fbc7c](https://github.com/Norgate-AV/genlinx/commit/a0fbc7cbe763c7b404cf97de7b528e93d13261f0))
- remove fs-extra dependency ([1067aed](https://github.com/Norgate-AV/genlinx/commit/1067aed4ec57c3d049f7cffaf4c1f31e07e8ff1d))
- replace string literals with name variable ([4e2d661](https://github.com/Norgate-AV/genlinx/commit/4e2d661bc9fc5354bcbc52ad2d44786c290278ce))
- **global-config:** return complete object rather than just the store ([4548838](https://github.com/Norgate-AV/genlinx/commit/45488381554ceba5330a2b5fcf3a9c7147007f13))
- **local-config:** return complete object rather than just the store ([899a052](https://github.com/Norgate-AV/genlinx/commit/899a05227395e6b35d94597018f756c9d7636876))
- return unique sets ([0e3182c](https://github.com/Norgate-AV/genlinx/commit/0e3182c557025ec9544cc76814cae51499955c9c))
- split app/global/local config function into separate files ([320d918](https://github.com/Norgate-AV/genlinx/commit/320d918366135ac260eff474d87435f942435992))
- split into smaller functions and resolve relative paths ([f2efc1c](https://github.com/Norgate-AV/genlinx/commit/f2efc1c9e9544519dfe4d48a90a57a50051d71a4))
- **config:** update cli logic ([a705024](https://github.com/Norgate-AV/genlinx/commit/a705024e7651859fc3896384f92bc41d0d2f0d8f))
- **build:** update command description ([b6e349b](https://github.com/Norgate-AV/genlinx/commit/b6e349b5466a1941e98682c6041ff04b744b44d4))
- **config:** update option descriptions ([b41944c](https://github.com/Norgate-AV/genlinx/commit/b41944cba0dbec0acb45dc66b44cbbf219ad38b7))
- use getModuleName ([a0e6939](https://github.com/Norgate-AV/genlinx/commit/a0e6939e3bb467c131a84a9ce5c9fc286e382f7d))
- use getPackageJson function when getting configs ([cbc59c6](https://github.com/Norgate-AV/genlinx/commit/cbc59c618972542ac14d4c41134582a53d9f03af))
- **config:** use new config member ([38be685](https://github.com/Norgate-AV/genlinx/commit/38be6851e7dd0027ad1d5fa5c9e3a7cd35f0e15a))
- use new function to get file path. return file path in object ([cf2149e](https://github.com/Norgate-AV/genlinx/commit/cf2149e9c9fbb047afedf6591859b7b298055ec2))
- **config:** use utils to print complex objects ([6af6c32](https://github.com/Norgate-AV/genlinx/commit/6af6c32ea1a4e355eb7466fc3cf082b918747819))

## [2.2.3](https://github.com/Norgate-AV/genlinx/compare/v2.2.2...v2.2.3) (2024-07-27)

### 🐛 Bug Fixes

- prepend arrays in config and dont sort ([86c8887](https://github.com/Norgate-AV/genlinx/commit/86c88872156b5a003562248a220ac897236e1476))

## [2.2.2](https://github.com/Norgate-AV/genlinx/compare/v2.2.1...v2.2.2) (2024-07-24)

### ⚠ BREAKING CHANGES

- **config:** remove RMS locations from the default config

### 🐛 Bug Fixes

- use @inquirer/prompts for selecting files ([4ad90cd](https://github.com/Norgate-AV/genlinx/commit/4ad90cd77da30a73e0b6dfe3dc4df5060e7f4e5d))

### ✨ Refactor

- **config:** remove RMS locations from the default config ([7eb65e6](https://github.com/Norgate-AV/genlinx/commit/7eb65e6eeab8f92285b726be81dbd31e26cc4680))

## [2.2.1](https://github.com/Norgate-AV/genlinx/compare/v2.2.0...v2.2.1) (2024-02-01)

### 🐛 Bug Fixes

- **archive:** correctly determine file type or return "Other" ([66e68df](https://github.com/Norgate-AV/genlinx/commit/66e68dfcb7202fbf5ba895e068acf0b09f4a8eee))

### 📖 Documentation

- fix broken links ([0976139](https://github.com/Norgate-AV/genlinx/commit/097613919a08d617bcccc8e4a6a56bc9efbe7b35))

## [2.2.0](https://github.com/Norgate-AV/genlinx/compare/v2.1.1...v2.2.0) (2024-02-01)

### 🌟 Features

- **build:** allow for building multiple source files ([f2f8243](https://github.com/Norgate-AV/genlinx/commit/f2f8243928ecdd5733ab0ea0b6924467c98693bc))

### 📖 Documentation

- add verbose flag ([aaaeb4e](https://github.com/Norgate-AV/genlinx/commit/aaaeb4e4024f4a1082a87bdd5354154d4556c93b))
- update build docs for multiple source file build ([07c7fec](https://github.com/Norgate-AV/genlinx/commit/07c7fec6a38876245a62e497d8e1f9f077cb2f54))
- update installation commands to include PNPM ([065ced3](https://github.com/Norgate-AV/genlinx/commit/065ced3c82b7d04536ab7767f9d29a757fbfe430))
- update to show latest config options ([6bc2815](https://github.com/Norgate-AV/genlinx/commit/6bc2815ac8d88491d57702b44ee3e60e36a7301c))

## [2.1.1](https://github.com/Norgate-AV/genlinx/compare/v2.1.0...v2.1.1) (2024-02-01)

### 🐛 Bug Fixes

- only log using local config if verbose ([f72744d](https://github.com/Norgate-AV/genlinx/commit/f72744de9cd7a7eca207caf37cc57cd63408965c))

## [2.1.0](https://github.com/Norgate-AV/genlinx/compare/v2.0.3...v2.1.0) (2024-02-01)

### 🌟 Features

- add --verbose flag to commands ([e822091](https://github.com/Norgate-AV/genlinx/commit/e82209108b50aa1d64aa4807a2f349ebc5d64df6))

## [2.0.3](https://github.com/Norgate-AV/genlinx/compare/v2.0.2...v2.0.3) (2024-02-01)

### 🐛 Bug Fixes

- add year to banner dynamically ([0425c24](https://github.com/Norgate-AV/genlinx/commit/0425c242d0e35336b50730ff44e9d97adeebb22d))
- open correct file to obtain latest version ([34617d2](https://github.com/Norgate-AV/genlinx/commit/34617d2167274b3ec84e2bee27b68e531be329d1))

## [2.0.2](https://github.com/Norgate-AV/genlinx/compare/v2.0.1...v2.0.2) (2024-01-31)

### 🐛 Bug Fixes

- get most upto date app version ([9c7c784](https://github.com/Norgate-AV/genlinx/commit/9c7c784555c799db47987aee1621f90faad29b1d))

## [2.0.1](https://github.com/Norgate-AV/genlinx/compare/v2.0.0...v2.0.1) (2024-01-31)

### 🐛 Bug Fixes

- replace \_\_dirname in code ([9432445](https://github.com/Norgate-AV/genlinx/commit/9432445ad4d26ab3e19212f444a9fdf4388e0663))

## [2.0.0](https://github.com/Norgate-AV/genlinx/compare/v1.1.4...v2.0.0) (2024-01-31)

### ⚠ BREAKING CHANGES

- convert codebase to typescript

### 🌟 Features

- convert codebase to typescript ([937de15](https://github.com/Norgate-AV/genlinx/commit/937de15f1a30b354d7f104785218865c2f7f58bf))

## [1.1.4](https://github.com/Norgate-AV/genlinx/compare/v1.1.3...v1.1.4) (2024-01-04)

### Bug Fixes

- **build:** remove accidental console.log ([a71b27e](https://github.com/Norgate-AV/genlinx/commit/a71b27e628f95d10d2fb1d8501a25fa71b3e313c))

## [1.1.3](https://github.com/Norgate-AV/genlinx/compare/v1.1.2...v1.1.3) (2024-01-04)

### Bug Fixes

- **build:** filter out falsy props from options ([e105abe](https://github.com/Norgate-AV/genlinx/commit/e105abe536ffe37c74ea299124835113ec11d067))

## [1.1.2](https://github.com/Norgate-AV/genlinx/compare/v1.1.1...v1.1.2) (2024-01-04)

### Bug Fixes

- **build:** add empty arrays to localconfig ([f54d80b](https://github.com/Norgate-AV/genlinx/commit/f54d80b846156397e9113b44dce998d992cb7b73))

## [1.1.1](https://github.com/Norgate-AV/genlinx/compare/v1.1.0...v1.1.1) (2024-01-04)

### Bug Fixes

- **build:** quick and dirty fix for local config issue ([d7cbf0f](https://github.com/Norgate-AV/genlinx/commit/d7cbf0fdb1f04db536fdcd167eff0384b4a25782))

## [1.1.0](https://github.com/Norgate-AV/genlinx/compare/v1.0.0...v1.1.0) (2023-09-20)

### Features

- add man page file ([957da87](https://github.com/Norgate-AV/genlinx/commit/957da87737d71f487c54d0e93e791c47d385e410))

## 1.0.0 (2023-09-18)

### Features

- accept array of workspace files. search if none provided. create cfg for all ([f4a4f58](https://github.com/Norgate-AV/genlinx/commit/f4a4f582b78507d7ef777f3a1862469d660480fe))
- add ability to ignore certain directories ([6862c77](https://github.com/Norgate-AV/genlinx/commit/6862c77df901ec0e67ffd38f922176a78d21dd56))
- add all option to bypass interactive prompts ([b4ca057](https://github.com/Norgate-AV/genlinx/commit/b4ca0579a9179e2f60a13177144b3f1658bdf3be))
- add APW class ([c5335d4](https://github.com/Norgate-AV/genlinx/commit/c5335d4597f18a164c5b764adca3186374987880))
- add apw directory to extra file locations for archive command options ([07716bf](https://github.com/Norgate-AV/genlinx/commit/07716bf5b3d70da576de8f073f0702472f867607))
- add apw discovered module paths to extra module paths ([bdef229](https://github.com/Norgate-AV/genlinx/commit/bdef2290287310e3e3c61761796236afd8328b52))
- add archive builder class ([84eee88](https://github.com/Norgate-AV/genlinx/commit/84eee8859261ca0f37ebd756582a947d75228222))
- add archive cli option for extra file archive location ([0338f3e](https://github.com/Norgate-AV/genlinx/commit/0338f3e6962e6d902078c3a23beda47ef855c5e4))
- add archive command boilerplate ([14ce8b9](https://github.com/Norgate-AV/genlinx/commit/14ce8b9e3d788b4be07e6a8c5daeba634495cd02))
- add archive config option to search for files not in workspace ([65c688a](https://github.com/Norgate-AV/genlinx/commit/65c688a0ddcac123a16ab90107d46aa425bfb0ec))
- add archive config options for including compiled files ([4291f1c](https://github.com/Norgate-AV/genlinx/commit/4291f1c8fea5bf9d3a8682c0130ce8fe8c3d278a))
- add archive extra file locations to config ([fe46df7](https://github.com/Norgate-AV/genlinx/commit/fe46df795d891da4794a260cc919c97b834ea24c))
- add archive extra files archive location to config ([dd9db81](https://github.com/Norgate-AV/genlinx/commit/dd9db81be6fa144fdc75d032bfae0afaa17a7319))
- add archive items file type reference object ([9daa63d](https://github.com/Norgate-AV/genlinx/commit/9daa63d512f6b42f25b7a9ee32bf4bedea5f5366))
- add archive options getter to options class ([0f1a8c6](https://github.com/Norgate-AV/genlinx/commit/0f1a8c6e70aaf61220eb5fcbb222cf13f82e3970))
- add ArchiveItemFactory ([98c5c89](https://github.com/Norgate-AV/genlinx/commit/98c5c896b5f29719fb48634e51a4294fd3736ab9))
- add banner to cli help output ([df46dfc](https://github.com/Norgate-AV/genlinx/commit/df46dfc47bd77b848072f30ed23c8e1b492f8194))
- add boilerplate for config command ([b7a4295](https://github.com/Norgate-AV/genlinx/commit/b7a42955fb1bf0c0798b8957af70fd586db7f527))
- add build command boilerplate ([ebe0048](https://github.com/Norgate-AV/genlinx/commit/ebe0048a4a162fdc945946394190f22278d6cdaa))
- add CfgBuilder class ([9a9143f](https://github.com/Norgate-AV/genlinx/commit/9a9143f9d8c461f18a35121fb122533f879fa491))
- add compiled file extension map ([282ed29](https://github.com/Norgate-AV/genlinx/commit/282ed292c782a3ef034447c3fb0cfbe8afe583ff))
- add create cfg option to global config ([41d8240](https://github.com/Norgate-AV/genlinx/commit/41d82401156e2665d53d62d1e73b37da58243793))
- add default archive output file to config ([d60ac1b](https://github.com/Norgate-AV/genlinx/commit/d60ac1bbe276d98c241023edd523e168d67e7456))
- add default config ([9970df2](https://github.com/Norgate-AV/genlinx/commit/9970df272b73efe28e9ce3094c2d0a9e57f495eb))
- **build-cmd:** add default include, module, library paths to global config for building axs files ([12e1cf6](https://github.com/Norgate-AV/genlinx/commit/12e1cf6ad91d6592c0399aefe87280fcf289f317))
- add default output file to config ([64a5a2e](https://github.com/Norgate-AV/genlinx/commit/64a5a2e1ab3251882f8b78254f7f8d446c8884d1))
- add Dockerfile ([2b0e103](https://github.com/Norgate-AV/genlinx/commit/2b0e10349d485f368713b3ca07f9dc4f6afda6d2))
- add env file to archive with extra files ([19d2580](https://github.com/Norgate-AV/genlinx/commit/19d2580e28c18dde3a4a57c8a32e3649e94e79cc))
- add EnvItem archive item class ([7eb5013](https://github.com/Norgate-AV/genlinx/commit/7eb50134bd4199f8efb38fb5709cb68be70775ce))
- add extra file archive location option ([0c377f5](https://github.com/Norgate-AV/genlinx/commit/0c377f5c1644810bcdd5393baca05af61ca70203))
- add extra file extensions to file type map ([b20c84a](https://github.com/Norgate-AV/genlinx/commit/b20c84a517d82dfc94d0a5201f1b481f95a91e31))
- add extra file type extensions ([75eacc8](https://github.com/Norgate-AV/genlinx/commit/75eacc8f6b707f9eb018781d38b8e30de814d559))
- **build-cmd:** add extra options to allow building source files ([d74e4d5](https://github.com/Norgate-AV/genlinx/commit/d74e4d5aab769fb03afd72596b29a7a0884c7a76))
- add extra options to archive cli command ([6a7489c](https://github.com/Norgate-AV/genlinx/commit/6a7489c673e01ad48e2e6bb9c7662f121b1d601c))
- add file extension map ([0ea9883](https://github.com/Norgate-AV/genlinx/commit/0ea988327cd0a4c5c414b6fd0bc4ca79c4dbe24f))
- add footer message to all help text directing to man page for extra help ([22c35a5](https://github.com/Norgate-AV/genlinx/commit/22c35a5d488d9c971f4b6b1a54aaaad173e58c6f))
- add found files to the archive in a .genlinx directory ([2b2c82e](https://github.com/Norgate-AV/genlinx/commit/2b2c82ee381753ba4b34c8ffbd3764a2db21a6c8))
- add function calls for build and archive ([6d8d234](https://github.com/Norgate-AV/genlinx/commit/6d8d234df1b372c8dcac308f646f0d455948855d))
- add GeneralItem class ([177a618](https://github.com/Norgate-AV/genlinx/commit/177a6184ec73267dd371f2ad6c62b332c73b6b86))
- add getFileDirectories method ([bc6b41e](https://github.com/Norgate-AV/genlinx/commit/bc6b41e2081e6882b175f47e6f03bf5e523ab004))
- **nlrc:** add getSourceBuildCommand function ([f386a53](https://github.com/Norgate-AV/genlinx/commit/f386a53f29c8f4a56a942144b6a1e51cb83a33be))
- add getter property to generate a list of file referenced that are not in the workspace ([0064e86](https://github.com/Norgate-AV/genlinx/commit/0064e86d0be2719e8da381d05b1ebfb4a026ba7f))
- add getter to return all files ([a53a5bc](https://github.com/Norgate-AV/genlinx/commit/a53a5bc0d2a1eb75245b5f1ff50205ede3c8efec))
- add ignored files list for archive command to global config ([1e1601d](https://github.com/Norgate-AV/genlinx/commit/1e1601d8dc630af2cc235cd51a5eca1d192faba0))
- add initial implementation of genlinx cfg command ([73550ea](https://github.com/Norgate-AV/genlinx/commit/73550ea253f8159d8b2f223d2b5218fdcada9be7))
- **utils:** add loadAPW helper function ([6d8bbd2](https://github.com/Norgate-AV/genlinx/commit/6d8bbd22ebc11fe29fb9ea3ba7e3943cfa730773))
- add local config options to each command ([a0a6fce](https://github.com/Norgate-AV/genlinx/commit/a0a6fcecd5893fbd6d96de7e69b454f7218eaeb4))
- add map to identify extra file types ([0418021](https://github.com/Norgate-AV/genlinx/commit/0418021bfcdfffe117e24f2fb67076dc56c5aa6f))
- add method for searching through extsing files based on regex pattern ([6f9dafe](https://github.com/Norgate-AV/genlinx/commit/6f9dafe9adcd4f4eb17960c04bd9744839365994))
- add method to check if a file already exists in the workspace ([98d15b6](https://github.com/Norgate-AV/genlinx/commit/98d15b6dfd5c67cb43e8e8e201fde9ef72d3367f))
- add method to check if file is a readable file ([5ca45a1](https://github.com/Norgate-AV/genlinx/commit/5ca45a12c7f460063bc9a234d1b8d7b61ab0d9a8))
- add method to get workspace id and add id getter ([1935a66](https://github.com/Norgate-AV/genlinx/commit/1935a66a384285936cfb0faa4e58f3d457e0f0e3))
- add method to setup options for builder from cli option and app config ([1f407ba](https://github.com/Norgate-AV/genlinx/commit/1f407bafa38a093e6aeddcf855ec04de4829dbab))
- add ModuleItem class ([f0efdfd](https://github.com/Norgate-AV/genlinx/commit/f0efdfdddf2b5ac971e9af3129abb9f4dd1bd008))
- add nlrc args array to config ([1beda73](https://github.com/Norgate-AV/genlinx/commit/1beda73c41c79527965c19978bde4630a21ec8ce))
- add NLRC class ([b669bf3](https://github.com/Norgate-AV/genlinx/commit/b669bf34bb6bffdce0bf0f37f816c654a93ea698))
- add placeholder method to get apw file type ([7949438](https://github.com/Norgate-AV/genlinx/commit/79494384e2b43a1e20a4fcb46da83049c32af9af))
- add placeholder to get the apw file type ([12e78b2](https://github.com/Norgate-AV/genlinx/commit/12e78b21eea34b3ff11e7810ed6c6b0cbc1b6b80))
- add property to get masterSrcPaths ([3469d05](https://github.com/Norgate-AV/genlinx/commit/3469d0566575aee7e295c47960feafd965d44464))
- add property to get moduleFilePaths ([7198a3d](https://github.com/Norgate-AV/genlinx/commit/7198a3d2335c9f2018b01ea348be565ce48f43f4))
- add property to get the apw file path ([5251752](https://github.com/Norgate-AV/genlinx/commit/525175268be2242581ddef8fb45e36f609995c38))
- add property to get unique include file directories ([34bcdb5](https://github.com/Norgate-AV/genlinx/commit/34bcdb5262ed9bc44db59b6fa60549bb9ae4525e))
- add script files for symlinking extra files ([6af9745](https://github.com/Norgate-AV/genlinx/commit/6af974553684f24e6aede60ff3129b9e3445ce19))
- **lib:** add selectWorkspaceFiles function ([6ddd29c](https://github.com/Norgate-AV/genlinx/commit/6ddd29ca69f6f7279770c5d75f1244533c2b56c2))
- add shell execuable command arg to default config ([91bb0a1](https://github.com/Norgate-AV/genlinx/commit/91bb0a1f7a1a3cdf947074a09a6a518fccd1dc6c))
- add shell executable to default config ([ebbbef1](https://github.com/Norgate-AV/genlinx/commit/ebbbef1067b1bb42cdc25bfceaf7605c43adb24e))
- add SourceItem class ([917b178](https://github.com/Norgate-AV/genlinx/commit/917b178d85ec47433a7a383e9bae7fbd792c22b8))
- add tpd file extension ([1ccd604](https://github.com/Norgate-AV/genlinx/commit/1ccd604d70c9efdfec0e3a282e7a61a635ab8613))
- add utility function to get app config file path ([d39d208](https://github.com/Norgate-AV/genlinx/commit/d39d208834241ae722ef3988a129f11d96558d01))
- add utility function to recursive walk directories for files ([0401cff](https://github.com/Norgate-AV/genlinx/commit/0401cffa953ae5b158f7abb70823fcbfc85c30ad))
- add utility functions for getting local config ([9168d5a](https://github.com/Norgate-AV/genlinx/commit/9168d5a4dde9efde9d175e89563cb0eb8cbf26bd))
- add utility helpers ([5f427ff](https://github.com/Norgate-AV/genlinx/commit/5f427ff7f9641b19794b543b81e2ae74abca259d))
- add WorkspaceItem class ([9c81b09](https://github.com/Norgate-AV/genlinx/commit/9c81b097988c1420b0d6d4686dda7f676c997b49))
- append additional include paths gleaned from apw ([a1b5ada](https://github.com/Norgate-AV/genlinx/commit/a1b5ada21de971dc5dfa0f75de99e13a4d067e90))
- **build-cmd:** build source file if passed via cli ([2a40133](https://github.com/Norgate-AV/genlinx/commit/2a40133ce125a3fcaf0f3653dfaa510f9349ba52))
- check for ignored files when gathering extra files ([e1c1034](https://github.com/Norgate-AV/genlinx/commit/e1c1034a623f9a29a4a7e05cc938fc5469200f8e))
- **build-cmd:** check os platform and exit if not being run on Windows ([9262a1d](https://github.com/Norgate-AV/genlinx/commit/9262a1d01d62f5626e9c084ae9c2a17cbcf7eb5e))
- define additional file types ([48538f4](https://github.com/Norgate-AV/genlinx/commit/48538f45b4d7099c3e34a798d249ec523cbaa4e2))
- export classes from lib ([a8e8b5f](https://github.com/Norgate-AV/genlinx/commit/a8e8b5fb9c2cef5d60b3bd30889664eecb76652e))
- get archive extra file locations option from options class ([b06d906](https://github.com/Norgate-AV/genlinx/commit/b06d906ac1e8cba7bb25e08bf0c55abc0f8801ee))
- get archive include files not in workspace option from options class ([cc5ee78](https://github.com/Norgate-AV/genlinx/commit/cc5ee78c88a1689b0f8ad29a15e080cf9993b39f))
- get executable from config and execute with execa ([afa09e8](https://github.com/Norgate-AV/genlinx/commit/afa09e8fce6109e732558afd255b541a8b0459f4))
- implement archive command logic ([f09d26c](https://github.com/Norgate-AV/genlinx/commit/f09d26c463743d059b71152ee8c76a5a331eca9b))
- include symlink scripts with extra files in archive ([4f2055c](https://github.com/Norgate-AV/genlinx/commit/4f2055cf5cd0c3c2da4378c8c269f27b094c0576))
- install default config if not found on disk ([d18ea95](https://github.com/Norgate-AV/genlinx/commit/d18ea956c816967bbc66a122a8b62f00a95246ff))
- print found workspace files ([abe7f27](https://github.com/Norgate-AV/genlinx/commit/abe7f271fc65157eae13c81c0326823542f3adb0))
- prompt user to select files if more than one found ([f5cd854](https://github.com/Norgate-AV/genlinx/commit/f5cd854445deeb0d0c051783f089780b5c4d3dc5))
- read stdout and throw exception if error found in build process ([9024270](https://github.com/Norgate-AV/genlinx/commit/90242703a2681bad042d1b178c1f204490f4cacb))
- return archive include compiled file options ([35dd2b4](https://github.com/Norgate-AV/genlinx/commit/35dd2b4e27df452a2a7a5c1d3e249dd59be7649f))
- return EnvItem for env file type ([2bfe7bc](https://github.com/Norgate-AV/genlinx/commit/2bfe7bcf64b7b54532c5c0f7409f7ddd5ee188b5))
- search each known file location to get list of files ([feef461](https://github.com/Norgate-AV/genlinx/commit/feef461a507a2908bbb8219b08bd14cd1c931af6))
- search for and archive all workspace file if none specified ([7915b87](https://github.com/Norgate-AV/genlinx/commit/7915b87b88d767e62c95b2d24dc1553019207531))
- take local config into accoutn when building options ([e25b3f0](https://github.com/Norgate-AV/genlinx/commit/e25b3f08f74571d001f4b754f815f33d74390d99))

### Bug Fixes

- add comma after \*.axi ([3c75c98](https://github.com/Norgate-AV/genlinx/commit/3c75c9836d6fbc88214ff37cb455e1a431d9c90c))
- amend -C to -CFG ([b2cf5f3](https://github.com/Norgate-AV/genlinx/commit/b2cf5f38cfbf1c90225654ce12aeb5b4e2572170))
- **build-cmd:** await for get local config before processing options ([88a9ed7](https://github.com/Norgate-AV/genlinx/commit/88a9ed7971a670bfe5bdb35d8654219d1629549e))
- catch exceptions when reading directories ([fac0238](https://github.com/Norgate-AV/genlinx/commit/fac023832d902bb5f302f296d1e502fd6ace121d))
- correct boolean expressions ([32b65a6](https://github.com/Norgate-AV/genlinx/commit/32b65a6b4a12d92f301fb0f3f8db239a91ae27a5))
- correctly resolve include directories ([616fd67](https://github.com/Norgate-AV/genlinx/commit/616fd6708bbea4cab026473c0dbdfb856a186664))
- **build-cmd:** fix copy/paste error referencing arhive configs ([6591ce4](https://github.com/Norgate-AV/genlinx/commit/6591ce4810630d1c3c07abd880720b3519d0e539))
- **build-cmd:** fix error regex pattern ([b2bc995](https://github.com/Norgate-AV/genlinx/commit/b2bc9952657cb24f7754dbdd1ba19d89c0e576b5))
- **build-cmd:** fix pluralisation for sourceFile ([28b1bac](https://github.com/Norgate-AV/genlinx/commit/28b1bac98f38fce5441724ac11617f0d15749776))
- fix regex pattern match for getting workspace id ([82b17ad](https://github.com/Norgate-AV/genlinx/commit/82b17ad648df1ab6c3eb5ef63cf00480f8689643))
- fix typo in duet file type ([0178176](https://github.com/Norgate-AV/genlinx/commit/0178176495244fa300be906ef5e6d1f0f60f6d11))
- fix typo outputLogFileSuffix ([9a60d67](https://github.com/Norgate-AV/genlinx/commit/9a60d6767dca96646fd2e4f4385d0faeeece355c))
- fix use of regex.test ([b015bba](https://github.com/Norgate-AV/genlinx/commit/b015bba20c703ac340f4a59f0f74beac67e03e5d))
- fix warning regex pattern ([3cc583f](https://github.com/Norgate-AV/genlinx/commit/3cc583fafb9a0be2d63f4082d6c15639b640bf50))
- **nlrc:** get extra paths from nlrc object ([6e2495e](https://github.com/Norgate-AV/genlinx/commit/6e2495ec3555d081d88f1b201a7d11bf422e8e06))
- include .axi files ([ededbd4](https://github.com/Norgate-AV/genlinx/commit/ededbd49caae5c4b108776ff2399eaec74d7ec25))
- **apw:** invert boolean expression when checking file existence ([4b04ead](https://github.com/Norgate-AV/genlinx/commit/4b04eadee929c79545dfe3bbc66baf6c02634ad1))
- join script to scripts path to get full file path ([21df134](https://github.com/Norgate-AV/genlinx/commit/21df134a88937d92ad75c947d796477a8a8aab33))
- only return file name of thr workspace when getting all files ([28c7d7f](https://github.com/Norgate-AV/genlinx/commit/28c7d7f697a09665006e29a72eba2850a10aa6ef))
- remove cfg property reference when returning extra file locations for archive ([5740417](https://github.com/Norgate-AV/genlinx/commit/5740417d82ebf6bd81c8add57d15d8345a1334ec))
- remove extra curly brace ([5cf9d03](https://github.com/Norgate-AV/genlinx/commit/5cf9d03d441f06ec75f841ae9cc3a7e54e11bd75))
- return full path to apw ([1314381](https://github.com/Norgate-AV/genlinx/commit/1314381a0acc72e341d250501d3548bf1f6e4413))
- set reject option for execa to false ([ee99896](https://github.com/Norgate-AV/genlinx/commit/ee99896c81335f7e9bc684da9d800f17fc215e38))
- **apw:** throw exception if the return array lengh is zero ([1341a10](https://github.com/Norgate-AV/genlinx/commit/1341a1095a5ac9283253cbe74e9516f496c6bff5))
- **apw:** update files reges to allow for optional cr and lf ([b604461](https://github.com/Norgate-AV/genlinx/commit/b604461be7060a8098bba66eae0277ea1ce20ccb))
- **apw:** update workspace id regex to allow for optional cr and lf ([361018e](https://github.com/Norgate-AV/genlinx/commit/361018e30802e7bae5741cca955625d9c79b4102))

### Documentation

- add basic boilerplate for man page ([6586321](https://github.com/Norgate-AV/genlinx/commit/658632186501107ebc29ec60a4027cd262768a90))
- add configuration file documentation and update toc ([9071e78](https://github.com/Norgate-AV/genlinx/commit/9071e786899cc6cf453768c1787ee0aa7d0c9cc6))
- add contributing.md ([5874905](https://github.com/Norgate-AV/genlinx/commit/5874905b0e7e75304b0af08b815f22c0ccf183bc))
- add netlinx studio logo ([fa7be71](https://github.com/Norgate-AV/genlinx/commit/fa7be713fc00312bdd18f0c3030706d3910c56e3))
- update archive command display options ([3d46f46](https://github.com/Norgate-AV/genlinx/commit/3d46f466c5bdcda99c9205dc707d4f0d8e2a578f))
- update archive command usage ([402681c](https://github.com/Norgate-AV/genlinx/commit/402681c7ac6444956c8380f7f54f749d8811c560))
- update cfg command usage ([2573b2b](https://github.com/Norgate-AV/genlinx/commit/2573b2b8136fffdf0aa85036ba58b1440f6cc6f7))
- update cfg command usage ([03794a9](https://github.com/Norgate-AV/genlinx/commit/03794a9eecbdab86d65c3506bf1db651b9c59607))
- update cli usage and sort command alphabetically ([e104f42](https://github.com/Norgate-AV/genlinx/commit/e104f420cbc23107d630d31cf19cfa0f4594b147))
- **readme:** update cli usage for each command ([44e7778](https://github.com/Norgate-AV/genlinx/commit/44e77789c05f6e03a8cb2dd19eb1d11b0c75cad4))
- update contributors ([8e2e06d](https://github.com/Norgate-AV/genlinx/commit/8e2e06d5c240585dfc9f21fcf6b809f99e8fcc97))
- update default configuration file ([bd4c0f4](https://github.com/Norgate-AV/genlinx/commit/bd4c0f43f624a038ef2749a4bdece0264ec8089d))
- update doc with config command usage ([3dd0cfb](https://github.com/Norgate-AV/genlinx/commit/3dd0cfba789ccc397311c428c779d7b8bdd14d75))
- update docs regaingd local configuration ([95bd60c](https://github.com/Norgate-AV/genlinx/commit/95bd60c0145c92304d1f6e11de966b29bb4465b6))
- update package name for installation ([85c41d0](https://github.com/Norgate-AV/genlinx/commit/85c41d017fe173451dcea684d1bdf2b4c73ec483))
- update readme ([e782017](https://github.com/Norgate-AV/genlinx/commit/e7820172525c5554baf1c1f50bf160afc9aaf5ad))
- update readme with command line displays ([fbb5d55](https://github.com/Norgate-AV/genlinx/commit/fbb5d556258a5bd93f78d52ba42b0c4192a15d56))

### Style

- apply prettier rules ([98b5ea3](https://github.com/Norgate-AV/genlinx/commit/98b5ea350d618f13869ae617e26128f79ccd5168))

### Refactor

- add "all" option to default global config ([2d74d89](https://github.com/Norgate-AV/genlinx/commit/2d74d894e7b81b5a661412b76936317d0c1fc0f2))
- add apw file directory to search locations ([75020af](https://github.com/Norgate-AV/genlinx/commit/75020affadf07bc2bb528cc5862b206404f8bc57))
- add command path and args to execa command ([1e51332](https://github.com/Norgate-AV/genlinx/commit/1e5133291c9774a31affd40526b14daae98c2bb7))
- add workspace file reference when printing type to console ([947fb17](https://github.com/Norgate-AV/genlinx/commit/947fb1756e2518baeae8a38fe7c23110751d7332))
- alphabetize util imports ([f89a21f](https://github.com/Norgate-AV/genlinx/commit/f89a21fd84c54f2f07670106cdf0f5875a5976aa))
- amend message when adding symlink scripts to archive ([58cd565](https://github.com/Norgate-AV/genlinx/commit/58cd5654feae6b3003db705df7a8aa0a0c7ad9f1))
- append apw include paths to existing option include paths on functional call ([9abeeb1](https://github.com/Norgate-AV/genlinx/commit/9abeeb14d20b4890563a01b361f7d8b668b72acf))
- await get local config ([5c4be57](https://github.com/Norgate-AV/genlinx/commit/5c4be572904d8f9cafe4a56958cac5e439ba5702))
- chain builder method for readability ([8c8573d](https://github.com/Norgate-AV/genlinx/commit/8c8573d3e648c46a89ce77c439bf825df977bd1b))
- change config nlrc args array to an options object ([eab942e](https://github.com/Norgate-AV/genlinx/commit/eab942ec007c8985ecec358ff4b0a94adc6442b2))
- change get extra files property to a method. refactor creation of array ([86cceab](https://github.com/Norgate-AV/genlinx/commit/86cceab0e5c9452f1055b598559663ad17762f1b))
- check args present before parsing. if not found, output help and exit gracefully ([f600b33](https://github.com/Norgate-AV/genlinx/commit/f600b33884be27033cc18d83e54d381d688d499a))
- **archive-builder:** convert all synchronous function call to async ([c50b88a](https://github.com/Norgate-AV/genlinx/commit/c50b88a1e9109753303835738b1d9df08804f6e2))
- **apw:** convert all synchronous function calls to async ([686e8a9](https://github.com/Norgate-AV/genlinx/commit/686e8a920c22fb327e20815134eebdf6916fa4be))
- **utils:** convert getGlobalAppConfig to be async ([d425f95](https://github.com/Norgate-AV/genlinx/commit/d425f95822aabb4b3601a4bd0c51b1af1fc3816d))
- **utils:** convert installDefaultGlobalAppConfig to async ([a523515](https://github.com/Norgate-AV/genlinx/commit/a5235153353bae2dc9014967bc6e214257fef62e))
- convert to static class ([2c652c2](https://github.com/Norgate-AV/genlinx/commit/2c652c2a04d62a0338b7aab379fa08529a271758))
- **utils:** convert walk directory to async ([f66c54b](https://github.com/Norgate-AV/genlinx/commit/f66c54b5bba1738bbad8312125311ebd6278cb4c))
- create file object prior to use ([2554531](https://github.com/Norgate-AV/genlinx/commit/2554531c1f6d2cc7af8abf9b7a10b4d48fede924))
- create method to add workspace files. add logic to start searching for extra files ([0c2d3ec](https://github.com/Norgate-AV/genlinx/commit/0c2d3ec23d20e73c254083d883d8024c08beef0f))
- create resolvesPaths function ([5ba349c](https://github.com/Norgate-AV/genlinx/commit/5ba349cc07e199fc8a88edce069a46f054cd6596))
- **cfg-builder:** default to "-R" for root directory ([127907d](https://github.com/Norgate-AV/genlinx/commit/127907dda5d1d65abde3798237aa236a0c1cb200))
- disable config option until fully implemented ([97cec26](https://github.com/Norgate-AV/genlinx/commit/97cec26deadb3ae7124201803082f639a7d8ddbe))
- do not pass global fields as args ([ef714ec](https://github.com/Norgate-AV/genlinx/commit/ef714eccb444d6ef95c16a1fe326be87b0cfacc1))
- explicitly check array length equal to zero ([11cbc0d](https://github.com/Norgate-AV/genlinx/commit/11cbc0df56ca92be0523d8da411dad6cfaea7b38))
- export commands from commands directory ([77e2680](https://github.com/Norgate-AV/genlinx/commit/77e26801d679daca088235ef26bfa1ff3ffcbfe4))
- extract file object creation to a seperate class ([e73cac1](https://github.com/Norgate-AV/genlinx/commit/e73cac1a2ac9884a9942621f70f0d882329ac223))
- extract getting refs from each file into a seperate method ([4827b00](https://github.com/Norgate-AV/genlinx/commit/4827b009b14661fb611cf8fbbe4e5357ca7eeefa))
- **cfgbuilder:** get axs root directory from apw.filePath dirname ([b58ad26](https://github.com/Norgate-AV/genlinx/commit/b58ad26b612686f95cff7a49300d4e020c263c2e))
- **local-config:** get local config by searching using find-up ([29f3aa3](https://github.com/Norgate-AV/genlinx/commit/29f3aa3b85552b2ac61133370202c37c4c3742df))
- get nlrc args from config ([16fa398](https://github.com/Norgate-AV/genlinx/commit/16fa398b6bb6be6ad98b6d2206bed20eb1cca30c))
- **build-cmd:** get number of error from errors array length ([a714288](https://github.com/Norgate-AV/genlinx/commit/a714288ddeb2fc5fa367be6a6264ffcf0a578b30))
- **build-cmd:** get number of warnings from warning array length ([3b2b563](https://github.com/Norgate-AV/genlinx/commit/3b2b563d4e4533e912212af6abbdb65904ea6b73))
- **nlrc:** get options for extra paths from config ([6773042](https://github.com/Norgate-AV/genlinx/commit/677304204af83e74d35522c771f4bed3925d71c2))
- get output file suffix from config ([5eb637c](https://github.com/Norgate-AV/genlinx/commit/5eb637ce33c638451f9ee0f2753c2a89d062ad1a))
- **apw:** hard code apw file path existence to true. this is already check when loading the apw ([7d0444f](https://github.com/Norgate-AV/genlinx/commit/7d0444f6d6bfbb06b8762f34c865bdb246088a84))
- if output file is not provided on cli append apw id to default ([ef7b618](https://github.com/Norgate-AV/genlinx/commit/ef7b61842a398c2ab8b5acc6139326eb26a84484))
- if output log file is not provided on cli, append apw id to default ([f9c0b78](https://github.com/Norgate-AV/genlinx/commit/f9c0b78826c59a7321617413955a7284a627d96f))
- import APW from parent directory ([6b45b28](https://github.com/Norgate-AV/genlinx/commit/6b45b281725fbf6ef88f4d84f6130b3a4bae2c97))
- make executable args an array ([86accc7](https://github.com/Norgate-AV/genlinx/commit/86accc771698e3697f2f0cfd233ee20859bd7aba))
- make function async ([842a83f](https://github.com/Norgate-AV/genlinx/commit/842a83f1cc799aad9ba8c532330c33e4a91eecf2))
- make function async and create file according workspace id ([d5f07f4](https://github.com/Norgate-AV/genlinx/commit/d5f07f42d0907d907d829b26f5dc244c9846f4e1))
- make into generic function returning list files matching specified extension ([0560496](https://github.com/Norgate-AV/genlinx/commit/05604964f6194bc77fd34db39957d292257dc39f))
- make more readable ([b9ecf39](https://github.com/Norgate-AV/genlinx/commit/b9ecf39d278a95d6af557b942429d29699782b11))
- make shorthand version flag lowercase ([aea2190](https://github.com/Norgate-AV/genlinx/commit/aea2190c3692ddd2a79f51741c4d4525aab95359))
- make workspace search more readable ([9368187](https://github.com/Norgate-AV/genlinx/commit/9368187e8f43e8374bf5e0bfcde5bff4af2b1d4a))
- move build error catching into a seperate function ([3d30b17](https://github.com/Norgate-AV/genlinx/commit/3d30b172cc202bfd6d0758daa7e8a46b22034112))
- move build options setup into options class ([dec3548](https://github.com/Norgate-AV/genlinx/commit/dec3548a4a4b557a4518a58ab9c8be908919a7e0))
- move cfg command into subdirectory ([d8a246d](https://github.com/Norgate-AV/genlinx/commit/d8a246d9af351c353a775c96f6f77095d6ebb913))
- move error parsing into seperate functions ([9abd32b](https://github.com/Norgate-AV/genlinx/commit/9abd32b6a4ec8765a8cdfd434bd6bffdb52a8f17))
- move getWorkspaceFiles function into a seperate file ([14e4700](https://github.com/Norgate-AV/genlinx/commit/14e470050a8aa304963497143178a5cd306b6ec5))
- **build-cmd:** only match the full error or warning ([00b5d1c](https://github.com/Norgate-AV/genlinx/commit/00b5d1c60dfd83ea3a8550afa6aea00519fa6835))
- only print file name to console ([12529f7](https://github.com/Norgate-AV/genlinx/commit/12529f72c89f10f0b3b06d6315b0431f7b17ec65))
- organise files into cli commands and cli actions ([e6c30b3](https://github.com/Norgate-AV/genlinx/commit/e6c30b3658b87153226baebe14aa379d1eca6971))
- prefix archive output file with apw id ([171ad1e](https://github.com/Norgate-AV/genlinx/commit/171ad1e39744fe28701f13ec547c2885170288a2))
- print warnings to console ([060b2df](https://github.com/Norgate-AV/genlinx/commit/060b2dfc6cca8997e535b8ba0a2511666e744a26))
- recursively look for extra file refs and find extra files ([e41a725](https://github.com/Norgate-AV/genlinx/commit/e41a725188e496e10a8dbc02425f45783ac1400c))
- refactor file sorting to be cleaner ([8ab2dfa](https://github.com/Norgate-AV/genlinx/commit/8ab2dfa21c6b48d0cd70a0fbe461f81f9947d9d2))
- **build:** refactor into separate functions ([4cedbfb](https://github.com/Norgate-AV/genlinx/commit/4cedbfbbf449e8285e49a5339fc94179d72d8d36))
- refactor into smaller chunks ([8ae5d91](https://github.com/Norgate-AV/genlinx/commit/8ae5d91c45247f1c06c680c434d14afc7d439cd2))
- refactor options getter into a simple utiltiy function ([a5fd7cf](https://github.com/Norgate-AV/genlinx/commit/a5fd7cffbce733dd85d5d46dc276aa9fe842cb67))
- refactor to use factory pattern to add extra files to archive ([4fd6b37](https://github.com/Norgate-AV/genlinx/commit/4fd6b37f0b6e828d3a4d3e6e07f2c9747d4ccd5b))
- reference file type object to get type ([5b28620](https://github.com/Norgate-AV/genlinx/commit/5b28620ef939f6305e165194ff2fe2d9b36b30c3))
- relocate options setup into a seperate class ([7d49fe3](https://github.com/Norgate-AV/genlinx/commit/7d49fe3b110842a0d6a154376d1bc01539b8958c))
- remove apw arg from action ([278bec1](https://github.com/Norgate-AV/genlinx/commit/278bec1d6ffa24199c99b6eb69a36428da8b6aa1))
- remove apw arg from action ([0abba30](https://github.com/Norgate-AV/genlinx/commit/0abba30110071a21efdce704fe90f1c37c4e7fc9))
- remove config build shell args ([4144fd4](https://github.com/Norgate-AV/genlinx/commit/4144fd4b61065d8b0ad805040a2f9a1cdc969a75))
- remove eslint rule disable comment ([5a9a2f8](https://github.com/Norgate-AV/genlinx/commit/5a9a2f8e6b68ad754b629ba9ed2e9a9b92b20965))
- remove file class in favour of simple js object ([4a63aec](https://github.com/Norgate-AV/genlinx/commit/4a63aec92877b7e41bf0b2b5ede513924c7cb739))
- remove find method ([c201eec](https://github.com/Norgate-AV/genlinx/commit/c201eecfa1206823a9148e6ef6593268d3589713))
- remove tenary operator for root directory option ([7fa635b](https://github.com/Norgate-AV/genlinx/commit/7fa635b34afe4184e76c5f864060bc8a1dd25a73))
- remove trailing backslashes from file paths ([6bf3e99](https://github.com/Norgate-AV/genlinx/commit/6bf3e992e42411d89aaf6a532d1ab1ace976c1cc))
- rename archive extra file locations key to be more specific ([db81e52](https://github.com/Norgate-AV/genlinx/commit/db81e52105495fef783b86b469206a01adb91ebf))
- rename build action to execute ([6d639ea](https://github.com/Norgate-AV/genlinx/commit/6d639ea0c4a937d60c9e172af8ee65dbb83dbed3))
- rename config build executable to shell ([137b3df](https://github.com/Norgate-AV/genlinx/commit/137b3df28a3e52b276825f2d791093b5634935f7))
- rename config to global config ([db7e031](https://github.com/Norgate-AV/genlinx/commit/db7e0317c800f8517511ae19fb41d6cbb4ae6ffc))
- rename extra file locations for archive options ([f27e6dd](https://github.com/Norgate-AV/genlinx/commit/f27e6dd78194c54cb24bf43e7c02be7f8ae15ffd))
- rename extraFiles to extraFileReferences. search for reference in files list ([f7d242b](https://github.com/Norgate-AV/genlinx/commit/f7d242b00fc26b8e710dada8364912139017480b))
- rename file type map and getter method ([010ae7c](https://github.com/Norgate-AV/genlinx/commit/010ae7c1dccdd81237a4be600cbad6c7657683e6))
- rename include path getter for consistency ([1661121](https://github.com/Norgate-AV/genlinx/commit/1661121ce2dfaa9c0dd98967d4baa335375a0a1f))
- replace addModuleFile and addMasterSrcFiles with addAxsFiles ([c9acc46](https://github.com/Norgate-AV/genlinx/commit/c9acc4699584dea2d760e42feccf43c5e619b339))
- replace build executable with build shell property ([5381d6e](https://github.com/Norgate-AV/genlinx/commit/5381d6ee8d5da0821b778dfa0313e2f70aa553a7))
- replace references to Norgate AV Solutions Ltd ([cdf6c63](https://github.com/Norgate-AV/genlinx/commit/cdf6c63f164db8bc57ccb3b1408b92bf50d357a2))
- resolve cfg root directory in CfgBuilder ([3e9c950](https://github.com/Norgate-AV/genlinx/commit/3e9c95003718121e0638ddcee59f99f314b65593))
- resolve relative file paths in local config ([44f6073](https://github.com/Norgate-AV/genlinx/commit/44f6073459837f2f7def5b79ac41236ef2056975))
- **cfg:** resolve relative paths from local config ([07ed70e](https://github.com/Norgate-AV/genlinx/commit/07ed70eb3c459e823e868139a180ac3b1f5b6049))
- return an error object with fullError and message props. print the fullError ([840d2f2](https://github.com/Norgate-AV/genlinx/commit/840d2f224d2785d335dd60ee74c6b791feda4061))
- return an object with command and args array ([1242574](https://github.com/Norgate-AV/genlinx/commit/1242574e8b50646c1d7488203550c0e19a020287))
- return include directories as a set ([dc3a279](https://github.com/Norgate-AV/genlinx/commit/dc3a2796e88ad1acf90813146b8c8debf1bb83e2))
- return result getFileDirectories ([f15dcce](https://github.com/Norgate-AV/genlinx/commit/f15dcce43ab70e386632375b68ba2321c76c7ac0))
- return unique arrays ([1b4c5dd](https://github.com/Norgate-AV/genlinx/commit/1b4c5dd7343fcd07944459af47572ff55f39d09a))
- run execa with shell option. set windows verbatim option. pipe to stdout ([03bb7d5](https://github.com/Norgate-AV/genlinx/commit/03bb7d5192be303aad5d3a482f1061d33e2ed70a))
- **nlrc:** shorten param name ([34e08f8](https://github.com/Norgate-AV/genlinx/commit/34e08f8658c9b8db024b31970b08f842c316f9b4))
- simplifiy replacement of spaces in apw id when returned ([021a4d7](https://github.com/Norgate-AV/genlinx/commit/021a4d72ed88c5a197c2a5de3999bb66e883c77b))
- **options:** simplify options building ([9bc34af](https://github.com/Norgate-AV/genlinx/commit/9bc34af4cff8c71be33658e36a81bfdcc749d651))
- split into smaller methods ([96bf762](https://github.com/Norgate-AV/genlinx/commit/96bf76250ae8f096027c49fb6c0cedae5caef23f))
- split public build method into smaller private methods ([e75ccc2](https://github.com/Norgate-AV/genlinx/commit/e75ccc231e45e87458ece471e502072006b455ca))
- spread args passed to execa ([655b8d9](https://github.com/Norgate-AV/genlinx/commit/655b8d902bb49203a21acaa793362084a49712a9))
- **cfg-cmd:** take root directory option from cli if present ([06bb808](https://github.com/Norgate-AV/genlinx/commit/06bb808a53edc44f06c4e2032a6863eb326f2d37))
- update all options to use Option class ([dd92dad](https://github.com/Norgate-AV/genlinx/commit/dd92dad69966903a4dcb85d02a63f333bf5eedc2))
- update archive cli option for extra file search locations ([6ec0d58](https://github.com/Norgate-AV/genlinx/commit/6ec0d583c1249da490972e90f6b57c2393064e1e))
- update archive command description ([c420090](https://github.com/Norgate-AV/genlinx/commit/c4200901fc51b2f162a247ef22edac57d1d26de5))
- update archive command wording to be more consistent ([6dbe890](https://github.com/Norgate-AV/genlinx/commit/6dbe89074f41fd87de78970275de33fa0a6bd669))
- update build command description ([e342447](https://github.com/Norgate-AV/genlinx/commit/e342447545999a29d40194f14fd7294c9b625dac))
- update build command wording to be more consistent ([9002314](https://github.com/Norgate-AV/genlinx/commit/9002314fc97ddb7d5bd9e63bdc2b0b63048355d8))
- update cfg command options ([8ceb54f](https://github.com/Norgate-AV/genlinx/commit/8ceb54f80f960f037a079b835f265bf34e3454ab))
- update cfg command options ([26056c2](https://github.com/Norgate-AV/genlinx/commit/26056c2dd9c440dac777025b4a02d784c0ef8343))
- update cfg command verbiage to be more consistent and explicit ([23c4c7c](https://github.com/Norgate-AV/genlinx/commit/23c4c7c1c3a68973428620f22a2564d0a38b6b7e))
- **build:** update command description ([576bb12](https://github.com/Norgate-AV/genlinx/commit/576bb12c0dcdc0e47a468da97101894f38798a5e))
- update option name to be more explicit ([bdfe932](https://github.com/Norgate-AV/genlinx/commit/bdfe932693f19269f878d094a5fd00b41f671417))
- update program description ([bab6dd2](https://github.com/Norgate-AV/genlinx/commit/bab6dd2a89dd06f8b89bb1cb16c20758712a4ea1))
- update variable name to be more readable ([f0515fc](https://github.com/Norgate-AV/genlinx/commit/f0515fc3c6382ddf1fb669415b182affd8d9210d))
- update wording for cfg command apw arg description ([6ff14ff](https://github.com/Norgate-AV/genlinx/commit/6ff14ff9327ddfec0a67e565dc00a463f54e3f52))
- use factory pattern to add workspace files to archive ([0314167](https://github.com/Norgate-AV/genlinx/commit/0314167465002aa4c1ef9b1c8ef27760c16b8a4f))
- use file type object to filter by file type ([fa64ffb](https://github.com/Norgate-AV/genlinx/commit/fa64ffb656d778b5cac7698e01cf053b515ad367))
- use getAppConfigFilePath utility function ([54427d1](https://github.com/Norgate-AV/genlinx/commit/54427d1ad84cb7b6b0a8cd31837bfae5cef5df06))
- use option property from config to get options flag for nlrc args ([0efa1cd](https://github.com/Norgate-AV/genlinx/commit/0efa1cd27b23826c835701b21f75db6b1aec1a0b))
- various refactors ([0fab4f2](https://github.com/Norgate-AV/genlinx/commit/0fab4f2c84e30eb6015e589566510f4272bfa874))
- **build-cmd:** various updates to make code a bit cleaner ([839cef7](https://github.com/Norgate-AV/genlinx/commit/839cef7eceb7c8c3060c5cd9e36d631e41ede962))
- **utils:** wrap logic in getFilesByExtension with a try/catch ([5833d7b](https://github.com/Norgate-AV/genlinx/commit/5833d7bb1e404343ee8d2888d65f3b77bb7b3bae))
