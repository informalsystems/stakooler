package querier

import (
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/bech32"
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
		Commissions:      make(map[string]float64),
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
					if strings.HasPrefix(strings.ToUpper(balance.Denom), "IBC/") {
						accountDetails.AvailableBalance[symbol+" (IBC)"] = convertedAmount
					} else {
						accountDetails.AvailableBalance[symbol] = convertedAmount
					}

				} else {
					denomMetadata, err := api.GetDenomMetadata(account, balance.Denom)
					if err != nil {
						return accountDetails, errors.New("cannot retrieve token denom metadata")
					} else {
						// Convert amount based on exponent
						exponent := denomMetadata.GetExponent()
						convertedAmount := amount / math.Pow10(exponent)
						// Check if it's an Osmosis Pool, if so then skip it. Will add support for LP later
						if !strings.HasPrefix(strings.ToUpper(balance.Denom), "GAMM/POOL/") {
							accountDetails.AvailableBalance[strings.ToUpper(denomMetadata.Metadata.Display)] = convertedAmount
						}
					}
				}
			}
		}
	}

	for i := range rewardsResponse.Total {
		reward := rewardsResponse.Total[i]
		amount, err := strconv.ParseFloat(reward.Amount, 1)
		if err != nil {
			return accountDetails, errors.New(fmt.Sprintf("Error converting amount: %s", err))
		} else {
			if amount > 0 {
				symbol, found := assets.GetSymbol(reward.Denom)
				if found {
					exponent := assets.GetExponent(symbol)
					convertedAmount := amount / math.Pow10(exponent)
					if strings.HasPrefix(strings.ToUpper(reward.Denom), "IBC/") {
						accountDetails.Rewards[symbol+" (IBC)"] = convertedAmount
					} else {
						accountDetails.Rewards[symbol] = convertedAmount
					}
				} else {
					denomMetadata, err := api.GetDenomMetadata(account, reward.Denom)
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

	totalAmount := 0.0
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
					totalAmount += convertedAmount
					if strings.HasPrefix(strings.ToUpper(delegation.Balance.Denom), "IBC/") {
						accountDetails.Delegations[symbol+" (IBC)"] = totalAmount
					} else {
						accountDetails.Delegations[symbol] = totalAmount
					}
				} else {
					denomMetadata, err := api.GetDenomMetadata(account, delegation.Balance.Denom)
					if err != nil {
						return accountDetails, errors.New("cannot retrieve token denom metadata")
					} else {
						//TODO: Convert amount based on exponent
						exponent := denomMetadata.GetExponent()
						convertedAmount := amount / math.Pow10(exponent)
						totalAmount += convertedAmount
						accountDetails.Delegations[strings.ToUpper(denomMetadata.Metadata.Display)] = totalAmount
					}
				}
			}
		}
	}

	totalAmount = 0.0
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
						totalAmount += convertedAmount
						if strings.HasPrefix(strings.ToUpper(params.ParamsResponse.BondDenom), "IBC/") {
							accountDetails.Unbondings[symbol+" (IBC)"] = totalAmount
						} else {
							accountDetails.Unbondings[symbol] = totalAmount
						}
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
							totalAmount += convertedAmount
							accountDetails.Unbondings[strings.ToUpper(denomMetadata.Metadata.Display)] = totalAmount
						}
					}
				}
			}
		}
	}

	// Get commissions
	totalAmount = 0.0
	validator, err := GetValidatorAccount(account)
	if err != nil {
		return accountDetails, errors.New("cannot retrieve validator account")
	} else {
		commissions, err := api.GetCommissions(account, validator)
		if err != nil {
			return accountDetails, errors.New(fmt.Sprintf("Failed to get commissions: %s", err))
		} else {
			for i := range commissions.Commissions.Commission {
				commission := commissions.Commissions.Commission[i]
				amount, err := strconv.ParseFloat(commission.Amount, 1)
				if err != nil {
					return accountDetails, errors.New(fmt.Sprintf("Error converting amount: %s", err))
				} else {
					if amount > 0 {
						symbol, found := assets.GetSymbol(commission.Denom)
						if found {
							exponent := assets.GetExponent(symbol)
							convertedAmount := amount / math.Pow10(exponent)
							totalAmount += convertedAmount
							if strings.HasPrefix(strings.ToUpper(commission.Denom), "IBC/") {
								accountDetails.Commissions[symbol+" (IBC)"] = totalAmount
							} else {
								accountDetails.Commissions[symbol] = totalAmount
							}
						} else {
							denomMetadata, err := api.GetDenomMetadata(account, commission.Denom)
							if err != nil {
								return accountDetails, errors.New("cannot retrieve token denom metadata")
							} else {
								// Convert amount based on exponent
								exponent := denomMetadata.GetExponent()
								convertedAmount := amount / math.Pow10(exponent)
								accountDetails.Commissions[strings.ToUpper(denomMetadata.Metadata.Display)] = convertedAmount
							}
						}
					}
				}
			}
		}
	}

	return accountDetails, nil
}

func GetValidatorAccount(account *model.Account) (string, error) {
	acct, acctBytes, err := bech32.DecodeAndConvert(account.Address)
	if err != nil {
		fmt.Println("Error decoding", account.Address, " Error:", err)
		return "", err
	}

	validatorAccount, err := bech32.ConvertAndEncode(acct+"valoper", acctBytes)
	if err != nil {
		fmt.Println("Error converting and encoding", account.Address, " Error:", err)
		return "", err
	}
	return validatorAccount, nil
}
