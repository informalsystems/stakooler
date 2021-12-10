package querier

import (
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/informalsystems/stakooler/client/cosmos/api"
	"github.com/informalsystems/stakooler/client/cosmos/api/osmosis"
	"github.com/informalsystems/stakooler/client/cosmos/api/sifchain"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"math"
	"strconv"
	"strings"
)

const zeroAmount = 0.00000

type TokenDetail struct {
	Symbol    string
	Precision int
}

func LoadAccountDetails(account *model.Account) (model.AccountDetails, error) {
	accountDetails := model.AccountDetails{
		AvailableBalance: make(map[string]float64),
		Rewards:          make(map[string]float64),
		Delegations:      make(map[string]float64),
		Unbondings:       make(map[string]float64),
		Commissions:      make(map[string]float64),
	}

	// Get Balances
	balancesResponse, err := api.GetBalances(account)
	if err != nil {
		return accountDetails, errors.New(fmt.Sprintf("failed to get balances: %s", err))
	}

	//totalAmount := 0.0
	for i := range balancesResponse.Balances {
		balance := balancesResponse.Balances[i]
		token := GetTokenDetails(balance.Denom, *account)
		amount, err := strconv.ParseFloat(balance.Amount, 1)
		if err != nil {
			return accountDetails, errors.New(fmt.Sprintf("error converting balance amount: %s", err))
		} else {
			if amount > zeroAmount {
				// Skip liquidity pools
				if !strings.HasPrefix(strings.ToUpper(token.Symbol), "GAMM/POOL/") {
					convertedAmount := amount / math.Pow10(token.Precision)
					//totalAmount += convertedAmount
					accountDetails.AvailableBalance[token.Symbol] = convertedAmount
				}
			}
		}
	}

	// Get Rewards
	rewardsResponse, err := api.GetRewards(account)
	if err != nil {
		return accountDetails, errors.New(fmt.Sprintf("failed to get rewards: %s", err))
	}

	totalAmount := 0.0
	for i := range rewardsResponse.Total {
		reward := rewardsResponse.Total[i]
		token := GetTokenDetails(reward.Denom, *account)
		amount, err := strconv.ParseFloat(reward.Amount, 1)
		if err != nil {
			return accountDetails, errors.New(fmt.Sprintf("error converting rewards amount: %s", err))
		} else {
			if amount > zeroAmount {
				convertedAmount := amount / math.Pow10(token.Precision)
				totalAmount += convertedAmount
				accountDetails.Rewards[token.Symbol] = totalAmount
			}
		}
	}

	// Get Delegations
	delegations, err := api.GetDelegations(account)
	if err != nil {
		return accountDetails, errors.New(fmt.Sprintf("failed to get delegations: %s", err))
	}

	totalAmount = 0.0
	for i := range delegations.DelegationResponses {
		delegation := delegations.DelegationResponses[i]
		token := GetTokenDetails(delegation.Balance.Denom, *account)
		amount, err := strconv.ParseFloat(delegation.Balance.Amount, 1)
		if err != nil {
			return accountDetails, errors.New(fmt.Sprintf("error converting delegation amount: %s", err))
		} else {
			if amount > zeroAmount {
				convertedAmount := amount / math.Pow10(token.Precision)
				totalAmount += convertedAmount
				accountDetails.Delegations[token.Symbol] = totalAmount
			}
		}
	}

	//Get Unbondings
	unbondings, err := api.GetUnbondings(account)
	if err != nil {
		return accountDetails, errors.New(fmt.Sprintf("failed to get unbondings: %s", err))
	}

	totalAmount = 0.0
	for i := range unbondings.UnbondingResponses {
		unbonding := unbondings.UnbondingResponses[i]
		for i := range unbonding.Entries {
			params, err := api.GetStakingParams(account)
			if err != nil {
				return accountDetails, errors.New(fmt.Sprintf("failed to get staking params: %s", err))
			}
			token := GetTokenDetails(params.ParamsResponse.BondDenom, *account)
			amount, err := strconv.ParseFloat(unbonding.Entries[i].Balance, 1)
			if err != nil {
				return accountDetails, errors.New(fmt.Sprintf("error converting unbonding amount: %s", err))
			} else {
				if amount > zeroAmount {
					convertedAmount := amount / math.Pow10(token.Precision)
					totalAmount += convertedAmount
					accountDetails.Unbondings[token.Symbol] = totalAmount
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
				token := GetTokenDetails(commission.Denom, *account)
				amount, err := strconv.ParseFloat(commission.Amount, 1)
				if err != nil {
					return accountDetails, errors.New(fmt.Sprintf("error converting commission amount: %s", err))
				} else {
					if amount > zeroAmount {
						convertedAmount := amount / math.Pow10(token.Precision)
						totalAmount += convertedAmount
						accountDetails.Commissions[token.Symbol] = totalAmount
					}
				}
			}
		}
	}

	return accountDetails, nil
}

// This function checks if the denom is for a chain (e.g. Osmosis or Sifchain)
// that keeps an asset list or registry for their denominations for the IBC denoms
// or the liquitiy pools. The function returns the UI friendly name and the exponent
// used by the denom. If there are any errors just return the denom and 0 for
// the precision exponent
func GetTokenDetails(denom string, account model.Account) TokenDetail {
	symbol := denom
	precision := 0
	bech32Prefix, _, _ := bech32.DecodeAndConvert(account.Address)

	// Check the chain and bech32Prefix and if matches one of the chains that have a registry or asset list
	// use that information to find the token metadata
	// TODO: In the future use the information from the chain registry instead of hard-coded values
	if strings.ToLower(account.Chain.ID) == "osmosis-1" && strings.ToLower(bech32Prefix) == "osmo" {
		// TODO: Don't fetch this for every account
		list, _ := osmosis.GetAssetsList()
		symbol, precision = list.GetSymbolExponent(denom)
	} else if strings.ToLower(account.Chain.ID) == "sifchain-1" && strings.ToLower(bech32Prefix) == "sif" {
		// TODO: Don't fetch this for every account
		tokenList, _ := sifchain.GetTokenList()
		symbol, precision = tokenList.GetSymbolExponent(denom)
	} else {
		// Try to get the denometadata from the chain
		denomMetadata, _ := api.GetDenomMetadata(&account, denom)
		symbol = denomMetadata.Metadata.Display
		precision = denomMetadata.GetExponent()
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
