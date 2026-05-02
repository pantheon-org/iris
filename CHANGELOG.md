# Changelog

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
