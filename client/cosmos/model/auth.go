package model

import "time"

type Bech32PrefixResponse struct {
	Bech32Prefix string `json:"bech32_prefix"`
}

type Accounts struct {
	Entries []*Account
}

type Account struct {
	Name        string
	Address     string
	Valoper     string
	BlockTime   time.Time
	BlockHeight string
	Tokens      map[string]*Token
}

type Token struct {
	DisplayName      string
	Denom            string
	BankBalance      float64
	Rewards          float64
	Delegation       float64
	Unbonding        float64
	Commission       float64
	OriginalVesting  float64
	DelegatedVesting float64
}
