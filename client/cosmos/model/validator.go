package model

import "time"

type Validators struct {
	Entries []*Validator
}

type Validator struct {
	Name       string
	Address    string
	Chain      []Chain
	Statistics []Statistic
}

type Statistic struct {
	BlockTime     time.Time
	BlockHeight   int64
	VotingPower   float64
	VotingPercent float64
	Ranking       int
	NumDelegators int
	TokensStaked  float64
}
