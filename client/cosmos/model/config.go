package model

type ZabbixConfig struct {
	Host string
	Port int
}

type Config struct {
	Accounts   Accounts
	Validators Validators
	Chains     Chains
	Zabbix     ZabbixConfig
}
