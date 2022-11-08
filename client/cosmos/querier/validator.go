package querier

import (
	"errors"
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/api"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"github.com/schollz/progressbar/v3"
	"sort"
	"strconv"
	"strings"
)

func LoadValidatorStats(validator *model.Validator, bar *progressbar.ProgressBar) error {

	// Get validators
	validators, err := api.GetValidators(validator)
	if err != nil {
		return err
	}
	bar.Add(1)

	var totalVotingPower int64

	// Sort validators by voting power (descending)
	sort.Slice(validators.ValidatorsResponse, func(i, j int) bool {
		tI, _ := strconv.ParseInt(validators.ValidatorsResponse[i].Tokens[:len(validators.ValidatorsResponse[i].Tokens)-validator.Chain.Exponent], 10, 64)
		tJ, _ := strconv.ParseInt(validators.ValidatorsResponse[j].Tokens[:len(validators.ValidatorsResponse[j].Tokens)-validator.Chain.Exponent], 10, 64)
		return tI >= tJ
	})
	
	// Get total voting power
	for i, val := range validators.ValidatorsResponse {
		tokenConverted, err := strconv.ParseInt(val.Tokens[:len(val.Tokens)-validator.Chain.Exponent], 10, 64)
		if err != nil {
			return errors.New(fmt.Sprintf("cannot convert tokens for voting power: %s", err))
		}
		if strings.ToLower(val.OperatorAddress) == strings.ToLower(validator.ValoperAddress) {
			validator.VotingPower = tokenConverted
			validator.Ranking = i + 1
			validator.Name = val.Description.Moniker
		}
		totalVotingPower += tokenConverted
	}

	// Find the voting power percent
	votingPowerShare := float64(validator.VotingPower) / float64(totalVotingPower) * 100.0
	validator.VotingPercent = votingPowerShare

	// Get unbondings
	unbondings, err := api.GetValidatorUnbondings(validator)
	if err != nil {
		return err
	}
	bar.Add(1)

	var totalUnbondings int64

	// Get total voting power
	for _, unbonding := range unbondings.UnbondingResponses {

		for _, entry := range unbonding.Entries {
			if len(entry.Balance) > validator.Chain.Exponent {
				unbonding := entry.Balance[:len(entry.Balance)-validator.Chain.Exponent]
				if unbonding != "" {
					unbondingConverted, err := strconv.ParseInt(unbonding, 10, 64)
					if err != nil {
						return errors.New(fmt.Sprintf("cannot convert unbondings: %s", err))
					}
					totalUnbondings += unbondingConverted
				}
			}
		}
	}
	validator.Unbondings = totalUnbondings

	// Get block time
	validator.BlockHeight = validators.BlockHeight
	block, _ := api.GetBlock(validator.BlockHeight, validator.Chain)
	bar.Add(1)

	validator.BlockTime = block.Block.Header.Time

	// Get number of delegators
	delegations, _ := api.GetValidatorDelegations(validator)
	bar.Add(1)
	validator.NumDelegators = delegations.Pagination.Total

	// Get number of total validators
	validator.NumValidators = validators.Pagination.Total

	return nil
}
