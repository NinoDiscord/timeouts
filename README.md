# Timeouts

> 🕰️ **Manages mute/ban timeouts**

## Installation

There's a couple of ways to install it

### Docker
```
docker pull ghcr.io/ninodiscord/timeouts/timeouts:latest
``` 

### Get it from git clone
```
git clone https://github.com/NinoDiscord/timeouts.git && cd timeouts
```

## Running Timeout Services

As of right now you'll need to provide the AUTH as a system enviromental variable before running it

### Linux 

```
export AUTH=<Value>
go run main.go
```
