package model

import "time"

type Accounts struct {
	Entries []*Account
}

type Account struct {
	Name        string
	Address     string
	Chain       Chain
	BlockTime   time.Time
	BlockHeight string
	TokensEntry []TokenEntry
}

type TokenEntry struct {
	DisplayName string
	Denom       string
	Balance     float64
	Reward      float64
	Delegation  float64
	Unbonding   float64
	Commission  float64
}
