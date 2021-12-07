package model

type Accounts struct {
	Entries	[]*Account
}

type Account struct {
	Name           string
	Address        string
	Chain	       Chain
	AccountDetails AccountDetails
}

type AccountDetails struct {
	AvailableBalance 	map[string]float64
	Rewards          	map[string]float64
	Delegations     	map[string]float64
	Unbondings          map[string]float64
	Commissions         map[string]float64
}