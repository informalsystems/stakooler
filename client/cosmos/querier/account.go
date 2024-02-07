package querier

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/informalsystems/stakooler/client/cosmos/api"
	"github.com/informalsystems/stakooler/client/cosmos/api/osmosis"
	"github.com/informalsystems/stakooler/client/cosmos/model"
)

const zeroAmount = 0.00000

type TokenDetail struct {
	Symbol    string
	Precision int
}

func LoadAuthData(account *model.Account, client *http.Client) error {
	var authResponse model.AuthResponse
	authResponse, err := api.GetAuth(account, client)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get auth info: %s", err))
	}
	for _, value := range authResponse.Account.BaseVestingAccount.OriginalVesting {
		metadata := GetDenomMetadata(value.Denom, *account, client)
		amount, err2 := strconv.ParseFloat(value.Amount, 1)
		if err2 != nil {
			return errors.New(fmt.Sprintf("error converting rewards amount: %s", err2))
		} else {
			if amount > zeroAmount {
				convertedAmount := amount / math.Pow10(metadata.Precision)
				foundToken := false
				for j := range account.TokensEntry {
					if strings.ToLower(account.TokensEntry[j].Denom) == strings.ToLower(value.Denom) {
						account.TokensEntry[j].Vesting += convertedAmount
						foundToken = true
					}
				}
				// If there were no tokens of this denom yet, create one
				if !foundToken {
					account.TokensEntry = append(account.TokensEntry, model.TokenEntry{
						DisplayName: metadata.Symbol,
						Denom:       value.Denom,
						Vesting:     convertedAmount,
					})
				}
			}
		}
	}

	for _, value := range authResponse.Account.BaseVestingAccount.DelegatedVesting {
		metadata := GetDenomMetadata(value.Denom, *account, client)
		amount, err2 := strconv.ParseFloat(value.Amount, 1)
		if err2 != nil {
			return errors.New(fmt.Sprintf("error converting rewards amount: %s", err2))
		} else {
			if amount > zeroAmount {
				convertedAmount := amount / math.Pow10(metadata.Precision)
				foundToken := false
				for j := range account.TokensEntry {
					if strings.ToLower(account.TokensEntry[j].Denom) == strings.ToLower(value.Denom) {
						account.TokensEntry[j].DelegatedVesting += convertedAmount
						foundToken = true
					}
				}
				// If there were no tokens of this denom yet, create one
				if !foundToken {
					account.TokensEntry = append(account.TokensEntry, model.TokenEntry{
						DisplayName:      metadata.Symbol,
						Denom:            value.Denom,
						DelegatedVesting: convertedAmount,
					})
				}
			}
		}
	}

	return nil
}

func LoadBankBalances(account *model.Account, client *http.Client) error {
	balancesResponse, err := api.GetBalances(account, client)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get balances: %s", err))
	}

	for i := range balancesResponse.Balances {
		balance := balancesResponse.Balances[i]
		// Skip liquidity pools and IBC tokens
		if strings.HasPrefix(strings.ToUpper(balance.Denom), "GAMM/POOL/") ||
			strings.HasPrefix(strings.ToUpper(balance.Denom), "IBC/") {
			continue
		}

		metadata := GetDenomMetadata(balance.Denom, *account, client)
		amount, err2 := strconv.ParseFloat(balance.Amount, 1)
		if err2 != nil {
			return errors.New(fmt.Sprintf("error converting balance amount: %s", err2))
		} else {
			var convertedAmount float64
			if amount > zeroAmount {
				convertedAmount = amount / math.Pow10(metadata.Precision)
			} else {
				convertedAmount = zeroAmount
			}
			foundToken := false
			for j := range account.TokensEntry {
				if strings.ToLower(account.TokensEntry[j].Denom) == strings.ToLower(balance.Denom) {
					account.TokensEntry[j].Balance += convertedAmount
					foundToken = true
				}
			}
			// If there were no tokens of this denom yet, create one
			if !foundToken {
				account.TokensEntry = append(account.TokensEntry, model.TokenEntry{
					DisplayName: metadata.Symbol,
					Denom:       balance.Denom,
					Balance:     convertedAmount,
				})
			}
		}
	}

	return nil
}

func LoadDistributionData(account *model.Account, client *http.Client) error {
	rewardsResponse, err := api.GetRewards(account, client)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get rewards: %s", err))
	}

	for i := range rewardsResponse.Total {
		reward := rewardsResponse.Total[i]
		metadata := GetDenomMetadata(reward.Denom, *account, client)

		// Skip liquidity pools and IBC tokens
		if strings.HasPrefix(strings.ToUpper(reward.Denom), "GAMM/POOL/") ||
			strings.HasPrefix(strings.ToUpper(reward.Denom), "IBC/") {
			continue
		}

		amount, err2 := strconv.ParseFloat(reward.Amount, 1)
		if err2 != nil {
			return errors.New(fmt.Sprintf("error converting rewards amount: %s", err2))
		} else {
			if amount > zeroAmount {
				convertedAmount := amount / math.Pow10(metadata.Precision)
				foundToken := false
				for j := range account.TokensEntry {
					if strings.ToLower(account.TokensEntry[j].Denom) == strings.ToLower(reward.Denom) {
						account.TokensEntry[j].Reward += convertedAmount
						foundToken = true
					}
				}
				// If there were no tokens of this denom yet, create one
				if !foundToken {
					account.TokensEntry = append(account.TokensEntry, model.TokenEntry{
						DisplayName: metadata.Symbol,
						Denom:       reward.Denom,
						Reward:      convertedAmount,
					})
				}
			}
		}
	}

	validator, err := GetValidatorAccount(account)
	if err != nil {
		return errors.New("cannot retrieve validator account")
	} else {
		commissions, err2 := api.GetCommissions(account, validator, client)
		if err2 != nil {
			return errors.New(fmt.Sprintf("Failed to get commissions: %s", err2))
		} else {
			for i := range commissions.Commissions.Commission {
				commission := commissions.Commissions.Commission[i]
				metadata := GetDenomMetadata(commission.Denom, *account, client)
				amount, err3 := strconv.ParseFloat(commission.Amount, 1)
				if err3 != nil {
					return errors.New(fmt.Sprintf("error converting commission amount: %s", err3))
				} else {
					if amount > zeroAmount {
						convertedAmount := amount / math.Pow10(metadata.Precision)
						foundToken := false
						for j := range account.TokensEntry {
							if strings.ToLower(account.TokensEntry[j].Denom) == strings.ToLower(commission.Denom) {
								account.TokensEntry[j].Commission += convertedAmount
								foundToken = true
							}
						}
						// If there were no tokens of this denom yet, create one
						if !foundToken {
							account.TokensEntry = append(account.TokensEntry, model.TokenEntry{
								DisplayName: metadata.Symbol,
								Denom:       commission.Denom,
								Commission:  convertedAmount,
							})
						}
					}
				}
			}
		}
	}

	return nil
}

