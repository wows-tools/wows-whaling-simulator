# Wows Whaling Simulator

Webapp + API + CLI tool simulating massive whaling with World of Warships lootboxes.

## Build Dependencies

- golang
- node+yarn

## Run dependencies

- redis (API)

## Build

To build everything:

```bash
make
```

To build just the API part:

```bash
# API
make wows-whaling-simulator
```

To build just the CLI part:

```bash
# CLI
make wows-whaling-simulator-cli
```

To build the ui part:

```bash
# React UI
make -C ui

# alternatively
cd ui
yarn build
```

## Run

### Wows API key

To use these tools/api you need a Wargaming API key: https://developers.wargaming.net/applications/

Don't forget to whitelist your public IP.

### API

#### Config

The default configuration file `config.yml` should work well enough to test the API.

But if you need to tweak it:

```bash
vim config.yml
```

The configuration loader is set to open the configuration file in either of the following location:

- `/etc/wows-whaling/config.yml`
- `$HOME/.wows-whaling/config.yml`
- `./config.yml`

additionally, each configuration parameter can be overriden by a corresponding environment variables:

```bash
export WOWS_WOWSAPIKEY=zzzzzzzzzzzzzzzzzzzzz
export WOWS_LISTEN=":8080"
export WOWS_REDISURI="redis://localhost:6379/"
export WOWS_LOGLEVEL="INFO"
export WOWS_STATICASSETS="./static"
export WOWS_RATES="./rates"
```

#### Launch

```bash
#export the API key (note: it can also be set in the configuration file)
export WOWS_WOWSAPIKEY=zzzzzzzzzzzzzzzzzzzzz

# launch the API
./wows-whaling-simulator
```

### CLI tool

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
