package querier

import (
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/api"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"github.com/schollz/progressbar/v3"
	"strconv"
	"strings"
)

func LoadValidatorStats(validator *model.Validator, bar *progressbar.ProgressBar) error {
	bar.Add(1)
	validators, err := api.GetValidators(validator)
	if err != nil {
		return err
	}

	var totalVotingPower int64

	// Get total voting power
	for i, val := range validators.ValidatorsResponse {
		if strings.ToLower(val.OperatorAddress) == strings.ToLower(validator.ValoperAddress) {
			validator.VotingPower = val.Tokens
			validator.Ranking = i + 1
		}
		totalVotingPower += val.Tokens
	}

	// Find the voting power percent
	votingPowerShare := float64(validator.VotingPower) / float64(totalVotingPower) * 100.0
	validator.VotingPercent = votingPowerShare
	bar.Add(1)
	unbondings, err := api.GetValidatorUnbondings(validator)
	if err != nil {
		return err
	}

	var totalUnbondings int64

	// Get total voting power
	for _, unbonding := range unbondings.UnbondingResponses {
		for _, entry := range unbonding.Entries {
			b, err := strconv.Atoi(entry.Balance)
			if err != nil {
				fmt.Sprintf("Error voting power: %s", err)
			}
			totalUnbondings += int64(b)
		}
	}
	validator.Unbondings = totalUnbondings
	bar.Add(1)
	// Get block time
	validator.BlockHeight = validators.BlockHeight
	block, _ := api.GetBlock(validator.BlockHeight, validator.Chain)
	validator.BlockTime = block.Block.Header.Time
	bar.Add(1)
	// Get number of delegators
	delegations, _ := api.GetValidatorDelegations(validator)
	validator.NumDelegators = delegations.Pagination.Total
	return nil
}
