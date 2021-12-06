package querier

import (
	"errors"
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/api"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"github.com/informalsystems/stakooler/client/osmosis"
	"math"
	"strconv"
	"strings"
)

func LoadAccountDetails(account *model.Account) (model.AccountDetails, error) {
	accountDetails := model.AccountDetails{
		AvailableBalance: make(map[string]float64),
		Rewards:          make(map[string]float64),
		Delegations:      make(map[string]float64),
		Unbondings:       make(map[string]float64),
	}

	// Get available balancesResponse
	balancesResponse, err := api.GetBalances(account)
	if err != nil {
		return accountDetails, errors.New(fmt.Sprintf("Failed to get balances: %s", err))
	}

	rewardsResponse, err := api.GetRewards(account)
	if err != nil {
		return accountDetails, errors.New(fmt.Sprintf("Failed to get rewards: %s", err))
	}

	delegations, err := api.GetDelegations(account)
	if err != nil {
		return accountDetails, errors.New(fmt.Sprintf("Failed to get delegations: %s", err))
	}

	unbondings, err := api.GetUnbondings(account)
	if err != nil {
		return accountDetails, errors.New(fmt.Sprintf("Failed to get unbondings: %s", err))
	}

	assets, err := osmosis.GetAssetsList()
	if err != nil {
		return accountDetails, errors.New(fmt.Sprintf("Error fetching assets list: %s", err))
	}

	for i := range balancesResponse.Balances {
		balance := balancesResponse.Balances[i]
		amount, err := strconv.ParseFloat(balance.Amount, 1)
		if err != nil {
			return accountDetails, errors.New(fmt.Sprintf("Error converting amount: %s", err))
		} else {
			if amount > 0 {
				symbol, found := assets.GetSymbol(balance.Denom)
				if found {
					exponent := assets.GetExponent(symbol)
					convertedAmount := amount / math.Pow10(exponent)
					accountDetails.AvailableBalance[symbol] = convertedAmount
				} else {
					denomMetadata, err := api.GetDenomMetadata(account, balance.Denom)
					if err != nil {
						return accountDetails, errors.New("cannot retrieve token denom metadata")
					} else {
						// Convert amount based on exponent
						exponent := denomMetadata.GetExponent()
						convertedAmount := amount / math.Pow10(exponent)
						accountDetails.AvailableBalance[strings.ToUpper(denomMetadata.Metadata.Display)] = convertedAmount
					}
				}
			}
		}
	}

	for i := range rewardsResponse.Rewards {
		reward := rewardsResponse.Rewards[i]
		for i := range reward.Reward {
			amount, err := strconv.ParseFloat(reward.Reward[i].Amount, 1)
			if err != nil {
				return accountDetails, errors.New(fmt.Sprintf("Error converting amount: %s", err))
			} else {
				if amount > 0 {
					symbol, found := assets.GetSymbol(reward.Reward[i].Denom)
					if found {
						exponent := assets.GetExponent(symbol)
						convertedAmount := amount / math.Pow10(exponent)
						accountDetails.Rewards[symbol] = convertedAmount
					} else {
						denomMetadata, err := api.GetDenomMetadata(account, reward.Reward[i].Denom)
						if err != nil {
							return accountDetails, errors.New("cannot retrieve token denom metadata")
						} else {
							// Convert amount based on exponent
							exponent := denomMetadata.GetExponent()
							convertedAmount := amount / math.Pow10(exponent)
							accountDetails.Rewards[strings.ToUpper(denomMetadata.Metadata.Display)] = convertedAmount
						}
					}
				}
			}
		}
	}

	for i := range delegations.DelegationResponses {
		delegation := delegations.DelegationResponses[i]
		amount, err := strconv.ParseFloat(delegation.Balance.Amount, 1)
		if err != nil {
			return accountDetails, errors.New(fmt.Sprintf("Error converting amount: %s", err))
		} else {
			if amount > 0 {
				symbol, found := assets.GetSymbol(delegation.Balance.Denom)
				if found {
					exponent := assets.GetExponent(symbol)
					convertedAmount := amount / math.Pow10(exponent)
					accountDetails.Delegations[symbol] = convertedAmount
				} else {
					denomMetadata, err := api.GetDenomMetadata(account, delegation.Balance.Denom)
					if err != nil {
						return accountDetails, errors.New("cannot retrieve token denom metadata")
					} else {
						//TODO: Convert amount based on exponent
						exponent := denomMetadata.GetExponent()
						convertedAmount := amount / math.Pow10(exponent)
						accountDetails.Delegations[strings.ToUpper(denomMetadata.Metadata.Display)] = convertedAmount
					}
				}
			}
		}
	}

	for i := range unbondings.UnbondingResponses {
		unbonding := unbondings.UnbondingResponses[i]
		for i := range unbonding.Entries {
			amount, err := strconv.ParseFloat(unbonding.Entries[i].Balance, 1)
			if err != nil {
				return accountDetails, errors.New(fmt.Sprintf("Error converting amount: %s", err))
			} else {
				if amount > 0 {
					params, err := api.GetStakingParams(account)
					if err != nil {
						return accountDetails, errors.New(fmt.Sprintf("Failed to get staking params: %s", err))
					}
					symbol, found := assets.GetSymbol(params.ParamsResponse.BondDenom)
					if found {
						exponent := assets.GetExponent(symbol)
						convertedAmount := amount / math.Pow10(exponent)
						accountDetails.Unbondings[symbol] = convertedAmount
					} else {
						// Use the mint params to get the denom since the unbonding response doesn't return that
						mintParams, err := api.GetMintParams(account)
						if err != nil {
							return accountDetails, errors.New("cannot retrieve mint params")
						}
						denomMetadata, err := api.GetDenomMetadata(account, mintParams.Params.MintDenom)
						if err != nil {
							return accountDetails, errors.New("cannot retrieve token denom metadata")
						} else {
							// Convert amount based on exponent
							exponent := denomMetadata.GetExponent()
							convertedAmount := amount / math.Pow10(exponent)
							accountDetails.Unbondings[strings.ToUpper(denomMetadata.Metadata.Display)] = convertedAmount
						}
					}
				}
			}
		}
	}
	return accountDetails, nil
}

