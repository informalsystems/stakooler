package querier

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/informalsystems/stakooler/client/cosmos/api"
	"github.com/informalsystems/stakooler/client/cosmos/model"
)

const zeroAmount = 0.00000

type TokenDetail struct {
	Symbol    string
	Precision int
}

func LoadAuthData(account *model.Account, client *http.Client, chain *model.Chain) error {
	acctResponse, err := api.GetAaccount(account.Address, chain.RestEndpoint, client)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get auth info: %s", err))
	}

	for _, value := range acctResponse.Account.BaseVestingAccount.OriginalVesting {
		metadata := GetDenomMetadata(value.Denom, chain, client)
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

	for _, value := range acctResponse.Account.BaseVestingAccount.DelegatedVesting {
		metadata := GetDenomMetadata(value.Denom, chain, client)
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

func LoadBankBalances(account *model.Account, client *http.Client, chain *model.Chain) error {
	balancesResponse, err := api.GetBalances(account.Address, chain.RestEndpoint, client)
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

		metadata := GetDenomMetadata(balance.Denom, chain, client)
		amount, err2 := strconv.ParseFloat(balance.Amount, 1)
		if err2 != nil {
			return errors.New(fmt.Sprintf("error converting balance amount: %s", err2))
		} else {
			if amount > zeroAmount {
				convertedAmount := amount / math.Pow10(metadata.Precision)
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
	}
	return nil
}

func LoadDistributionData(account *model.Account, client *http.Client, chain *model.Chain) error {
	rewardsResponse, err := api.GetRewards(account.Address, chain.RestEndpoint, client)
	if err != nil {
		return errors.New(fmt.Sprintf("querying %s for rewards failed: %s", chain.Id, err))
	}

	for i := range rewardsResponse.Total {
		reward := rewardsResponse.Total[i]
		metadata := GetDenomMetadata(reward.Denom, chain, client)

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

	commissions, err := api.GetCommissions(account.Valoper, chain.RestEndpoint, client)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to get commissions: %s", err))
	} else {
		for i := range commissions.Commissions.Commission {
			commission := commissions.Commissions.Commission[i]
			metadata := GetDenomMetadata(commission.Denom, chain, client)

			// Skip liquidity pools and IBC tokens
			if strings.HasPrefix(strings.ToUpper(commission.Denom), "GAMM/POOL/") ||
				strings.HasPrefix(strings.ToUpper(commission.Denom), "IBC/") {
				continue
			}

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
	return nil
}

/*
func LoadStakingData(account *model.Account, client *http.Client) error {
	delegations, err := api.GetDelegations(account, client)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get delegations: %s", err))
	}

	params, err := api.GetStakingParams(account.Chain.RestEndpoint, client)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get staking params: %s", err))
	}

	metadata := GetDenomMetadataFromBank(params.ParamsResponse.BondDenom, *account, client)

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
*/

// GetDenomMetadata This function checks if the denom is for a chain (e.g. Osmosis or Sifchain)
// that keeps an asset list or registry for their denominations for the IBC denoms
// or the liquidity pools. The function returns the UI friendly name and the exponent
// used by the denom. If there are any errors just return the denom and 0 for
// the precision exponent
func GetDenomMetadata(denom string, chain *model.Chain, client *http.Client) TokenDetail {
	var symbol string
	exponent := 0

	if chain.AssetList != nil {
		symbol, exponent = chain.AssetList.GetAssetDetails(denom)
	}
	// in case asset details are missing
	if exponent == 0 {
		denomMetadata, _ := api.GetDenomMetadataFromBank(denom, chain.RestEndpoint, client)
		if denomMetadata.Metadata.Base == "" {
			symbol = denom
			exponent = 6
		} else {
			symbol = strings.ToUpper(denomMetadata.Metadata.Display)
			exponent = api.GetExponent(&denomMetadata)
		}
	}

	// If it's an IBC denom add '(IBC)' to the symbol
	if strings.HasPrefix(strings.ToUpper(denom), "IBC/") {
		symbol = symbol + " (IBC)"
	}
	return TokenDetail{symbol, exponent}
}
