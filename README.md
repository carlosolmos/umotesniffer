# umotesniffer

Connects to the TCP servers of two micromotes to inspect the cot traffic.

It runs in your console.


![Alt text](./umotesniffer1.png?raw=true "Umote Sniffer in gui mode")


# install

1. download the binary from the releases page https://umotesniffer/releases
2. create the configuration file (see below)

# build

```azure
$ go build -o umotesniffer .
```

# Run

```azure
$./umotesniffer 
Tool to probe the umote network.

Usage:
  umotesniffer [command]

Available Commands:
  gui         Launch in GUI mode
  help        Help about any command
  version     Shows the version of the system

Flags:
      --config string   config file (default is $HOME/umotesniffer.properties)
  -h, --help            help for umotesniffer

Use "umotesniffer [command] --help" for more information about a command.

```

# Launch in Console UI mode

```azure
$./umotesniffer gui --config config.properties
```

# Config file

`config.properties`

define the host:port and alias for two umotes (top and bottom)
```azure
log_level=debug
log_out=file

TopHost=127.0.0.1:9088
TopAlias="Field 00000002"

BottomHost=127.0.0.1:8088
BottomAlias="GW 00000001"

```
