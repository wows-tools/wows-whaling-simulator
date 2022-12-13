# Wows Whaling Simulator

Webapp + API + CLI tool simulating whaling sessions with lootboxes

## API

## Dependencies

* golang
* node+yarn

## Build

To build everything:
```bash
make
```

To build just the API + CLI part:

```bash
# API
make wows-whaling-simulator

# CLI
make wows-whaling-simulator-cli
```

To build the ui part:

```bash
cd ui
yarn build
```

# Run

```
export WOWS_WOWSAPIKEY=zzzzzzzzzzzzzzzzzzzzz
./wows-whaling-simulator
```

## CLI tool

To use the CLI tool you need a Wargaming API key: https://developers.wargaming.net/applications/

Once you got the API key (don't forget to whitelist your public IP), do the following.

Help:

```
export WOWS_WOWSAPIKEY=zzzzzzzzzzzzzzzzzzzzz

./wows-whaling-simulator-cli -help
Usage of ./wows-whaling-simulator-cli:
  -lootbox string
    	Lootbox type
  -n int
    	number of containers opened
  -nick string
    	Nickname of the player
  -realm string
    	Wows realm (eu, na, asia) (default "eu")
  -target string
    	ship targeted (exclusive with -n


```

Quantity Run:
```
wows-whaling-simulator-cli -n 5 -lootbox santa_big_2022 -realm eu -nick kakwa

```

Target Run:
```
wows-whaling-simulator-cli -t Enterprise -lootbox santa_big_2022 -realm eu -nick kakwa

```

It will output something similar to:

```json
{
 "container_opened": 5,
 "pities": 0,
 "game_money_spent": 3750,
 "euro_spent": 12.15,
 "dollar_spent": 13.36,
 "collectables_items": null,
 "other_items": [
  {
   "quantity": 20,
   "name": "Commander XP Blue Economic Bonus",
   "id": "commanderxp_blue_economic_bonus",
   "attributes": {
    "type": "economic bonus"
   }
  },
  {
   "quantity": 75,
   "name": "Victor Lima",
   "id": "victor_lima",
   "attributes": {
    "type": "signal"
   }
  },
  {
   "quantity": 75,
   "name": "India Delta",
   "id": "india_delta",
   "attributes": {
    "type": "signal"
   }
  },
  {
   "quantity": 75,
   "name": "Sierra Mike",
   "id": "sierra_mike",
   "attributes": {
    "type": "signal"
   }
  }
 ]
}
```
