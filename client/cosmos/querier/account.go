package querier

import (
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/informalsystems/stakooler/client/cosmos/api"
	"github.com/informalsystems/stakooler/client/cosmos/api/osmosis"
	"github.com/informalsystems/stakooler/client/cosmos/api/sifchain"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"github.com/schollz/progressbar/v3"
	"math"
	"strconv"
	"strings"
)

const zeroAmount = 0.00000

type TokenDetail struct {
	Symbol    string
	Precision int
}

func LoadTokenInfo(account *model.Account, bar *progressbar.ProgressBar) error {

	var tokens []model.TokenEntry

	// Get Latest Block Information
	// Use the same block information for all the entries
	blockResponse, err := api.GetLatestBlock(account.Chain)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get latest block: %s", err))
	}
	bar.Add(1)

	// Get Balances
	balancesResponse, err := api.GetBalances(account)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get balances: %s", err))
	}

	bar.Add(1)
	for i := range balancesResponse.Balances {
		// Skip liquidity pools and IBC tokens
		balance := balancesResponse.Balances[i]
		if !strings.HasPrefix(strings.ToUpper(balance.Denom), "GAMM/POOL/") &&
			!strings.HasPrefix(strings.ToUpper(balance.Denom), "IBC/") {

			metadata := GetTokenMetadata(balance.Denom, *account)
			token := model.TokenEntry{}
			token.DisplayName = metadata.Symbol
			token.Denom = balance.Denom
			token.BlockTime = blockResponse.Block.Header.Time
			token.BlockHeight = blockResponse.Block.Header.Height
			amount, err := strconv.ParseFloat(balance.Amount, 1)
			if err != nil {
				return errors.New(fmt.Sprintf("error converting balance amount: %s", err))
			} else {
				if amount > zeroAmount {
					convertedAmount := amount / math.Pow10(metadata.Precision)
					token.Balance = convertedAmount
					tokens = append(tokens, token)
				}
			}
		}
	}

	// Get Rewards
	rewardsResponse, err := api.GetRewards(account)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get rewards: %s", err))
	}

	bar.Add(1)
	totalAmount := 0.0
	for i := range rewardsResponse.Total {
		reward := rewardsResponse.Total[i]
		metadata := GetTokenMetadata(reward.Denom, *account)
		amount, err := strconv.ParseFloat(reward.Amount, 1)
		if err != nil {
			return errors.New(fmt.Sprintf("error converting rewards amount: %s", err))
		} else {
			if amount > zeroAmount {
				convertedAmount := amount / math.Pow10(metadata.Precision)
				totalAmount += convertedAmount
				for i := range tokens {
					if strings.ToLower(tokens[i].Denom) == strings.ToLower(reward.Denom) {
						tokens[i].Reward = totalAmount
					}
				}
			}
		}
	}

	// Get Delegations
	delegations, err := api.GetDelegations(account)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get delegations: %s", err))
	}

	bar.Add(1)
	totalAmount = 0.0
	for i := range delegations.DelegationResponses {
		delegation := delegations.DelegationResponses[i]
		metadata := GetTokenMetadata(delegation.Balance.Denom, *account)
		amount, err := strconv.ParseFloat(delegation.Balance.Amount, 1)
		if err != nil {
			return errors.New(fmt.Sprintf("error converting delegation amount: %s", err))
		} else {
			if amount > zeroAmount {
				convertedAmount := amount / math.Pow10(metadata.Precision)
				totalAmount += convertedAmount
				for i := range tokens {
					if strings.ToLower(tokens[i].Denom) == strings.ToLower(delegation.Balance.Denom) {
						tokens[i].Delegation = totalAmount
					}
				}
			}
		}
	}

	//Get Unbondings
	unbondings, err := api.GetUnbondings(account)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get unbondings: %s", err))
	}

	bar.Add(1)
	totalAmount = 0.0
	for i := range unbondings.UnbondingResponses {
		unbonding := unbondings.UnbondingResponses[i]
		for i := range unbonding.Entries {
			params, err := api.GetStakingParams(account.Chain.LCD)
			if err != nil {
				return errors.New(fmt.Sprintf("failed to get staking params: %s", err))
			}
			metadata := GetTokenMetadata(params.ParamsResponse.BondDenom, *account)
			amount, err := strconv.ParseFloat(unbonding.Entries[i].Balance, 1)
			if err != nil {
				return errors.New(fmt.Sprintf("error converting unbonding amount: %s", err))
			} else {
				if amount > zeroAmount {
					convertedAmount := amount / math.Pow10(metadata.Precision)
					totalAmount += convertedAmount
					for i := range tokens {
						if strings.ToLower(tokens[i].Denom) == strings.ToLower(params.ParamsResponse.BondDenom) {
							tokens[i].Unbonding = totalAmount
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
		return errors.New("cannot retrieve validator account")
	} else {
		commissions, err := api.GetCommissions(account, validator)
		bar.Add(1)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to get commissions: %s", err))
		} else {
			for i := range commissions.Commissions.Commission {
				commission := commissions.Commissions.Commission[i]
				metadata := GetTokenMetadata(commission.Denom, *account)
				amount, err := strconv.ParseFloat(commission.Amount, 1)
				if err != nil {
					return errors.New(fmt.Sprintf("error converting commission amount: %s", err))
				} else {
					if amount > zeroAmount {
						convertedAmount := amount / math.Pow10(metadata.Precision)
						totalAmount += convertedAmount
						for i := range tokens {
							if strings.ToLower(tokens[i].Denom) == strings.ToLower(commission.Denom) {
								tokens[i].Commission = totalAmount
							}
						}
					}
				}
			}
		}
	}
	account.TokensEntry = tokens
	return nil
}

// GetTokenMetadata This function checks if the denom is for a chain (e.g. Osmosis or Sifchain)
// that keeps an asset list or registry for their denominations for the IBC denoms
// or the liquidity pools. The function returns the UI friendly name and the exponent
// used by the denom. If there are any errors just return the denom and 0 for
// the precision exponent
func GetTokenMetadata(denom string, account model.Account) TokenDetail {
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
	} else if strings.ToLower(bech32Prefix) == "sif" {
		// TODO: Don't fetch this for every account
		tokenList, _ := sifchain.GetTokenList()
		symbol, precision = tokenList.GetSymbolExponent(denom)
	} else {
		// Try to get the denometadata from the chain
		//if strings.Contains(denom, "ibc/") {
		//	denomM
		//}
		denomMetadata, _ := api.GetDenomMetadata(&account, denom)
		// In case no denom metadata is available just use the denom - 'u' and precision 6
		if denomMetadata.Metadata.Base == "" {
			symbol = strings.ToUpper(denom[1:])
			precision = 6
		} else {
			symbol = denomMetadata.Metadata.Display
			precision = denomMetadata.GetExponent()
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
