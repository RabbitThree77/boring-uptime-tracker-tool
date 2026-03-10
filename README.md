# Boring Uptime Tracker Tool (BUTT)
A TOML configurable uptime tracker written in Go

> [!CAUTION]
> This repo is new and under active development.
> I do not recommend using it in its current state.
> BUTT doesn't yet provide any solid way to detect an outage other than text logs

## Features
- fully configure the tracker in TOML
- Customize
  - interval
  - Retry Count

## Quick Start
1. **Clone Repository**
  ```bash
  git clone https://github.com/RabbitThree77/boring-uptime-tracker-tool
  ```
2. **Configure**
  Edit `Config.toml`
3. **Run The Script**
  `go run main.go`

## Configuration
A short guide on configuration, sample configuration provided in repo
### Server-wide Configuration
```TOML
[Server]
verbose = false # a flag that decides whether or not to print extra infomration for each request [true | false]
timeout = 1.0   # the value that decides how long to wait until a request fails [min 1, float]
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
- [ ] create more customization settings
- [ ] store historical uptime data
- [ ] create an interface to view uptime data
- [ ] implement external means of reporting outages (discord, slack, ...)
- [ ] ...
