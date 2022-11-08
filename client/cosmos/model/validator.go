package model

import "time"

type Validators struct {
	Entries []*Validator
}

type Validator struct {
	Moniker        string
	ValoperAddress string
	Chain          Chain
	BlockTime      time.Time
	BlockHeight    string
	VotingPower    int64
	VotingPercent  float64
	Ranking        int
	NumValidators  string
	NumDelegators  string
	Unbondings     int64
}
