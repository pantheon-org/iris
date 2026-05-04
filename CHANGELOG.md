# Changelog

## [2026.5.11](https://github.com/pantheon-org/iris/compare/v2026.5.10...v2026.5.11) (2026-05-04)


### Features

* add Gemini, Qwen, Windsurf, and Copilot provider implementations ([#104](https://github.com/pantheon-org/iris/issues/104)) ([9688c9b](https://github.com/pantheon-org/iris/commit/9688c9bfec6ce4b7f4422c5a52afdf2650671311))
* **cli:** add --json flag to list, sync, and status commands ([#102](https://github.com/pantheon-org/iris/issues/102)) ([ff9e4dc](https://github.com/pantheon-org/iris/commit/ff9e4dc8281af173fc0f5f0167828240e9eb71d7))
* **opencode:** add doc-review plugin for idle session doc checks ([#90](https://github.com/pantheon-org/iris/issues/90)) ([711bba0](https://github.com/pantheon-org/iris/commit/711bba02ccd6b869b392a05b75122b95697be8f5))
* translate status.synced/desync and add semantic exit codes ([#98](https://github.com/pantheon-org/iris/issues/98)) ([11cc2b2](https://github.com/pantheon-org/iris/commit/11cc2b25af016e3989104cb97b1fe011bc0f88d9))


### Bug Fixes

* correct opencode Enabled parse bug and remove os.Exit from sync RunE ([#94](https://github.com/pantheon-org/iris/issues/94)) ([743d024](https://github.com/pantheon-org/iris/commit/743d024502a8eac246d95265f1c0e6e05ac08ab5))
* harden config store and startup error handling ([#95](https://github.com/pantheon-org/iris/issues/95)) ([c6584e1](https://github.com/pantheon-org/iris/commit/c6584e142c6acacf8175724b6b2b273fd5fefc1c))
* move symlink check before MkdirAll and compute dynamic column widths ([#96](https://github.com/pantheon-org/iris/issues/96)) ([7501ffc](https://github.com/pantheon-org/iris/commit/7501ffc4e4bd3d11a64de2a970428a7fec2f0368))
* **providers:** revert provider names to short CLI-friendly identifiers ([#92](https://github.com/pantheon-org/iris/issues/92)) ([8d4f211](https://github.com/pantheon-org/iris/commit/8d4f211e531e90fd5b814a8830344435d2137117))
* warn on UserHomeDir failure and clarify ValidateProjectRoot docs ([#99](https://github.com/pantheon-org/iris/issues/99)) ([f2965a7](https://github.com/pantheon-org/iris/commit/f2965a7bbd6eaa605a202e4435d55eb10c3e4308))
* wizard batch saves, rename TerminalRunner, unify PromptConfirm ([#97](https://github.com/pantheon-org/iris/issues/97)) ([2db9901](https://github.com/pantheon-org/iris/commit/2db9901e7b43668786629717783ba0cc6d8f28ff))

## [2026.5.10](https://github.com/pantheon-org/iris/compare/v2026.5.9...v2026.5.10) (2026-05-04)


### Features

* **ci:** auto-update all open PR branches when main is pushed ([#83](https://github.com/pantheon-org/iris/issues/83)) ([fcc39bc](https://github.com/pantheon-org/iris/commit/fcc39bcaa25ac2e3b50da1775a09aa044ef1eda5))
* **cli:** add --provider flag to sync and init subcommands ([#56](https://github.com/pantheon-org/iris/issues/56)) ([8f1d5e1](https://github.com/pantheon-org/iris/commit/8f1d5e17aef729108067ff7295d7664062782c22))
* **cli:** add add and remove commands ([#27](https://github.com/pantheon-org/iris/issues/27)) ([cefd0af](https://github.com/pantheon-org/iris/commit/cefd0af74dc4d291b6cf0fb81eb86ab4a1a96288))
* **cli:** add init command ([#30](https://github.com/pantheon-org/iris/issues/30)) ([9ddd9e9](https://github.com/pantheon-org/iris/commit/9ddd9e9b51f13dbfb0f0acac9eb5ec08a9395494))
* **cli:** add list and status commands ([#26](https://github.com/pantheon-org/iris/issues/26)) ([159b70f](https://github.com/pantheon-org/iris/commit/159b70f4ed93013df714fcec7c4a41ef56ce2b66))
* **cli:** add short flag aliases for all subcommand options ([#44](https://github.com/pantheon-org/iris/issues/44)) ([cd8e4d7](https://github.com/pantheon-org/iris/commit/cd8e4d74a2ae15e2fea577bab48317d81c7fd513))
* **cli:** add sync command ([#28](https://github.com/pantheon-org/iris/issues/28)) ([0d45a66](https://github.com/pantheon-org/iris/commit/0d45a66525dbc0680e7710759ab91c68785c782f))
* **config:** add codec interface and store ([#19](https://github.com/pantheon-org/iris/issues/19)) ([3a3f722](https://github.com/pantheon-org/iris/commit/3a3f722f49b960a3d9f3965f6ec7d0ca2dd57fc4))
* **detector:** add project root provider detector ([#24](https://github.com/pantheon-org/iris/issues/24)) ([2a085f6](https://github.com/pantheon-org/iris/commit/2a085f69190f80a650a405b5999d91651aced02e))
* enhance MCP server configuration with additional fields and preserve remote server attributes ([#61](https://github.com/pantheon-org/iris/issues/61)) ([4ff7393](https://github.com/pantheon-org/iris/commit/4ff739338c60ffe72c55abca6dcbfa5e40370b12))
* **i18n:** allow lang to be set in .iris.json config ([#59](https://github.com/pantheon-org/iris/issues/59)) ([350f9ee](https://github.com/pantheon-org/iris/commit/350f9ee538920b48784fdabe34697c0366e29908))
* **i18n:** internationalise CLI with 14 languages via embedded JSON locales ([#57](https://github.com/pantheon-org/iris/issues/57)) ([c2b799d](https://github.com/pantheon-org/iris/commit/c2b799d66745a0bfd998ff3732ab1d999286ed33))
* **ierrors:** add sentinel errors ([#17](https://github.com/pantheon-org/iris/issues/17)) ([9b46186](https://github.com/pantheon-org/iris/commit/9b461863fba011267649d0b582588d525de6a189))
* **integration:** add end-to-end test and update README ([#31](https://github.com/pantheon-org/iris/issues/31)) ([b5ae7d0](https://github.com/pantheon-org/iris/commit/b5ae7d0ce4328c332b0be4aadfd02b7c02ddb3ca))
* **merger:** add SyncProvider and SyncAllProviders ([#25](https://github.com/pantheon-org/iris/issues/25)) ([13fc1eb](https://github.com/pantheon-org/iris/commit/13fc1eb710f5c07484f95790733058c556d245ce))
* **providers:** add 10 new providers + docs ([#38](https://github.com/pantheon-org/iris/issues/38)) ([4db743c](https://github.com/pantheon-org/iris/commit/4db743c2ae0de0b1599cc92bf4c7763dd9c8fa67))
* **providers:** add Claude and Gemini provider implementations ([#22](https://github.com/pantheon-org/iris/issues/22)) ([4bd5d7d](https://github.com/pantheon-org/iris/commit/4bd5d7d77e4ae33fe4cd3f2d2d15a8bacf1a167b))
* **providers:** add Codex provider implementation ([#21](https://github.com/pantheon-org/iris/issues/21)) ([6d6cf7b](https://github.com/pantheon-org/iris/commit/6d6cf7b04294b57c7a3bf672330f18a262484ff1))
* **providers:** add IntelliJ IDEA provider (.idea/mcp.json) ([#41](https://github.com/pantheon-org/iris/issues/41)) ([3d39508](https://github.com/pantheon-org/iris/commit/3d39508197808563145d90be8b24d31193bbef3d))
* **providers:** add named constants for provider names ([#81](https://github.com/pantheon-org/iris/issues/81)) ([dcf0831](https://github.com/pantheon-org/iris/commit/dcf0831b5ee2672a268d8b0e3bbbe4ec59a6857b))
* **providers:** add OpenCode provider implementation ([#23](https://github.com/pantheon-org/iris/issues/23)) ([8b87974](https://github.com/pantheon-org/iris/commit/8b879741ddf4ced6e66b05d25b8f98b7ab8f29b4))
* **providers:** add Provider interface and Registry ([#20](https://github.com/pantheon-org/iris/issues/20)) ([b9cecfb](https://github.com/pantheon-org/iris/commit/b9cecfbe2e49b6b2b399669139d68fec9766f8f2))
* **scaffold:** add dependencies and wire six no-op cobra subcommands ([#15](https://github.com/pantheon-org/iris/issues/15)) ([afb6350](https://github.com/pantheon-org/iris/commit/afb635079ae60242de4077444853a5884dd099ea))
* switch to CalVer (YYYY.M.PATCH) release versioning ([#3](https://github.com/pantheon-org/iris/issues/3)) ([bd9e668](https://github.com/pantheon-org/iris/commit/bd9e6688c10723066384d80fb1384c16904ce871))
* **types:** add canonical MCPServer and IrisConfig types ([#18](https://github.com/pantheon-org/iris/issues/18)) ([d73f79a](https://github.com/pantheon-org/iris/commit/d73f79a94e93cebd412303e07adb41e3c1ead35e))
* **types:** add MCPServer.Validate to reject invalid transport and URL values ([#72](https://github.com/pantheon-org/iris/issues/72)) ([6d7ccfe](https://github.com/pantheon-org/iris/commit/6d7ccfef235fb02ef402bb7854ca348371764886))
* **wizard:** add Runner interface, ScriptedRunner, and RunInit ([#29](https://github.com/pantheon-org/iris/issues/29)) ([b3725df](https://github.com/pantheon-org/iris/commit/b3725dfbcc0ca6b5eb3ece55385a24f5dfa61e40))
* **wizard:** add URL prompt for SSE transport in server configuration ([#63](https://github.com/pantheon-org/iris/issues/63)) ([e004f35](https://github.com/pantheon-org/iris/commit/e004f3535237efd71f4f1a4c148e96e83f890ef2))
* **wizard:** detect installed harnesses and offer to import servers on init ([#55](https://github.com/pantheon-org/iris/issues/55)) ([8e3ff46](https://github.com/pantheon-org/iris/commit/8e3ff4628b4756c32e619461d344a59aa4dfedec))


### Bug Fixes

* **ci:** run checks on release-please branches to satisfy branch protection ([#34](https://github.com/pantheon-org/iris/issues/34)) ([2ec203c](https://github.com/pantheon-org/iris/commit/2ec203c43ed2425d1d5c0342b71c0be3f78cb498))
* **cli:** improve error handling in RunStatus for missing files ([#64](https://github.com/pantheon-org/iris/issues/64)) ([0a93562](https://github.com/pantheon-org/iris/commit/0a93562348adacf2d84bfe473ba23e5349fafbcc))
* **config:** add mutex to Store.Save and handle temp file cleanup errors ([#68](https://github.com/pantheon-org/iris/issues/68)) ([155ce29](https://github.com/pantheon-org/iris/commit/155ce2936665233c08cf37dfa20003bfd981758b))
* **config:** remove dead Providers field and validate Version on load ([#79](https://github.com/pantheon-org/iris/issues/79)) ([29328ca](https://github.com/pantheon-org/iris/commit/29328cae6a1928a8cd36bd9b7e47c266c2aa8447))
* **detector:** surface IO errors from Exists() instead of silently skipping providers ([#74](https://github.com/pantheon-org/iris/issues/74)) ([c9dfce1](https://github.com/pantheon-org/iris/commit/c9dfce109eda3891ae45b2f6b0d3607ab9caa439))
* **i18n:** propagate load errors and optimise normalize loop ([#76](https://github.com/pantheon-org/iris/issues/76)) ([c786dff](https://github.com/pantheon-org/iris/commit/c786dffeb91c913003ee080c05a846aa6facdb5a))
* **merger:** guard against symlink targets in SyncProvider ([#80](https://github.com/pantheon-org/iris/issues/80)) ([54bd93f](https://github.com/pantheon-org/iris/commit/54bd93fde550dd27257fdb6423ef98226f90dde3))
* **providers:** add missing opencode_expected.json test fixture ([#33](https://github.com/pantheon-org/iris/issues/33)) ([b8fb5e6](https://github.com/pantheon-org/iris/commit/b8fb5e6bac97b45c01a5d6dbf5b58662db62c683))
* **providers:** correct Gemini config path and enable project-level config ([#46](https://github.com/pantheon-org/iris/issues/46)) ([242a839](https://github.com/pantheon-org/iris/commit/242a83968295e15203e3ddd3e145f27d2e289b19))
* **providers:** enable project-level config support for Codex (.codex/config.toml) ([#50](https://github.com/pantheon-org/iris/issues/50)) ([f75b2a6](https://github.com/pantheon-org/iris/commit/f75b2a63bfc45af82b48af29707adc8b14be8104))
* **providers:** enable project-level config support for Mistral Vibe (.vibe/config.toml) ([#52](https://github.com/pantheon-org/iris/issues/52)) ([3224ead](https://github.com/pantheon-org/iris/commit/3224eadec90d8428c2e33217c7b1449b780367b3))
* **providers:** enable project-level config support for Qwen Code (.qwen/settings.json) ([#51](https://github.com/pantheon-org/iris/issues/51)) ([a83ea7a](https://github.com/pantheon-org/iris/commit/a83ea7a5fdd86b86cba7e4d6e3bf70cc89205d18))
* **providers:** validate projectRoot to prevent path traversal in ConfigFilePath ([#73](https://github.com/pantheon-org/iris/issues/73)) ([a56dbdc](https://github.com/pantheon-org/iris/commit/a56dbdc1d6ea9a71a12abe50ae7fdbc840b3b694))
* read GH_APP_ID from vars not secrets in release-please workflow ([#85](https://github.com/pantheon-org/iris/issues/85)) ([b2a8a5d](https://github.com/pantheon-org/iris/commit/b2a8a5d14b7fbff9537599db90d3c4a7b50fc038))
* **registry:** include provider name in Filter error for easier debugging ([#71](https://github.com/pantheon-org/iris/issues/71)) ([3bafcb7](https://github.com/pantheon-org/iris/commit/3bafcb790ad9cc3d032551560a1da4e12e80468e))
* remove bootstrap_sha from release-please manifest ([#4](https://github.com/pantheon-org/iris/issues/4)) ([23afa3d](https://github.com/pantheon-org/iris/commit/23afa3d67406987376e47f57358faffbb2657609))
* resolve misleading file paths in status and sync commands ([#65](https://github.com/pantheon-org/iris/issues/65)) ([40261a8](https://github.com/pantheon-org/iris/commit/40261a87f4080564b484b237371fdfad3e417168))
* set versioning-strategy as workflow input for release-please ([#6](https://github.com/pantheon-org/iris/issues/6)) ([2af9e9c](https://github.com/pantheon-org/iris/commit/2af9e9c293ca2ebaa8c1193f24a9b92c1962fa30))
* **types:** guarantee non-nil IrisConfig.Servers and remove defensive nil checks ([#67](https://github.com/pantheon-org/iris/issues/67)) ([ea52c21](https://github.com/pantheon-org/iris/commit/ea52c212e10861e9decf28019b602b61de961d00))
* use correct versioning key in release-please config ([#10](https://github.com/pantheon-org/iris/issues/10)) ([4a14e75](https://github.com/pantheon-org/iris/commit/4a14e75b2369b773737c3197bc081c9cd98b4593))
* use GitHub App token in release-please workflow to trigger CI checks ([#84](https://github.com/pantheon-org/iris/issues/84)) ([4a5cfd7](https://github.com/pantheon-org/iris/commit/4a5cfd7a6d7ed2c56d537c23e0825eeccdc771c1))
* **wizard:** propagate readFile error in RunInit instead of silent fallback ([#75](https://github.com/pantheon-org/iris/issues/75)) ([3d788a8](https://github.com/pantheon-org/iris/commit/3d788a83597d8bfb8b631e7cb5f7c2b1b3addb5b))

## [2026.5.9](https://github.com/pantheon-org/iris/compare/v2026.5.8...v2026.5.9) (2026-05-04)


### Features

* **ci:** auto-update all open PR branches when main is pushed ([#83](https://github.com/pantheon-org/iris/issues/83)) ([cfb605b](https://github.com/pantheon-org/iris/commit/cfb605b8e4aa7cbc37669306c3e925eb874d7b31))
* **providers:** add named constants for provider names ([#81](https://github.com/pantheon-org/iris/issues/81)) ([50e343b](https://github.com/pantheon-org/iris/commit/50e343b72e1430fcd92189b522ffb41d8e9643b7))
* **types:** add MCPServer.Validate to reject invalid transport and URL values ([#72](https://github.com/pantheon-org/iris/issues/72)) ([3d1db95](https://github.com/pantheon-org/iris/commit/3d1db95c2827733765bc5eb376657081bdc04c57))


### Bug Fixes

* **config:** remove dead Providers field and validate Version on load ([#79](https://github.com/pantheon-org/iris/issues/79)) ([0a0a979](https://github.com/pantheon-org/iris/commit/0a0a979e326f16a09b2cfb909fab81dac72f8ae2))
* **detector:** surface IO errors from Exists() instead of silently skipping providers ([#74](https://github.com/pantheon-org/iris/issues/74)) ([40a2757](https://github.com/pantheon-org/iris/commit/40a275766704063cf5e136c52160f15af45c3669))
* **merger:** guard against symlink targets in SyncProvider ([#80](https://github.com/pantheon-org/iris/issues/80)) ([ac33b5f](https://github.com/pantheon-org/iris/commit/ac33b5f06b3b06edda6a25b150a10484ed21213e))
* read GH_APP_ID from vars not secrets in release-please workflow ([#85](https://github.com/pantheon-org/iris/issues/85)) ([e79e0f2](https://github.com/pantheon-org/iris/commit/e79e0f2eb5e9290be0e5d51ab9e5fa5e23312227))
* use GitHub App token in release-please workflow to trigger CI checks ([#84](https://github.com/pantheon-org/iris/issues/84)) ([c6d2611](https://github.com/pantheon-org/iris/commit/c6d2611e381b587fbfc2bcd784c7e2e3641cbba5))
* **wizard:** propagate readFile error in RunInit instead of silent fallback ([#75](https://github.com/pantheon-org/iris/issues/75)) ([bf10637](https://github.com/pantheon-org/iris/commit/bf10637919760e2ae61029f380a89f41ff8e1e77))

## [2026.5.8](https://github.com/pantheon-org/iris/compare/v2026.5.7...v2026.5.8) (2026-05-02)


### Bug Fixes

* **config:** add mutex to Store.Save and handle temp file cleanup errors ([#68](https://github.com/pantheon-org/iris/issues/68)) ([45049c2](https://github.com/pantheon-org/iris/commit/45049c2b070d6a0c17155e58c84ea620b9aaa35f))
* **i18n:** propagate load errors and optimise normalize loop ([#76](https://github.com/pantheon-org/iris/issues/76)) ([8d81e4f](https://github.com/pantheon-org/iris/commit/8d81e4f0f39b2bba0f96a3fc86bb9225da644502))
* **providers:** validate projectRoot to prevent path traversal in ConfigFilePath ([#73](https://github.com/pantheon-org/iris/issues/73)) ([bbc56af](https://github.com/pantheon-org/iris/commit/bbc56af13116333f21093174bd4b7d99af8b0ae3))
* **registry:** include provider name in Filter error for easier debugging ([#71](https://github.com/pantheon-org/iris/issues/71)) ([dbc6880](https://github.com/pantheon-org/iris/commit/dbc6880e7a48517bb7c5c779a7460f716255557b))
* **types:** guarantee non-nil IrisConfig.Servers and remove defensive nil checks ([#67](https://github.com/pantheon-org/iris/issues/67)) ([af685c0](https://github.com/pantheon-org/iris/commit/af685c0c67598e23dc4fc9598db0752f791e1453))

## [2026.5.7](https://github.com/pantheon-org/iris/compare/v2026.5.6...v2026.5.7) (2026-05-02)


### Features

* enhance MCP server configuration with additional fields and preserve remote server attributes ([#61](https://github.com/pantheon-org/iris/issues/61)) ([403ccf4](https://github.com/pantheon-org/iris/commit/403ccf492da96eb2169a688d804fb203dc7985c4))
* **wizard:** add URL prompt for SSE transport in server configuration ([#63](https://github.com/pantheon-org/iris/issues/63)) ([51341cd](https://github.com/pantheon-org/iris/commit/51341cd9fe41c02b10d85cace277f98937fe9b3e))


### Bug Fixes

* **cli:** improve error handling in RunStatus for missing files ([#64](https://github.com/pantheon-org/iris/issues/64)) ([0a4f503](https://github.com/pantheon-org/iris/commit/0a4f503603608b59aed5a8a5119a505490aeecd1))
* resolve misleading file paths in status and sync commands ([#65](https://github.com/pantheon-org/iris/issues/65)) ([79a1137](https://github.com/pantheon-org/iris/commit/79a11379886bab3c9f02c0546be827fbdd8214c0))

## [2026.5.6](https://github.com/pantheon-org/iris/compare/v2026.5.5...v2026.5.6) (2026-05-01)


### Features

* **i18n:** allow lang to be set in .iris.json config ([#59](https://github.com/pantheon-org/iris/issues/59)) ([01eec69](https://github.com/pantheon-org/iris/commit/01eec6981abfd6b1a00fb7d42fb2c8f16dac404f))
* **i18n:** internationalise CLI with 14 languages via embedded JSON locales ([#57](https://github.com/pantheon-org/iris/issues/57)) ([cb9c46e](https://github.com/pantheon-org/iris/commit/cb9c46e5dd080338618c15c1e2cb7f478f67a029))

## [2026.5.5](https://github.com/pantheon-org/iris/compare/v2026.5.4...v2026.5.5) (2026-05-01)


### Features

* **cli:** add --provider flag to sync and init subcommands ([#56](https://github.com/pantheon-org/iris/issues/56)) ([d5bdd8d](https://github.com/pantheon-org/iris/commit/d5bdd8d8f55233fa03aa72ced48da68a0555b777))
* **wizard:** detect installed harnesses and offer to import servers on init ([#55](https://github.com/pantheon-org/iris/issues/55)) ([87d144c](https://github.com/pantheon-org/iris/commit/87d144ce403673eb465b61e1f7b9605714473edc))


### Bug Fixes

* **providers:** correct Gemini config path and enable project-level config ([#46](https://github.com/pantheon-org/iris/issues/46)) ([0ea5383](https://github.com/pantheon-org/iris/commit/0ea5383620ecbbfcbe1436305c6d1849b39d04f3))
* **providers:** enable project-level config support for Codex (.codex/config.toml) ([#50](https://github.com/pantheon-org/iris/issues/50)) ([c11340b](https://github.com/pantheon-org/iris/commit/c11340b7a938fe928503921c6761aa0067e3b6e6))
* **providers:** enable project-level config support for Mistral Vibe (.vibe/config.toml) ([#52](https://github.com/pantheon-org/iris/issues/52)) ([eaea507](https://github.com/pantheon-org/iris/commit/eaea507cc68ca849abc42eac070afe25c1721eda))
* **providers:** enable project-level config support for Qwen Code (.qwen/settings.json) ([#51](https://github.com/pantheon-org/iris/issues/51)) ([4326257](https://github.com/pantheon-org/iris/commit/4326257726981e55e07ba402cbfe50e0a83110d2))

## [2026.5.4](https://github.com/pantheon-org/iris/compare/v2026.5.3...v2026.5.4) (2026-05-01)


### Features

* **cli:** add short flag aliases for all subcommand options ([#44](https://github.com/pantheon-org/iris/issues/44)) ([99ed9cc](https://github.com/pantheon-org/iris/commit/99ed9ccfe90663b5ac7d85a95d34ae5ebef6a93c))
* **providers:** add IntelliJ IDEA provider (.idea/mcp.json) ([#41](https://github.com/pantheon-org/iris/issues/41)) ([758e6e1](https://github.com/pantheon-org/iris/commit/758e6e1e5171e58f912ebf9e02004c0a59f5fde8))

## [2026.5.3](https://github.com/pantheon-org/iris/compare/v2026.5.2...v2026.5.3) (2026-05-01)


### Features

* **providers:** add 10 new providers + docs ([#38](https://github.com/pantheon-org/iris/issues/38)) ([30874dd](https://github.com/pantheon-org/iris/commit/30874ddc23c035f8d7fa902e6ece8b0ea5afd2c7))

## [2026.5.2](https://github.com/pantheon-org/iris/compare/v2026.5.1...v2026.5.2) (2026-05-01)


### Features

* **cli:** add add and remove commands ([#27](https://github.com/pantheon-org/iris/issues/27)) ([a6785a3](https://github.com/pantheon-org/iris/commit/a6785a3e479456bdb766d7838981bad996b1262d))
* **cli:** add init command ([#30](https://github.com/pantheon-org/iris/issues/30)) ([448bd42](https://github.com/pantheon-org/iris/commit/448bd42ed0fa9909cbf6adbd9af143723398ce45))
* **cli:** add list and status commands ([#26](https://github.com/pantheon-org/iris/issues/26)) ([c3d8007](https://github.com/pantheon-org/iris/commit/c3d80073140245b48db13ac0e3a8493be6e220e2))
* **cli:** add sync command ([#28](https://github.com/pantheon-org/iris/issues/28)) ([671067c](https://github.com/pantheon-org/iris/commit/671067c10904f69d8fc2c06359ef54e120cd6444))
* **config:** add codec interface and store ([#19](https://github.com/pantheon-org/iris/issues/19)) ([7cea957](https://github.com/pantheon-org/iris/commit/7cea95786b8ef0c7aacc5ef679243ce1f469bfe6))
* **detector:** add project root provider detector ([#24](https://github.com/pantheon-org/iris/issues/24)) ([55a6f8a](https://github.com/pantheon-org/iris/commit/55a6f8a73c019cef4908a477861cdd48bcff28ca))
* **ierrors:** add sentinel errors ([#17](https://github.com/pantheon-org/iris/issues/17)) ([46c08e7](https://github.com/pantheon-org/iris/commit/46c08e7e937eee9dfd8959774b57e1274cd7577d))
* **integration:** add end-to-end test and update README ([#31](https://github.com/pantheon-org/iris/issues/31)) ([3c80809](https://github.com/pantheon-org/iris/commit/3c8080977899abe64c5a28ae3c13e78849ce295c))
* **merger:** add SyncProvider and SyncAllProviders ([#25](https://github.com/pantheon-org/iris/issues/25)) ([ed4d1e2](https://github.com/pantheon-org/iris/commit/ed4d1e2ef65b943305342d887b86e1b8d3a0e2bc))
* **providers:** add Claude and Gemini provider implementations ([#22](https://github.com/pantheon-org/iris/issues/22)) ([5e0a44f](https://github.com/pantheon-org/iris/commit/5e0a44f3544a3783b4fcbf5c11f393dbc6d358fa))
* **providers:** add Codex provider implementation ([#21](https://github.com/pantheon-org/iris/issues/21)) ([da7f104](https://github.com/pantheon-org/iris/commit/da7f104b3bdc1f5e6e1112747b9f4f4f1a68d734))
* **providers:** add OpenCode provider implementation ([#23](https://github.com/pantheon-org/iris/issues/23)) ([af46969](https://github.com/pantheon-org/iris/commit/af46969bcf69b2cc6c79c29d7ca55f899bf0a02f))
* **providers:** add Provider interface and Registry ([#20](https://github.com/pantheon-org/iris/issues/20)) ([2fe3627](https://github.com/pantheon-org/iris/commit/2fe3627c47d04e837a9306e1747eec79bd0410d0))
* **scaffold:** add dependencies and wire six no-op cobra subcommands ([#15](https://github.com/pantheon-org/iris/issues/15)) ([3aa00aa](https://github.com/pantheon-org/iris/commit/3aa00aaf0d993fbe04e4cd6be96af05e1dfa732e))
* **types:** add canonical MCPServer and IrisConfig types ([#18](https://github.com/pantheon-org/iris/issues/18)) ([fbfd484](https://github.com/pantheon-org/iris/commit/fbfd48447a73f9a7ec2b07bc7126fcfb8d9d2bb3))
* **wizard:** add Runner interface, ScriptedRunner, and RunInit ([#29](https://github.com/pantheon-org/iris/issues/29)) ([5a403b5](https://github.com/pantheon-org/iris/commit/5a403b5d12222928c9bc45c4f70acfa7e32722d6))


### Bug Fixes

* **ci:** run checks on release-please branches to satisfy branch protection ([#34](https://github.com/pantheon-org/iris/issues/34)) ([8f806f7](https://github.com/pantheon-org/iris/commit/8f806f7da2c6b2a37df37f62a0f9a22fe1777610))
* **providers:** add missing opencode_expected.json test fixture ([#33](https://github.com/pantheon-org/iris/issues/33)) ([44e9ea4](https://github.com/pantheon-org/iris/commit/44e9ea4fd00468b88400507ff72d90b082f960c5))

## [2026.5.1](https://github.com/pantheon-org/iris/compare/v2026.5.0...v2026.5.1) (2026-05-01)


### Features

* switch to CalVer (YYYY.M.PATCH) release versioning ([#3](https://github.com/pantheon-org/iris/issues/3)) ([fe2d73a](https://github.com/pantheon-org/iris/commit/fe2d73aadae7c4fb5e9bb5416942ade033bf61d8))


### Bug Fixes

* remove bootstrap_sha from release-please manifest ([#4](https://github.com/pantheon-org/iris/issues/4)) ([3a80f86](https://github.com/pantheon-org/iris/commit/3a80f86d359e9e7677abfa1a27def3686039eabe))
* set versioning-strategy as workflow input for release-please ([#6](https://github.com/pantheon-org/iris/issues/6)) ([2027cc5](https://github.com/pantheon-org/iris/commit/2027cc5d0433cf6a2e7a4dc76d0d4f47e7ad52f7))
* use correct versioning key in release-please config ([#10](https://github.com/pantheon-org/iris/issues/10)) ([18bea57](https://github.com/pantheon-org/iris/commit/18bea571f52081c931d769b309464efc651319ce))
