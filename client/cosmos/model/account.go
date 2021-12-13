package model

import "time"

type Accounts struct {
	Entries []*Account
}

type Account struct {
	Name        string
	Address     string
	Chain       Chain
	TokensEntry []TokenEntry
}

type TokenEntry struct {
	DisplayName string
	Denom       string
	BlockTime   time.Time
	BlockHeight string
	Balance     float64
	Reward      float64
	Delegation  float64
	Unbonding   float64
	Commission  float64
}
