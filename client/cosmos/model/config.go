package model

type ZabbixConfig struct {
	Server string
	Port   int
	Host   string
}

type Config struct {
	Accounts   Accounts
	Validators ValidatorList
	Chains     Chains
	Zabbix     ZabbixConfig
}
