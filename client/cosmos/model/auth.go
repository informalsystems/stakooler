package model

import "time"

type Account struct {
	Name        string
	Address     string
	Valoper     string
	BlockTime   time.Time
	BlockHeight string
	Tokens      map[string]*Token
}

type Token struct {
	DisplayName string
	Denom       string
	Balances    struct {
		Bank             float64
		Rewards          float64
		Commission       float64
		Delegated        float64
		Unbonding        float64
		OriginalVesting  float64
		DelegatedVesting float64
	}
}
