package model

type ZabbixConfig struct {
	Server string
	Port   int
	Host   string
}

type Config struct {
	Accounts   Accounts
	Validators Validators
	Chains     Chains
	Zabbix     ZabbixConfig
}
