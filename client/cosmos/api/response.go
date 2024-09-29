package api

const OriginalVesting = 0
const DelegatedVesting = 1
const Bank = 3
const Rewards = 4
const Commission = 5
const Delegation = 6
const Unbonding = 7

type Response interface {
	GetBalances() map[int]map[string]string
}
