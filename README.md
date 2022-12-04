# Wows Whaling Simulator

Webapp + API + CLI tool simulating whaling sessions with lootboxes

## Status

Work in Progress

## API

TODO

## CLI tool

To use the CLI tool you need a Wargaming API key: https://developers.wargaming.net/applications/

Once you got the API key (don't forget to whitelist your public IP), do the following:

```
export WG_API_KEY=zzzzzzzzzzzzzzzzzzzzz

go run misc/cli/main.go --help
Usage of /tmp/go-build4003606908/b001/exe/main:
  -n int
    	number of containers opened (default 10)
  -nick string
    	Nickname of the player
  -realm string
    	Wows realm (eu, na, asia) (default "eu")


go run misc/cli/main.go -n 20 -nick "kakwa" -realm eu | sort -n
```

It will output something similar to:

```
     1 Item: Cossack
     1 Item: Mikhail Kutuzov
     5 Item: Spring Sky
    20 Item: New Year Streamer
    25 Item: Type 3 — New Year
    30 Item: Basilisk
    30 Item: Leviathan
    30 Item: Ouroboros
    30 Item: Red Dragon
    30 Item: Wyvern
    60 Item: New Year Sky camouflages
    60 Item: Scylla
  5000 Item: Doubloons
 37500 Item: Coal
```
