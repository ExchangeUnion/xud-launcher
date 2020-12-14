xud-launcher
============

The xud-launcher translates xud-docker configuration into a docker-compose.yml file and provides a bunch of commands to
tweak your XUD setup.

- `setup`: brings up an XUD environment
- `cleanup`: remove the whole XUD environment (with all files)
- `info`: show basic launcher information
- `version`: show the launcher version
- `attach`: attach this launcher to an XUD environment (not implemented yet)
- `console`: start up your current shell like xud-docker xud-ctl (not implemented yet)

More advanced commands comes with:

- `gen`: generates a docker-compose.yml file
- The docker-compose wrapper commands `up`, `down`, `start`, `stop`, `restart`, `logs`, `exec`, `pull`

### Build

```sh
go build .
```

### Run

```sh
./xud-launcher setup
```
