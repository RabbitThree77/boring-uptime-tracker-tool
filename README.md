# Boring Uptime Tracker Tool (BUTT)
A lightweight, configurable uptime tracker written in Go

> [!WARNING]
> This repo is new and under active development.
> I do not recommend using it in serious projects.

## Features
- tiny size
- made to be simple
- concurrency
- local only
- fully configurable in TOML
- Discord Webhook Notifications

## Quick Start
1. **Clone Repository:**
  ```bash
  git clone https://github.com/RabbitThree77/boring-uptime-tracker-tool
  ```
2. **Configure:**
  Edit `Config.toml`
3. **Run The Script:**
  `go run main.go`

## Configuration
A short guide on configuration, sample configuration provided in repo
### Server-wide Configuration
```TOML
[Server]
verbose = false # a flag that decides whether or not to print extra infomration for each request [true | false]
timeout = 1.0   # the value that decides how long to wait until a request fails [min 1, float]

[Notifications]
discord_webhook = "BUTT_DISCORD_WEBHOOK" # the name of an enviromental variable that holds a discord webhook url
```
### Individual Website Settings
```TOML
[[Websites]]
name = "Google"            # the arbitrary display name of the website [string]
url = "https://google.com" # the URL that is tested for uptime [string]
interval = 15.0            # the interval between requests, suggest to keep as high as possible [float]
retry = 0                  # the amount of times to try again before reporting an outage to the logs [integer]
```

## Roadmap
- [x] implement external means of reporting outages (discord)
- [ ] create more customization settings
- [ ] store historical uptime data
- [ ] create an interface to view uptime data
- [ ] implement external means of reporting outages (other)
- [ ] ...
