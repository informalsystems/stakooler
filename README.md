# Stakooler

The koolest (light) tool for Cosmos stakers !

## Installation

* Install a recent version of [Golang installed on the machine](https://go.dev/doc/install)
* Clone this repository
* Build the tool with `go build`

## Configuration

stakooler needs a configuration file in order to run.

create and save a `config.toml` file in the current directory or under `$HOME/.stakooler/config.toml`

> You can also save in a different location and use the `--config [config file full path]`

Sample `config.toml`

```toml
########################
# Accounts             #
########################

[[accounts]]
name = "[account name]" # name of this account
address = "[cosmos address bech-32]" # account address
chain = "[a chain id matching one from the chains section]" # this should match the chain id of one of the chains configured

########################
# Validators           #
########################

[[validators]]
valoper = "[validator's valoper address]"
chain = "[a chain id matching one from the chains section]"

########################
# Chains               #
########################

[[chains]]
id = "[chain id]" # chain-id
lcd = "[lcd address of the node]" # the REST endpoint of the node e.g. http://myosmonode.com:1317


########################
# Zabbix - enabled when running with the --zbx flag
########################

[zabbix]
server = "[IP/URL]" # Zabbix server IP or URL
host = "[zabbix host]" # Host defined in zabbix with a trapper item
port = "[zabbix trapper port]" # Port used by Zabbix server for trapper items. Default 10051
```

> the [chain id] has to match one that is available in the [Cosmos Directory](https://cosmos.directory). Select the chain, in the tab Chain look for the  Chain ID property

## Running

### Accounts Details

In order to fetch accounts details (the ones from the configuration file) use:

```stakooler accounts details```

This will show balance, rewards, staked and unbonding tokens for each account

### Validator Statistics

In order to fetch the validator's statistics across different chains use:

```stakooler validator stats```

> NOTE: Still Work In Progress (WIP)! First release coming soon!
