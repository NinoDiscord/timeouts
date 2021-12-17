# 🕰️ @nino/timeouts
> **Microservice to relay and dispatch timeouts.**

## Features
- :heart: **Easy to Use** — You can create a WebSocket client for this in ~1 minute or so.
- :zap: **Robust** — Made to be fast and performant

## Why?
Originally, this was made for Nino and it should be small and easy to use, but as the bot *will* grows,
it will need to dispatch more and more on different nodes and such, so **v2** and onwards will persist and
apply timeouts on every server restart! In **v1**, you would have to tell the server what needs to be applied.

## Installation
### Prerequisites 
Before you can run the **timeouts** service, you will need the following:

- **Redis v6.2+**
- **Go v1.17+**

Optional tools:

- **Docker**
- **Sentry**

#### Locally
```shell
# 1. Pull the repository
$ git clone https://github.com/NinoDiscord/timeouts && cd timeouts

# 2. Build the project
# If you're using Windows, install Make with chocolatey: `choco install make`
$ make build

# 3. Run the project
$ ./build/timeouts # Attach `.exe` at the end of this command if using Windows.
```

#### Docker
You can use the image we provide on the GitHub Container Registry:

```shell
# 1. Pull the image down to your system
$ docker pull ghcr.io/ninodiscord/timeouts/timeouts:latest

# 2. Run the image
$ docker run -d -p 4250:4250 -e AUTH=... ghcr.io/ninodiscord/timeouts/timeouts:latest
```

## License
**@nino/timeouts** is released under the **MIT** License, read [here](/LICENSE) for more information.
