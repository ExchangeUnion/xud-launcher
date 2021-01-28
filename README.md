xud-launcher
============

The xud-launcher is a thin wrapper of xud-docker launcher which enables running any branch of xud-docker since version 2. It will keep a low update frequency, and it will be embedded in our GUI and CLI applications. 

### Build

On *nix platform
```sh
make
```

On Windows platform
```
mingw32-make
```

### Run

On *nix platform
```sh
export BRANCH=master
export NETWORK=mainnet
./xud-launcher setup
```

On Windows platform (with CMD)
```
set BRANCH=master
set NETWORK=mainnet
./xud-launcher setup
```

On Windows platform (with Powershell)
```
$Env:BRANCH = "master"
$Env:NETWORK = "mainnet"
./xud-launcher setup
```
