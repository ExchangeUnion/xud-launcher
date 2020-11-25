xud-launcher
============

```sh
go build

# Generate docker-compose.yml file
./xud-launcher -n simnet gen

# Delegate to docker-compose commands
./xud-launcher -n simnet up -d
./xud-launcher -n simnet down
./xud-launcher -n simnet start xud
./xud-launcher -n simnet stop xud
./xud-launcher -n simnet restart xud
./xud-launcher -n simnet logs xud
```