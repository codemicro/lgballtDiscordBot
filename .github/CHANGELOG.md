# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

(Dates are in YYYY-MM-DD format. This message is mainly for my own sake.)

## [Unreleased]

## [4.9.5] - 2021-08-23
### Changed
* Update `dgo-toolkit` to fix a restrictions bug

## [4.9.4] - 2021-08-23
### Changed
* Changed the `adminRole` config parameter to `adminRoles` and make it accept a `[]string`

## [4.9.3] - 2021-08-16
### Changed
* Revert changes made to verification system

## [4.9.2] - 2021-08-16
### Changed
* Remove explicit words and food emojis from emojify dataset
* Change density of added emojis to emojify messages
  * This means that only 80% of opportunities to add an emoji to a word will be taken instead of the former 100%
* Verification now outputs pending requests in a separate channel to processed requests to make it easier to see all pending verification requests
  * This now means that `MessageLink` fields of `db.VerificationFail` may be blank when stored in the database
  * Adds the `verificationIds.archiveChannel` config parameter

## [4.9.1] - 2021-08-06
### Removed
* Remove "fuck" and words containing "fuck" from the the emojify dataset

## [4.9.0] - 2021-08-06
### Added
* `<p>emojify` command
### Changed
* Switch to different logging system (from a homebrew, pretty bad system to [Zerolog](https://github.com/rs/zerolog) and [Lumberjack](https://github.com/natefinch/lumberjack))
* `<p>avatar` no longer outputs tiny avatars, and instead will produce avatars at a maximum of 2048x2048 pixels.
* Emoji commands now produce images from the correct Discord domain (`discord.com`, not `discordapp.com`)
### Fixed
* Make edited PluralKit proxied messages show up in the action log
  * Previously, since webhook messages are classed as being sent by bot accounts, they would be indiscriminately ignored.
### Removed
* Functionality to react to PK proxied messages in bio components

## [4.8.8] - 2021-07-30
### Fixed
* Removed goroutine/memory leak from PluralKit API package
  * In some cases when a HTTP request errored out, the response body would not be closed.
  * For each response body, there are two goroutines that are run (`net/http.(*persistConn).readLoop`, `net/http.(*persistConn).writeLoop`).
  * These will run forever and create a goroutine and memory leak unless explicitly stopped by `resp.Body.Close()`.
* Prevent multiple goroutines being started when reporting analytics events instead of just a single goroutine

## [4.8.7] - 2021-07-29
### Fixed
* `<p>restart` and `<p>goroutinestack` can now only be run by the user with the ID `config.OwnerId`
  * Before, it was anyone that *wasn't* that user could run it. Which was very wrong.

## [4.8.6] - 2021-07-29
### Added
* `<p>restart` command
* `<p>goroutinestack` goroutine stack dump command

## [4.8.5] - 2021-07-28
### Changed
* Replaced old DB-based analytics with Prometheus-based monitoring
### Removed
* `db.AnalyticsEvent` type

## [4.8.4] - 2021-07-03
### Fixed
* Updated ping detection regex to detect pings in the format `<@userid>`

## [4.8.3] - 2021-07-03
### Changed
* Ensure that all generated user pings are in the format `<@userid>`
### Fixed
* Applied linter recommendations

## [4.8.2] - 2021-06-30
### Fixed
* Reddit: fixed use of loop variable in a goroutine
  * This could cause a given Reddit feed to have multiple monitors running on it, causing individual posts to be posted multiple times 

## [4.8.1] - 2021-06-30
### Fixed
* Verification reaction handler: add handling for users that have left before being verified
* Attempted to fix Reddit rate-limiting problems
  * Stagger subreddit monitor starts
  * Add proper user agent
  * Add delay and retry for the first rate-limit response per action run
  * Don't run feed monitors when running in debug mode

## [4.8.0] - 2021-06-24
### Added
* `<p>spoiler` command

## [4.7.7] - 2021-06-20
### Fixed
* Newlines are no longer replaced with spaces in verification messages
* Remove useless steps from verification pronoun detection

## [4.7.6] - 2021-06-08
### Changed
* Add `LABEL com.centurylinklabs.watchtower.stop-signal="SIGINT"` to Dockerfile
* Replace occurrences of `0x414b` with `0xAb1`

## [4.7.5] - 2021-06-06
### Added
* Added channel link to `#rules-welcome` on verification fail text
### Fixed
* Bios now show even if the associated PluralKit system cannot be found

## [4.7.4] - 2021-06-05
*No changes, CI test release*

## [4.7.3] - 2021-06-05
### Fixed
* Forced verifications no longer strip the first word off the message

## [4.7.2] - 2021-06-05
### Added
* Tone tag delete command (`<p>tonetag delete`)
### Fixed
* Admin commands page of the help embed will now show
  * This was due to the shutdown command not having any help text
  * This has been mitigated by adding filler text in place of the command description if none is provided
* Update DiscordGo patch to correctly copy `*Message` instances

## [4.7.1] - 2021-06-01
### Fixed
* `<p>uwu` now converts `l>w` and `r>w`

## [4.7.0] - 2021-06-01
### Added
* Tone tag lookup commands
  * `<p>toneTag lookup <shorthand>`, `<p>toneTag list`, `<p>toneTag new <shorthand> <description>`
  * Also adds `toV47.py` migration script to automatically fill the database with a set of tone tags prior to first use
* `<p>uwu`
### Changed
* The previous incident counter streak is now shown when the `<p>indicents reset` command is run
* Improved chatchart chart and response formatting

## [4.6.1] - 2021-05-16
### Fixed
* Verification error message actually refers to the right thing

## [4.6.0] - 2021-05-16
### Added
* Basic command and PluralKit request analytics
  * Adds `analytics_events` table to database
### Changed
* Verification: updated error message when verification text is missing
* Bios: update system help dialog to include information about new bio picker introduced in [4.0.0].
* Move bio help dialog text to external Markdown files for easy editing
  * This also introduces the `internal/markdown` package
### Fixed
* Checks for PluralKit proxied messages in bios no longer generate warnings for messages not found (404) 
  * A 404 just means that the message wasn't proxied

## [4.5.1] - 2021-05-14
### Changed
* Bio commands that deal with adding reactions/deleting messages now interact with webhook messages sent by PluralKit if possible (as opposed to the account message that has subsequently been deleted)

## [4.5.0] - 2021-05-13
### Added
* `<p>shutdown` command
### Changed
* Reduce amount of clutter in action log messages
  * Adds filters for when message content or author information cannot be received and changes the log output accordingly
* Improve `state.State` API
  * Removed `state.(*state).AddGoroutine` and inlined added its functionality into `state.(*state).WaitUntilShutdownTrigger`

## [4.4.0] - 2021-05-12
### Added
* `<p>incidents` and `<p>incidents reset` commands
### Changed
* `<p>purgeUnverified` now confirms before executing any kicks

## [4.3.0] - 2021-05-09
### Added
* Action log
  * Alongside `actionLogChannel` config parameter

## [4.2.0] - 2021-05-05
### Added
* `<p>purgeUnverified`
* `verificationIds.extraValidRoles` config parameter

## [4.1.4] - 2021-04-23
### Fixed
* `<p>verifyf` no longer errors when run
  * When retrieving a message using `discordgo.(*Session).ChannelMessage`, the `GuildID` field is not set. This was used by `verification.(*Verification).coreVerify` to retrieve guild roles.
  * When a message is received from a websocket update, it would have `GuildID` and be okay
  * Fixed by manually adding the guild ID to the `Message` instance in question.
  * Introduced in [3.7.0]

## [4.1.3] - 2021-04-22
### Fixed
* Role positions are now properly taken into account when determining what colour to show on bio embeds

## [4.1.2] - 2021-04-22
### Fixed
* Bio embed colours are now taken from higher priority roles first
  * Previously the colour of the lowest priority role with a colour would be used

## [4.1.1] - 2021-04-21
### Fixed
* Bios now always show the user colour, irrespective of if the top assigned role has a colour or not (it wouldn't do before)

## [4.1.0] - 2021-04-21
### Added
* User colours to bios
* Changelog is now included in `<p>info` command
  * This is done by uploading the changelog to a 3rd party service before building
### Changed
* Logging outputs now go to `os.Stderr` in all cases
### Fixed
* Bios for systems no longer fail to load any bios if the PluralKit system member list is private
* Left-hand bio carousel scroll reaction no longer does nothing if the carousel didn't start at the first bio

## [4.0.0] - 2021-04-20
### Added
* Added initial selection to bios (**BREAKING**)
  * For accounts with multiple bios registered, the bio to start at in the carousel is now chosen by the user
* `<p>forgetme` command
* `<p>mydata` command
### Changed
* Switch to using hashed user IDs in kick, ban and verification fail tables (**BREAKING**)
  * This means this data can be retained while complying with a request from a user to delete all identifying data
  
## [3.7.2] - 2021-04-11
### Fixed
* Remove erroneous `\t` in role ping regexp that was breaking verification pronoun roles

## [3.7.1] - 2021-04-11
### Added
* Filtering to verification pronoun roles

## [3.7.0] - 2021-04-11
### Added
* Pronouns are now automatically detected in verification messages and the relevant roles applied on verification acceptance.

## [3.6.2] - 2021-04-08
### Changed
* Switch verification log messages to use embeds and pings instead of plaintext and weird Base64 encoded data
### Fixed
* Chatchart messages no longer fail to send when the requesting message has been deleted (eg. if a PluralKit proxy was used)

## [3.6.1] - 2021-04-06
### Changed
* Test change of command precedence for potential bug fix?

## [3.6.0] - 2021-04-02
### Changed
* Only the user that requested the bio of another user can control the carousel generated
* User errors (those caused by something the user did) in commands now use `route.(*MessageContext).SendErrorMessage`
### Fixed
* Custom emoji regular expression now detects multiple custom emojis in one message
  * This fixes `<p>steal` by extension
* Pings in the format `<@userid>` are correctly detected

## [3.5.3] - 2021-04-01
### Changed
* Build: `v` prefix is trimmed from version numbers

## [3.5.2] - 2021-04-01
### Changed
* Switched to DiscordGo for webhook execution
### Fixed
* Reddit RSS feeds no longer cause a panic if no publish time is provided

## [3.5.1] - 2021-03-19
### Changed
* `bioFieldType` help text
### Fixed
* Allowed command overloading on bio clear commands
* Verification now requires a message again

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

[Unreleased]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.9.5...HEAD
[4.9.5]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.9.4...v4.9.5
[4.9.4]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.9.3...v4.9.4
[4.9.3]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.9.2...v4.9.3
[4.9.2]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.9.1...v4.9.2
[4.9.1]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.9.0...v4.9.1
[4.9.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.8.8...v4.9.0
[4.8.8]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.8.7...v4.8.8
[4.8.7]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.8.6...v4.8.7
[4.8.6]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.8.5...v4.8.6
[4.8.5]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.8.4...v4.8.5
[4.8.4]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.8.3...v4.8.4
[4.8.3]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.8.2...v4.8.3
[4.8.2]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.8.1...v4.8.2
[4.8.1]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.8.0...v4.8.1
[4.8.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.7.7...v4.8.0
[4.7.7]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.7.6...v4.7.7
[4.7.6]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.7.5...v4.7.6
[4.7.5]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.7.4...v4.7.5
[4.7.4]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.7.3...v4.7.4
[4.7.3]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.7.2...v4.7.3
[4.7.2]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.7.1...v4.7.2
[4.7.1]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.7.0...v4.7.1
[4.7.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.6.1...v4.7.0
[4.6.1]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.6.0...v4.6.1
[4.6.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.5.1...v4.6.0
[4.5.1]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.5.0...v4.5.1
[4.5.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.4.0...v4.5.0
[4.4.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.3.0...v4.4.0
[4.3.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.2.0...v4.3.0
[4.2.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.1.3...v4.2.0
[4.1.4]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.1.3...v4.1.4
[4.1.3]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.1.2...v4.1.3
[4.1.2]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.1.1...v4.1.2
[4.1.1]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.1.0...v4.1.1
[4.1.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v4.0.0...v4.1.0
[4.0.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v3.7.2...v4.0.0
[3.7.2]: https://github.com/codemicro/lgballtDiscordBot/compare/v3.7.1...v3.7.2
[3.7.1]: https://github.com/codemicro/lgballtDiscordBot/compare/v3.7.0...v3.7.1
[3.7.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v3.6.2...v3.7.0
[3.6.2]: https://github.com/codemicro/lgballtDiscordBot/compare/v3.6.1...v3.6.2
[3.6.1]: https://github.com/codemicro/lgballtDiscordBot/compare/v3.6.0...v3.6.1
[3.6.0]: https://github.com/codemicro/lgballtDiscordBot/compare/v3.5.3...v3.6.0
[3.5.3]: https://github.com/codemicro/lgballtDiscordBot/compare/v3.5.2...v3.5.3
[3.5.2]: https://github.com/codemicro/lgballtDiscordBot/compare/v3.5.1...v3.5.2
[3.5.1]: https://github.com/codemicro/lgballtDiscordBot/compare/v3.5.0...v3.5.1
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
