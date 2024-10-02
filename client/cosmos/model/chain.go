package model

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/informalsystems/stakooler/client/cosmos/api"
	"github.com/rs/zerolog/log"
)

const zeroAmount = 0.00000

type Chain struct {
	Name         string
	Id           string
	RestEndpoint string
	Bech32Prefix string
	Accounts     []*Account
	BondDenom    string
	Exponent     int
	AssetList    *api.AssetList
}

func (c *Chain) FetchAccountBalances(blockInfo api.BlockResponse, client *http.Client) error {
	for idx := range c.Accounts {
		c.Accounts[idx].BlockTime = blockInfo.Block.Header.Time
		c.Accounts[idx].BlockHeight = blockInfo.Block.Header.Height

		acct := api.AcctResponse{}
		if err := acct.QueryAuth(c.Accounts[idx].Address, c.RestEndpoint, client); err != nil {
			if strings.Contains(err.Error(), "404") {
				log.Error().Msg(fmt.Sprintf("account %s not found", c.Accounts[idx].Name))
				continue
			}
			return errors.New(fmt.Sprintf("query account: %s", err))
		} else {
			if err = c.ParseAcctQueryResp(&acct, idx, client); err != nil {
				return errors.New(fmt.Sprintf("process vesting: %s", err))
			}
		}

		bank := &api.BankResponse{}
		if err := bank.QueryBankBalances(c.Accounts[idx].Address, c.RestEndpoint, client); err != nil {
			return errors.New(fmt.Sprintf("query bank balances: %s", err))
		} else {
			if err = c.ParseAcctQueryResp(bank, idx, client); err != nil {
				return errors.New(fmt.Sprintf("process bank balances: %s", err))
			}
		}

		rewards := &api.RewardsResponse{}
		if err := rewards.QueryRewards(c.Accounts[idx].Address, c.RestEndpoint, client); err != nil {
			return errors.New(fmt.Sprintf("query rewards: %s", err))
		} else {
			if err = c.ParseAcctQueryResp(rewards, idx, client); err != nil {
				return errors.New(fmt.Sprintf("process rewards: %s", err))
			}
		}

		commission := &api.CommissionResponse{}
		if err := commission.QueryCommission(c.Accounts[idx].Valoper, c.RestEndpoint, client); err != nil {
			return errors.New(fmt.Sprintf("query commissions: %s", err))
		} else if commission.Commissions.Commission != nil {
			if err = c.ParseAcctQueryResp(commission, idx, client); err != nil {
				return errors.New(fmt.Sprintf("process commissions: %s", err))
			}
		}

		delegation := &api.Delegations{}
		if err := delegation.QueryDelegations(c.Accounts[idx].Address, c.RestEndpoint, client); err != nil {
			return errors.New(fmt.Sprintf("query delegations: %s", err))
		} else {
			if err = c.ParseAcctQueryResp(delegation, idx, client); err != nil {
				return errors.New(fmt.Sprintf("process delegations: %s", err))
			}
		}

		unbondings := &api.Unbondings{}
		if err := unbondings.QueryUnbondings(c.Accounts[idx].Address, c.RestEndpoint, client); err != nil {
			return errors.New(fmt.Sprintf("query unbondings: %s", err))
		} else {
			if err = c.ParseAcctQueryResp(unbondings, idx, client); err != nil {
				return errors.New(fmt.Sprintf("process unbondings: %s", err))
			}
		}
	}
	return nil
}

func (c *Chain) ParseAcctQueryResp(resp api.AccountQueryResponse, idx int, client *http.Client) error {
	for balanceType, balance := range resp.GetBalances() {
		for denom, amount := range balance {
			if strings.HasPrefix(strings.ToUpper(denom), "GAMM/POOL/") ||
				strings.HasPrefix(strings.ToUpper(denom), "IBC/") ||
				strings.HasPrefix(strings.ToUpper(denom), "FACTORY/") ||
				strings.HasPrefix(strings.ToUpper(denom), "ST") {
				continue
			}

			symbol, exponent := GetDenomMetadata(denom, c, client)
			floatAmount, err := strconv.ParseFloat(amount, 1)
			if err != nil {
				return err
			}

			if floatAmount > zeroAmount {
				convertedAmmount := floatAmount / math.Pow10(exponent)

				if _, ok := c.Accounts[idx].Tokens[denom]; !ok {
					c.Accounts[idx].Tokens[denom] = &Token{
						DisplayName: symbol,
						Denom:       denom,
					}
				}

				switch balanceType {
				case api.OriginalVesting:
					c.Accounts[idx].Tokens[denom].Balances.OriginalVesting += convertedAmmount
				case api.DelegatedVesting:
					c.Accounts[idx].Tokens[denom].Balances.DelegatedVesting += convertedAmmount
				case api.Bank:
					c.Accounts[idx].Tokens[denom].Balances.Bank += convertedAmmount
				case api.Rewards:
					c.Accounts[idx].Tokens[denom].Balances.Rewards += convertedAmmount
				case api.Commission:
					c.Accounts[idx].Tokens[denom].Balances.Commission += convertedAmmount
				case api.Delegation:
					c.Accounts[idx].Tokens[denom].Balances.Delegated += convertedAmmount
				case api.Unbonding:
					c.Accounts[idx].Tokens[denom].Balances.Unbonding += convertedAmmount
				}
			}
		}
	}
	return nil
}

// GetDenomMetadata checks if the provided denom is part of a chain's external asset list
// and returns the UI friendly name and exponent
func GetDenomMetadata(denom string, chain *Chain, client *http.Client) (string, int) {
	var symbol string
	exponent := 0

	if chain.AssetList != nil {
		symbol, exponent = chain.AssetList.SearchForAsset(denom)
	}
	// in case asset details are missing
	if exponent == 0 {
		denomMetadata := &api.DenomMetadataResponse{}
		if err := denomMetadata.QueryMetadataFromBank(denom, chain.RestEndpoint, client); err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("getting metadata from bank for: %s", denom))
		}
		if denomMetadata.Metadata.Base == "" {
			symbol = denom
			exponent = 6
		} else {
			symbol = strings.ToUpper(denomMetadata.Metadata.Display)
			exponent = denomMetadata.GetExponent()
		}

	}

	// If it's an IBC denom add '(IBC)' to the symbol
	if strings.HasPrefix(strings.ToUpper(denom), "IBC/") {
		symbol = symbol + " (IBC)"
	}
	return symbol, exponent
}
