# gershik-feed-bot
A Telegram feed bot made mainly for my friend [@gershik](https://t.me/gershik).

## What is it?
This bot polls different Telegram channels and collects all its messages into a single feed.
You can create multiple feeds, each with its own channels.
Also, there is a filtering feature, so you can filter out messages that contain certain text.

## Getting started
1. Create a bot using [@BotFather](https://t.me/BotFather).
2. Grab your API ID and API hash from [here](https://my.telegram.org/apps).
3. Get the binary from [releases](https://github.com/altfoxie/gershik-feed-bot/releases) or build it yourself.
4. Run it without any arguments to generate a config file.
5. Fill in the config file `config.yml`.
> **Note:** You can get your Telegram ID by sending `/getid` to [@myidbot](https://t.me/myidbot).
6. Run it again and enjoy!

## Building
1. Install [Go](https://go.dev).
2. Clone this repository.
3. Run `go build` in the repository directory.

## Configuration file
### Bot
| Name | Type | Description |
| --- | --- | --- |
| token | string | Bot token. |
| whitelist | []int64 | List of users who can use the bot. |

### User
| Name | Type | Description |
| --- | --- | --- |
| session_path | string | Path to the session json file. |
| app_id | int | API ID from [here](https://my.telegram.org/apps). |
| app_hash | string | API hash from [here](https://my.telegram.org/apps). |

### Poller
| Name | Type | Description |
| --- | --- | --- |
| storage_path | string | Path to the storage json file. |
| inteval | int | Polling interval in seconds. |
| channels_interval | int | Interval between polling different channels in seconds. |
| mark_as_unread | bool | Mark messages as unread after forwarding. |
| filters | []Filter | List of filters. |

### Filter object
| Name | Type | Description |
| --- | --- | --- |
| channels | []int64 | List of channels IDs where the filter should be applied. |
| text | string | Text to filter, not case sensitive. |