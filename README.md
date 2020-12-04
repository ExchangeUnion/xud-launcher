xud-launcher
============

### Build

```sh
go build .
```

### Run synthetic commands

```sh
./xud-launcher setup
./xud-launcher cleanup
```

### Generate `docker-compose.yml` and `config.json`

```sh
./xud-launcher gen
```

### Run docker-compose commands

```sh
./xud-launcher up -- -d
./xud-launcher down
./xud-launcher start xud
./xud-launcher stop xud
./xud-launcher restart xud
./xud-launcher logs xud
```
