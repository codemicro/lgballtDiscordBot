# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

(Dates are in YYYY-MM-DD format. This message is mainly for my own sake.)

## [Unreleased]

## [1.6.1] - 2020-12-30
### Added
* Rejected verification and ban/kick check to verification tools

## [1.6.0] - 2020-12-29
### Added
* `verification` component
* New user verification commands
### Changed
* Message link detection regular expression
* Config file format

## [1.5.0] - 2020-12-29
### Added
* `<p>emoji` command

## [1.4.1] - 2020-12-27
Now it compiles!

## [1.4.0] - 2020-12-27
### Added
* `misc` components
* `<p>avatar` command
* `tools.ParsePing`

## [1.3.3] - 2020-12-27
### Added
* Debug mode
### Changed
* Optimise check for own user ID using `sync.Once`
### Fixed
* Bot now has an exclusion for itself on reaction add/remove
* Emojis are now validated when registering new role reactions

## [1.3.2] - 2020-12-27
### Changed
* Switch chat charts to bar charts from pie charts

## [1.3.1] - 2020-12-27
### Changed
* Chatchart user inclusion threshold

## [1.3.0] - 2020-12-27
### Added
* `chartchart` component
* Chatchart commands
### Fixed
* Panic on new reaction role registration with standard emoji

## [1.2.2] - 2020-12-26
### Changed
* Bot no longer accepts commands via DM
### Fixed
* Bot now reacts with target emoji when a new reaction role is registered

## [1.2.1] - 2020-12-26
### Fixed
* Fix panic on bad number of arguments to role track command

## [1.2.0] - 2020-12-26
### Added 
* `roles` component
* Role reaction tools
* `broken` command

## [1.1.1] - 2020-12-23
### Added
* Admin bio modification commands to force set/clear user bios

## [1.1.0] - 2020-12-23
### Changed
* Update bio help text
* Switch to SQLite database

## [1.0.3] - 2020-12-03
### Fixed
* Fix disconnect and deadlock on gateway reconnect command (Harmony update)

## [1.0.2] - 2020-12-02
### Added
* Bio field clearing functionality
* `info` component
* `<p>info prefix` command
### Changed
* Large internal refactor
### Fixed
* Command parsing now correctly splits arguments


## [1.0.1] - 2020-11-30
### Changed
* Updated bios help message
### Fixed
* Added bot exclusion (bots cannot trigger this bot, nor can PluralKit proxies)

## [1.0.0] - 2020-11-30
* Initial release with `bio` component

[Unreleased]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.6.1...HEAD
[1.6.1]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.6.0...v1.6.1
[1.6.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.5.0...v1.6.0
[1.5.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.4.1...v1.5.0
[1.4.1]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.4.0...v1.4.1
[1.4.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.3.3...v1.4.0
[1.3.3]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.3.2...v1.3.3
[1.3.2]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.3.1...v1.3.2
[1.3.1]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.3.0...v1.3.1
[1.3.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.2.2...v1.3.0
[1.2.2]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.2.1...v1.2.2
[1.2.1]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.2.0...v1.2.1
[1.2.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.1.1...v1.2.0
[1.1.1]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.1.0...v1.1.1
[1.1.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.0.3...v1.1.0
[1.0.3]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.0.2...v1.0.3
[1.0.2]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/codemicro/lgballtDiscordBot/releases/tag/v1.0.0