package model

import "time"

type Validators struct {
	Entries []*Validator
}

type Validator struct {
	Name           string
	ValoperAddress string
	Chain          Chain
	BlockTime      time.Time
	BlockHeight    string
	VotingPower    int64
	VotingPercent  string
	Ranking        int
	NumDelegators  string
	Unbondings     int64
}
