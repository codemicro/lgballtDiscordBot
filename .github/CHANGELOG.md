# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

(Dates are in YYYY-MM-DD format. This message is mainly for my own sake.)

## [Unreleased]
### Changed
* `bioFieldType` help text
### Fixed
* Allowed command overloading on bio clear commands

## [3.5.0] - 2021-03-19
### Added
* Help menu
### Changed
* Command parsing system
* Bio commands 
  * `$bio set` and `$bio clear` added to remove complex command parsing logic
### Removed
* `<p>biof` commands
* an outdated section of the bio help text

## [3.4.0] - 2021-03-03
### Added
* `<p>muteme <duration>` command
### Changed
* Removed lingering use of "their" in `<p>pressf`

## [3.3.2] - 2021-02-26
### Changed
* `<p>pressf` no longer uses a pronoun
### Fixed
* Panic on start if JSON data is invalid for embedded CLOC info

## [3.3.1] - 2021-02-17
### Changed 
* Updated to Go version 1.16
* Switch to Docker for deployment

## [3.3.0] - 2021-02-16
### Added
* `<p>steal` command for custom emojis
### Removed
* Custom image for `<p>avatar` command for user with ID equal to `config.OwnerId`
* Reddit account age functionality in verification

## [3.2.0] - 2021-02-15
### Added
* Verification now states the Reddit account age if one is specified
* Admin command to send/edit bot messages 
### Fixed
* Bios for systems now don't fail if the member has been deleted

## [3.1.0] - 2021-02-03
### Added
* Listener confirmation

## [3.0.1] - 2021-01-28
### Fixed
* Account bios now show the correct nickname and avatar for the account (not that of the person who ran the bio command)

## [3.0.0] - 2021-01-28
### Added
* Bios for systems (**BREAKING**)
  * Migration script available - `migration/toV3.py`
* `pkApi` config field

## [2.0.2] - 2021-01-23
### Added
* `<p>shutdown` command
* `ownerId` config field
### Changed
* `<p>pressf` now uses the "her" pronoun for the user with ID `ownerId`

## [2.0.1] - 2021-01-15
### Fixed
* Fix command parsing when dealing with newlines
  * Newlines aren't randomly removed from the middle of a bio field now, for example

## [2.0.0] - 2021-01-15
### Added
* Bio fields now have a built in length limit that reflects that imposed by Discord ([#4](https://github.com/codemicro/lgballtDiscordBot/issues/4))
* A config field for chat chart channel exclusions has been added ([#5](https://github.com/codemicro/lgballtDiscordBot/issues/5))
### Changed
* Modified database schema (**BREAKING**)
  * Migration script available - `migration/toV2.py`
* Commands that are followed by a newline and then the rest of the command are now correctly parsed and processed
### Fixed
* Bio JSON can now forcibly changed when the existing JSON is invalid ([#1](https://github.com/codemicro/lgballtDiscordBot/issues/1))
* `<p>pressf` now only tracks reactions that are 🇫 ([#3](https://github.com/codemicro/lgballtDiscordBot/issues/3))

## [1.8.4] - 2021-01-07
### Added
* Verification ratelimiting
### Fixed
* Fixed improper ping parsing when recording user removals

## [1.8.3] - 2021-01-03
### Fixed
* Fixed message content filter

## [1.8.2] - 2021-01-03
### Changed
* Add filter to new message content to strip extra stuff off the end
  * `hello world submitted by /u/blah [link] [comments]` -> `hello world`

## [1.8.1] - 2021-01-03
### Changed
* Reddit feed watcher now respects defined interval
### Fixed
* Role pings are now filtered by `<p>pressf`

## [1.8.0] - 2021-01-03
### Added
* Number of rotations of the earth since the start time to the info command
* Added Reddit feed config options
* Reddit feed watcher

## [1.7.3] - 2021-01-02
### Added
* `<p>info` command
### Removed
* `<p>broken` command

## [1.7.2] - 2021-01-01
### Fixed
* Users can now only press F once per message.
* Kick/ban log trigger now only works for admin users

## [1.7.1] - 2021-01-01
### Changed
* `core.Bot.SendMessage` now filters @everyone and @here pings... whoops.
* `<p>verifyf` can now only be used by admins

## [1.7.0] - 2021-01-01
### Added
* `<p>pressf` command
* `<p>verifyf` command
### Changed
* Change kick/ban detection to command based system (audit logs were unreliable)
* Tweak message fail commands
* Switch bios to using user nickname when possible

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

[Unreleased]: https://github.com/codemicro/lgballtDiscordBot/compare/v3.5.0...HEAD
[3.5.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v3.4.0...v3.5.0
[3.4.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v3.3.2...v3.4.0
[3.3.2]: https://github.com/codemicro/lgballtDiscordBot/compare/v3.3.1...v3.3.2
[3.3.1]: https://github.com/codemicro/lgballtDiscordBot/compare/v3.3.0...v3.3.1
[3.3.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v3.2.0...v3.3.0
[3.2.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v3.1.0...v3.2.0
[3.1.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v3.0.1...v3.1.0
[3.0.1]: https://github.com/codemicro/lgballtDiscordBot/compare/v3.0.0...v3.0.1
[3.0.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v2.0.2...v3.0.0
[2.0.2]: https://github.com/codemicro/lgballtDiscordBot/compare/v2.0.1...v2.0.2
[2.0.1]: https://github.com/codemicro/lgballtDiscordBot/compare/v2.0.0...v2.0.1
[2.0.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.8.4...v2.0.0
[1.8.4]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.8.3...v1.8.4
[1.8.3]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.8.2...v1.8.3
[1.8.2]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.8.1...v1.8.2
[1.8.1]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.8.0...v1.8.1
[1.8.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.7.3...v1.8.0
[1.7.3]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.7.2...v1.7.3
[1.7.2]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.7.1...v1.7.2
[1.7.1]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.7.0...v1.7.1
[1.7.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v1.6.1...v1.7.0
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
