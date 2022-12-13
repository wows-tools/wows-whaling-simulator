# Wows Whaling Simulator

Webapp + API + CLI tool simulating whaling sessions with lootboxes

## Status

Work in Progress

## API

build

```
make
```

launch:

```
export WOWS_WOWSAPIKEY=zzzzzzzzzzzzzzzzzzzzz
./wows-whaling-simulator
```

TODO

## CLI tool

To use the CLI tool you need a Wargaming API key: https://developers.wargaming.net/applications/

Once you got the API key (don't forget to whitelist your public IP), do the following:

```
export WOWS_WOWSAPIKEY=zzzzzzzzzzzzzzzzzzzzz

go run misc/cli/main.go --help
Usage of /tmp/go-build4003606908/b001/exe/main:
  -n int
    	number of containers opened (default 10)
  -nick string
    	Nickname of the player
  -realm string
    	Wows realm (eu, na, asia) (default "eu")


go run misc/cli/main.go -n 20 -nick "kakwa" -realm eu -lootbox santa_mega_2022 | sort -n
```

It will output something similar to:

```
     1 Item: Béarn
     1 Item: Collingwood
     1 Item: FR25
     1 Item: Karl von Schönberg
     1 Item: Rio de Janeiro
     1 Item: Toulon
     5 Item: Free XP Red Economic Bonus
    20 Item: Ship XP Blue Economic Bonus
    40 Item: Free XP Blue Economic Bonus
    60 Item: Commander XP Blue Economic Bonus
   120 Item: Commander XP Green Economic Bonus
   120 Item: Credits Green Economic Bonus
   180 Item: Free XP Green Economic Bonus
   180 Item: Premium Days
   300 Item: Ship XP Green Economic Bonus
  2500 Item: Doubloons
 37500 Item: Coal
```
