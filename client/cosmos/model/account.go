package model

import (
	"time"
)

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
	Vesting     float64
}

type AuthResponse struct {
	Account struct {
		Type               string `json:"@type"`
		BaseVestingAccount struct {
			BaseAccount struct {
				Address       string `json:"address,omitempty"`
				PubKey        string `json:"public_key,omitempty"`
				AccountNumber string `json:"account_number,omitempty"`
				Sequence      string `json:"sequence,omitempty"`
			}
			OriginalVesting []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"original_vesting"`
			DelegatedFree []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"delegated_free"`
			DelegatedVesting []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"delegated_vesting"`
			EndTime string `json:"end_time"`
		} `json:"base_vesting_account"`
		StartTime      string `json:"start_time"`
		VestingPeriods []struct {
			Length string `json:"length"`
			Amount []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"amount"`
		} `json:"vesting_periods"`
	} `json:"account"`
}
