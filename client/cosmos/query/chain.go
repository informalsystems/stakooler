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
	AssetList    *api.AssetList
}

func (c *Chain) FetchAccountBalances(blockInfo api.BlockResponse, client *http.Client) error {
	for idx := range c.Accounts {
		c.Accounts[idx].BlockTime = blockInfo.Block.Header.Time
		c.Accounts[idx].BlockHeight = blockInfo.Block.Header.Height

		if acctResponse, err := api.GetAccount(c.Accounts[idx].Address, c.RestEndpoint, client); err != nil {
			return errors.New(fmt.Sprintf("query account: %s", err))
		} else {
			if err = c.ProcessResponse(acctResponse, idx, client); err != nil {
				return errors.New(fmt.Sprintf("process vesting: %s", err))
			}
		}

		if bankResponse, err := api.GetBalances(c.Accounts[idx].Address, c.RestEndpoint, client); err != nil {
			return errors.New(fmt.Sprintf("query bank balances: %s", err))
		} else {
			if err = c.ProcessResponse(bankResponse, idx, client); err != nil {
				return errors.New(fmt.Sprintf("process bank balances: %s", err))
			}
		}

		if rewardsResponse, err := api.GetRewards(c.Accounts[idx].Address, c.RestEndpoint, client); err != nil {
			return errors.New(fmt.Sprintf("query rewards: %s", err))
		} else {
			if err = c.ProcessResponse(rewardsResponse, idx, client); err != nil {
				return errors.New(fmt.Sprintf("process rewards: %s", err))
			}
		}

		if commissionResponse, err := api.GetCommissions(c.Accounts[idx].Valoper, c.RestEndpoint, client); err != nil {
			return errors.New(fmt.Sprintf("query commissions: %s", err))
		} else if commissionResponse != nil {
			if err = c.ProcessResponse(commissionResponse, idx, client); err != nil {
				return errors.New(fmt.Sprintf("process commissions: %s", err))
			}
		}

		if delegationsResponse, err := api.GetDelegations(c.Accounts[idx].Address, c.RestEndpoint, client); err != nil {
			return errors.New(fmt.Sprintf("query delegations: %s", err))
		} else {
			if err = c.ProcessResponse(delegationsResponse, idx, client); err != nil {
				return errors.New(fmt.Sprintf("process delegations: %s", err))
			}
		}

		if unbondingResponse, err := api.GetUnbondings(c.Accounts[idx].Address, c.RestEndpoint, client); err != nil {
			return errors.New(fmt.Sprintf("query unbondings: %s", err))
		} else {
			if err = c.ProcessResponse(unbondingResponse, idx, client); err != nil {
				return errors.New(fmt.Sprintf("process unbondings: %s", err))
			}
		}
	}
	return nil
}

func (c *Chain) ProcessResponse(resp api.Response, idx int, client *http.Client) error {
	for balanceType, balance := range resp.GetBalances() {
		for denom, amount := range balance {
			if strings.HasPrefix(strings.ToUpper(denom), "GAMM/POOL/") ||
				strings.HasPrefix(strings.ToUpper(denom), "IBC/") {
				continue
			}

			metadata := GetDenomMetadata(denom, c, client)
			floatAmount, err := strconv.ParseFloat(amount, 1)
			if err != nil {
				return err
			}

			if floatAmount > zeroAmount {
				convertedAmmount := floatAmount / math.Pow10(metadata.Precision)

				if _, ok := c.Accounts[idx].Tokens[denom]; !ok {
					c.Accounts[idx].Tokens[denom] = &model.Token{
						DisplayName: metadata.Symbol,
						Denom:       denom,
					}
				}

				switch balanceType {
				case api.OriginalVesting:
					c.Accounts[idx].Tokens[denom].OriginalVesting += convertedAmmount
				case api.DelegatedVesting:
					c.Accounts[idx].Tokens[denom].DelegatedVesting += convertedAmmount
				case api.Bank:
					c.Accounts[idx].Tokens[denom].BankBalance += convertedAmmount
				case api.Rewards:
					c.Accounts[idx].Tokens[denom].Rewards += convertedAmmount
				case api.Commission:
					c.Accounts[idx].Tokens[denom].Commission += convertedAmmount
				case api.Delegation:
					c.Accounts[idx].Tokens[denom].Delegation += convertedAmmount
				case api.Unbonding:
					c.Accounts[idx].Tokens[denom].Unbonding += convertedAmmount
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
		symbol, exponent = chain.AssetList.SearchForAsset(denom)
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
