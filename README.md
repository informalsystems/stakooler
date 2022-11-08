# Stakooler

The koolest (light) tool for stakers !

### Installation

* Install a recent version of [Golang installed on the machine](https://go.dev/doc/install)
* Clone this repository
* Build the tool with `go build`

### Configuration

stakooler needs a configuration file in order to run.

create and save a `config.toml` file in the current directory or under `$HOME/.stakooler/config.toml`

`config.toml`
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
valoper = "cosmosvaloper16k579jk6yt2cwmqx9dz5xvq9fug2tekvlu9qdv"
chain = "cosmoshub"

########################
# Chains               #
########################

[[chains]]
id = "[chain id]" # chain-id
lcd = "[lcd address of the node]" # the REST endpoint of the node e.g. http://myosmonode.com:1317
```

### Running

#### Accounts Details

In order to fetch accounts details (the ones from the configuration file) use:

```stakooler accounts details```

This will show balance, rewards, staked and unbonding tokens for each account


> NOTE: Still Work In Progress (WIP)! First release coming soon!
