package model

import (
	"time"
)

type ValidatorList struct {
	Entries []*Validator
}

type Validator struct {
	Moniker        string
	ValoperAddress string
	BlockTime      time.Time
	BlockHeight    string
	VotingPower    int64
	VotingPercent  float64
	Ranking        int
	NumValidators  string
	NumDelegators  string
	Unbondings     int64
	Commission     float64
}
