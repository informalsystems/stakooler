package model

type ZabbixConfig struct {
	Host string
	Port string
}

type Config struct {
	Accounts   Accounts
	Validators Validators
	Chains     Chains
	Zabbix     ZabbixConfig
}
