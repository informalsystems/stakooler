package model

import "time"

type Delegations struct {
	DelegationResponses []struct {
		Delegation struct {
			DelegatorAddress string `json:"delegator_address"`
			ValidatorAddress string `json:"validator_address"`
			Shares           string `json:"shares"`
		} `json:"delegation"`
		Balance struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"balance"`
	} `json:"delegation_responses"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}

type Unbondings struct {
	UnbondingResponses []struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		Entries          []struct {
			CreationHeight string    `json:"creation_height"`
			CompletionTime time.Time `json:"completion_time"`
			InitialBalance string    `json:"initial_balance"`
			Balance        string    `json:"balance"`
		} `json:"entries"`
	} `json:"unbonding_responses"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}

type Params struct {
	ParamsResponse struct {
		UnbondingTime     string `json:"unbonding_time"`
		MaxValidators     int    `json:"max_validators"`
		MaxEntries        int    `json:"max_entries"`
		HistoricalEntries int    `json:"historical_entries"`
		BondDenom         string `json:"bond_denom"`
		MinCommissionRate string `json:"min_commission_rate"`
	} `json:"params"`
}

type ValidatorSet struct {
	BlockHeight string `json:"block_height"`
	Validators  []struct {
		Address string `json:"address"`
		PubKey  struct {
			Type string `json:"@type"`
			Key  string `json:"key"`
		} `json:"pub_key"`
		VotingPower      string `json:"voting_power"`
		ProposerPriority string `json:"proposer_priority"`
	} `json:"validators"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}

type Validators struct {
	BlockHeight        string `json:"block_height,omitempty"`
	ValidatorsResponse []struct {
		OperatorAddress string `json:"operator_address"`
		ConsensusPubkey struct {
			Type string `json:"@type"`
			Key  string `json:"key"`
		} `json:"consensus_pubkey"`
		Jailed          bool   `json:"jailed"`
		Status          string `json:"status"`
		Tokens          string `json:"tokens"`
		DelegatorShares string `json:"delegator_shares"`
		Description     struct {
			Moniker         string `json:"moniker"`
			Identity        string `json:"identity"`
			Website         string `json:"website"`
			SecurityContact string `json:"security_contact"`
			Details         string `json:"details"`
		} `json:"description"`
		UnbondingHeight string    `json:"unbonding_height"`
		UnbondingTime   time.Time `json:"unbonding_time"`
		Commission      struct {
			CommissionRates struct {
				Rate          string `json:"rate"`
				MaxRate       string `json:"max_rate"`
				MaxChangeRate string `json:"max_change_rate"`
			} `json:"commission_rates"`
			UpdateTime time.Time `json:"update_time"`
		} `json:"commission"`
		MinSelfDelegation string `json:"min_self_delegation"`
	} `json:"validators"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}
