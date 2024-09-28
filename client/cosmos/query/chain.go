package query

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

type Chains struct {
	Entries []*Chain
}

type Chain struct {
	Name         string
	Id           string
	RestEndpoint string
	Bech32Prefix string
	Accounts     []*model.Account
	BondDenom    string
	Exponent     int
	AssetList    *model.AssetList
}

func (c *Chain) FetchAccountBalances(blockInfo model.BlockResponse, client *http.Client) error {
	for idx := range c.Accounts {
		c.Accounts[idx].BlockTime = blockInfo.Block.Header.Time
		c.Accounts[idx].BlockHeight = blockInfo.Block.Header.Height

		if acctResponse, err := api.GetAccount(c.Accounts[idx].Address, c.RestEndpoint, client); err != nil {
			return errors.New(fmt.Sprintf("query account: %s", err))
		} else {
			if err = c.ProcessOriginalVesting(acctResponse, idx, client); err != nil {
				return errors.New(fmt.Sprintf("process original vesting: %s", err))
			}

			if err = c.ProcessDelegatedVesting(acctResponse, idx, client); err != nil {
				return errors.New(fmt.Sprintf("process delegated vesting: %s", err))
			}
		}

		if bankResponse, err := api.GetBalances(c.Accounts[idx].Address, c.RestEndpoint, client); err != nil {
			return errors.New(fmt.Sprintf("query bank balances: %s", err))
		} else {
			if err = c.ProcessBankBalances(bankResponse, idx, client); err != nil {
				return errors.New(fmt.Sprintf("process bank balances: %s", err))
			}
		}

		if rewardsResponse, err := api.GetRewards(c.Accounts[idx].Address, c.RestEndpoint, client); err != nil {
			return errors.New(fmt.Sprintf("query rewards: %s", err))
		} else {
			if err = c.ProcessRewards(rewardsResponse, idx, client); err != nil {
				return errors.New(fmt.Sprintf("process rewards: %s", err))
			}
		}

		if commissionResponse, err := api.GetCommissions(c.Accounts[idx].Valoper, c.RestEndpoint, client); err != nil {
			return errors.New(fmt.Sprintf("query commissions: %s", err))
		} else if commissionResponse != nil {
			if err = c.ProcessCommission(commissionResponse, idx, client); err != nil {
				return errors.New(fmt.Sprintf("process commissions: %s", err))
			}
		}

		if delegationsResponse, err := api.GetDelegations(c.Accounts[idx].Address, c.RestEndpoint, client); err != nil {
			return errors.New(fmt.Sprintf("query delegations: %s", err))
		} else {
			if err = c.ProcessDelegations(delegationsResponse, idx, client); err != nil {
				return errors.New(fmt.Sprintf("process delegations: %s", err))
			}
		}

		if unbondingResponse, err := api.GetUnbondings(c.Accounts[idx].Address, c.RestEndpoint, client); err != nil {
			return errors.New(fmt.Sprintf("query unbondings: %s", err))
		} else {
			if err = c.ProcessUnbondings(unbondingResponse, idx, client); err != nil {
				return errors.New(fmt.Sprintf("process unbondings: %s", err))
			}
		}
	}
	return nil
}

func (c *Chain) ProcessOriginalVesting(resp *model.AcctResponse, idx int, client *http.Client) error {
	for _, value := range resp.Account.BaseVestingAccount.OriginalVesting {
		if strings.HasPrefix(strings.ToUpper(value.Denom), "GAMM/POOL/") ||
			strings.HasPrefix(strings.ToUpper(value.Denom), "IBC/") {
			continue
		}
		metadata := GetDenomMetadata(value.Denom, c, client)
		amount, err := strconv.ParseFloat(value.Amount, 1)
		if err != nil {
			return err
		}

		if amount > zeroAmount {
			convertedAmount := amount / math.Pow10(metadata.Precision)
			foundToken := false
			for j := range c.Accounts[idx].Tokens {
				if strings.ToLower(c.Accounts[idx].Tokens[j].Denom) == strings.ToLower(value.Denom) {
					c.Accounts[idx].Tokens[j].Vesting += convertedAmount
					foundToken = true
				}
			}

			if !foundToken {
				c.Accounts[idx].Tokens = append(c.Accounts[idx].Tokens, &model.Token{
					DisplayName: metadata.Symbol,
					Denom:       value.Denom,
					Vesting:     convertedAmount,
				})
			}
		}
	}
	return nil
}