func LoadStakingData(account *model.Account, client *http.Client) error {
	delegations, err := api.GetDelegations(account, client)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get delegations: %s", err))
	}

	params, err := api.GetStakingParams(account.Chain.LCD, client)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get staking params: %s", err))
	}

	metadata := GetDenomMetadata(params.ParamsResponse.BondDenom, *account, client)

	for i := range delegations.DelegationResponses {
		delegation := delegations.DelegationResponses[i]
		amount, err2 := strconv.ParseFloat(delegation.Balance.Amount, 1)
		if err2 != nil {
			return errors.New(fmt.Sprintf("error converting delegation amount: %s", err2))
		} else {
			if amount > zeroAmount {
				convertedAmount := amount / math.Pow10(metadata.Precision)
				foundToken := false
				for j := range account.TokensEntry {
					if strings.ToLower(account.TokensEntry[j].Denom) == strings.ToLower(params.ParamsResponse.BondDenom) {
						account.TokensEntry[j].Delegation += convertedAmount
						foundToken = true
					}
				}

				if !foundToken {
					account.TokensEntry = append(account.TokensEntry, model.TokenEntry{
						DisplayName: metadata.Symbol,
						Denom:       params.ParamsResponse.BondDenom,
						Delegation:  convertedAmount,
					})
				}
			}
		}
	}

	unbondings, err := api.GetUnbondings(account, client)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get unbondings: %s", err))
	}

	for i := range unbondings.UnbondingResponses {
		unbonding := unbondings.UnbondingResponses[i]
		for j := range unbonding.Entries {
			amount, err2 := strconv.ParseFloat(unbonding.Entries[j].Balance, 1)
			if err2 != nil {
				return errors.New(fmt.Sprintf("error converting unbonding amount: %s", err2))
			} else {
				if amount > zeroAmount {
					convertedAmount := amount / math.Pow10(metadata.Precision)
					foundToken := false
					for k := range account.TokensEntry {
						if strings.ToLower(account.TokensEntry[k].Denom) == strings.ToLower(params.ParamsResponse.BondDenom) {
							account.TokensEntry[k].Unbonding += convertedAmount
							foundToken = true
						}
					}

					if !foundToken {
						account.TokensEntry = append(account.TokensEntry, model.TokenEntry{
							DisplayName: metadata.Symbol,
							Denom:       params.ParamsResponse.BondDenom,
							Unbonding:   convertedAmount,
						})
					}
				}
			}
		}
	}

	return nil
}

// GetDenomMetadata This function checks if the denom is for a chain (e.g. Osmosis or Sifchain)
// that keeps an asset list or registry for their denominations for the IBC denoms
// or the liquidity pools. The function returns the UI friendly name and the exponent
// used by the denom. If there are any errors just return the denom and 0 for
// the precision exponent
func GetDenomMetadata(denom string, account model.Account, client *http.Client) TokenDetail {
	symbol := denom
	precision := 0
	bech32Prefix, _, _ := bech32.DecodeAndConvert(account.Address)

	// Check the chain and bech32Prefix and if matches one of the chains that have a registry or asset list
	// use that information to find the token metadata
	// TODO: In the future use the information from the chain registry instead of hard-coded values
	if strings.ToLower(bech32Prefix) == "osmo" {
		// TODO: Don't fetch this for every account
		list, _ := osmosis.GetAssetsList()
		symbol, precision = list.GetSymbolExponent(denom)
	} else {
		// Try to get the denometadata from the chain
		//if strings.Contains(denom, "ibc/") {
		//	denomM
		//}
		denomMetadata, _ := api.GetDenomMetadata(&account, denom, client)
		// In case no denom metadata is available just use the denom - 'u' and precision 6
		if denomMetadata.Metadata.Base == "" {
			symbol = strings.ToUpper(denom[1:])
			precision = 6
		} else {
			symbol = denomMetadata.Metadata.Display
			precision = api.GetExponent(&denomMetadata)
		}

	}

	// If it's an IBC denom add '(IBC)' to the symbol
	if strings.HasPrefix(strings.ToUpper(denom), "IBC/") {
		symbol = symbol + " (ibc)"
	}
	return TokenDetail{strings.ToUpper(symbol), precision}
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