func (c *Chain) ProcessDelegatedVesting(resp *model.AcctResponse, idx int, client *http.Client) error {
	for _, value := range resp.Account.BaseVestingAccount.DelegatedVesting {
		if strings.HasPrefix(strings.ToUpper(value.Denom), "GAMM/POOL/") ||
			strings.HasPrefix(strings.ToUpper(value.Denom), "IBC/") {
			continue
		}
		metadata := GetDenomMetadata(value.Denom, c, client)
		amount, err := strconv.ParseFloat(value.Amount, 1)
		if err != nil {
			return err
		} else {
			if amount > zeroAmount {
				convertedAmount := amount / math.Pow10(metadata.Precision)
				foundToken := false
				for j := range c.Accounts[idx].Tokens {
					if strings.ToLower(c.Accounts[idx].Tokens[j].Denom) == strings.ToLower(value.Denom) {
						c.Accounts[idx].Tokens[j].DelegatedVesting += convertedAmount
						foundToken = true
					}
				}

				if !foundToken {
					c.Accounts[idx].Tokens = append(c.Accounts[idx].Tokens, &model.Token{
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

func (c *Chain) ProcessBankBalances(resp *model.BalancesResponse, idx int, client *http.Client) error {
	for _, value := range resp.Balances {
		if strings.HasPrefix(strings.ToUpper(value.Denom), "GAMM/POOL/") ||
			strings.HasPrefix(strings.ToUpper(value.Denom), "IBC/") {
			continue
		}
		metadata := GetDenomMetadata(value.Denom, c, client)
		amount, err := strconv.ParseFloat(value.Amount, 1)
		if err != nil {
			return err
		} else {
			if amount > zeroAmount {
				convertedAmount := amount / math.Pow10(metadata.Precision)
				foundToken := false
				for j := range c.Accounts[idx].Tokens {
					if strings.ToLower(c.Accounts[idx].Tokens[j].Denom) == strings.ToLower(value.Denom) {
						c.Accounts[idx].Tokens[j].Balance += convertedAmount
						foundToken = true
					}
				}

				if !foundToken {
					c.Accounts[idx].Tokens = append(c.Accounts[idx].Tokens, &model.Token{
						DisplayName: metadata.Symbol,
						Denom:       value.Denom,
						Balance:     convertedAmount,
					})
				}
			}
		}
	}
	return nil
}

func (c *Chain) ProcessRewards(resp *model.RewardsResponse, idx int, client *http.Client) error {
	for _, value := range resp.Total {
		if strings.HasPrefix(strings.ToUpper(value.Denom), "GAMM/POOL/") ||
			strings.HasPrefix(strings.ToUpper(value.Denom), "IBC/") {
			continue
		}
		metadata := GetDenomMetadata(value.Denom, c, client)
		amount, err := strconv.ParseFloat(value.Amount, 1)
		if err != nil {
			return err
		} else {
			if amount > zeroAmount {
				convertedAmount := amount / math.Pow10(metadata.Precision)
				foundToken := false
				for j := range c.Accounts[idx].Tokens {
					if strings.ToLower(c.Accounts[idx].Tokens[j].Denom) == strings.ToLower(value.Denom) {
						c.Accounts[idx].Tokens[j].Reward += convertedAmount
						foundToken = true
					}
				}

				if !foundToken {
					c.Accounts[idx].Tokens = append(c.Accounts[idx].Tokens, &model.Token{
						DisplayName: metadata.Symbol,
						Denom:       value.Denom,
						Reward:      convertedAmount,
					})
				}
			}
		}
	}
	return nil
}

func (c *Chain) ProcessCommission(resp *model.CommissionResponse, idx int, client *http.Client) error {
	for _, value := range resp.Commissions.Commission {
		if strings.HasPrefix(strings.ToUpper(value.Denom), "GAMM/POOL/") ||
			strings.HasPrefix(strings.ToUpper(value.Denom), "IBC/") {
			continue
		}
		metadata := GetDenomMetadata(value.Denom, c, client)
		amount, err := strconv.ParseFloat(value.Amount, 1)
		if err != nil {
			return err
		} else {
			if amount > zeroAmount {
				convertedAmount := amount / math.Pow10(metadata.Precision)
				foundToken := false
				for j := range c.Accounts[idx].Tokens {
					if strings.ToLower(c.Accounts[idx].Tokens[j].Denom) == strings.ToLower(value.Denom) {
						c.Accounts[idx].Tokens[j].Commission += convertedAmount
						foundToken = true
					}
				}

				if !foundToken {
					c.Accounts[idx].Tokens = append(c.Accounts[idx].Tokens, &model.Token{
						DisplayName: metadata.Symbol,
						Denom:       value.Denom,
						Commission:  convertedAmount,
					})
				}
			}
		}
	}
	return nil
}

func (c *Chain) ProcessDelegations(resp *model.Delegations, idx int, client *http.Client) error {
	for _, value := range resp.DelegationResponses {
		if strings.HasPrefix(strings.ToUpper(value.Balance.Denom), "GAMM/POOL/") ||
			strings.HasPrefix(strings.ToUpper(value.Balance.Denom), "IBC/") {
			continue
		}
		metadata := GetDenomMetadata(value.Balance.Denom, c, client)
		amount, err := strconv.ParseFloat(value.Balance.Amount, 1)
		if err != nil {
			return err
		} else {
			if amount > zeroAmount {
				convertedAmount := amount / math.Pow10(metadata.Precision)
				foundToken := false
				for j := range c.Accounts[idx].Tokens {
					if strings.ToLower(c.Accounts[idx].Tokens[j].Denom) == strings.ToLower(value.Balance.Denom) {
						c.Accounts[idx].Tokens[j].Delegation += convertedAmount
						foundToken = true
					}
				}

				if !foundToken {
					c.Accounts[idx].Tokens = append(c.Accounts[idx].Tokens, &model.Token{
						DisplayName: metadata.Symbol,
						Denom:       value.Balance.Denom,
						Delegation:  convertedAmount,
					})
				}
			}
		}
	}
	return nil
}

func (c *Chain) ProcessUnbondings(resp *model.Unbondings, idx int, client *http.Client) error {
	for _, value := range resp.UnbondingResponses {
		for _, entry := range value.Entries {
			metadata := GetDenomMetadata(c.BondDenom, c, client)
			amount, err := strconv.ParseFloat(entry.Balance, 1)
			if err != nil {
				return err
			} else {
				if amount > zeroAmount {
					convertedAmount := amount / math.Pow10(metadata.Precision)
					foundToken := false
					for j := range c.Accounts[idx].Tokens {
						if strings.ToLower(c.Accounts[idx].Tokens[j].Denom) == strings.ToLower(c.BondDenom) {
							c.Accounts[idx].Tokens[j].Unbonding += convertedAmount
							foundToken = true
						}
					}

					if !foundToken {
						c.Accounts[idx].Tokens = append(c.Accounts[idx].Tokens, &model.Token{
							DisplayName: metadata.Symbol,
							Denom:       c.BondDenom,
							Delegation:  convertedAmount,
						})
					}
				}
			}
		}
	}
	return nil
}

// GetDenomMetadata checks if the provided denom is part of a chain's external asset list
// and returns the UI friendly name and exponent
func GetDenomMetadata(denom string, chain *Chain, client *http.Client) TokenDetail {
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
