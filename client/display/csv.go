package display

import (
	"encoding/csv"
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"log"
	"os"
	"time"
)

func WriteAccountsCSV(accounts *model.Accounts) {

	// Outputs to Stdout
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	header := []string{"account_name", "account_address", "chain_id", "block_height", "block_time", "token", "balance", "rewards", "staked", "unbonding", "commissions", "total"}
	if err := w.Write(header); err != nil {
		log.Fatalln("error writing record to file", err)
	}
	for acctIdx := range accounts.Entries {
		entries := accounts.Entries[acctIdx].TokensEntry

		// In case there is no token information
		if len(entries) == 0 {
			record := []string{
				accounts.Entries[acctIdx].Name, accounts.Entries[acctIdx].Address, "na", "na", "na", "na", "na", "na",
				"na", "na", "na", "na",
			}
			if err := w.Write(record); err != nil {
				log.Fatalln("error writing record", err)
			}
		} else {
			for i := range accounts.Entries[acctIdx].TokensEntry {
				total := entries[i].Balance + entries[i].Reward + entries[i].Delegation + entries[i].Unbonding + entries[i].Commission
				record := []string{
					accounts.Entries[acctIdx].Name,
					accounts.Entries[acctIdx].Address,
					accounts.Entries[acctIdx].Chain.ID,
					accounts.Entries[acctIdx].TokensEntry[i].BlockHeight,
					accounts.Entries[acctIdx].TokensEntry[i].BlockTime.Format(time.RFC3339Nano),
					accounts.Entries[acctIdx].TokensEntry[i].DisplayName,
					fmt.Sprintf("%f", accounts.Entries[acctIdx].TokensEntry[i].Balance),
					fmt.Sprintf("%f", accounts.Entries[acctIdx].TokensEntry[i].Reward),
					fmt.Sprintf("%f", accounts.Entries[acctIdx].TokensEntry[i].Delegation),
					fmt.Sprintf("%f", accounts.Entries[acctIdx].TokensEntry[i].Unbonding),
					fmt.Sprintf("%f", accounts.Entries[acctIdx].TokensEntry[i].Commission),
					fmt.Sprintf("%f", total),
				}
				if err := w.Write(record); err != nil {
					log.Fatalln("error writing record", err)
				}
			}
		}
	}
}

func WriteValidatorCSV(validators *model.Validators) {

	// Outputs to Stdout
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	header := []string{"moniker", "chain_id", "valoper_address", "block_time", "block_height", "voting_power_tokens", "voting_power_percent", "ranking", "delegators", "unbondings"}
	if err := w.Write(header); err != nil {
		log.Fatalln("error writing record to file", err)
	}
	for idx := range validators.Entries {
		validator := validators.Entries[idx]

		record := []string{
			validator.Moniker,
			validator.Chain.ID,
			validator.ValoperAddress,
			validator.BlockTime.Format(time.RFC822),
			validator.BlockHeight,
			fmt.Sprintf("%d", validator.VotingPower),
			fmt.Sprintf("%.2f", validator.VotingPercent),
			fmt.Sprintf("%d", validator.Ranking),
			validator.NumDelegators,
			fmt.Sprintf("%d", validator.Unbondings),
		}
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record", err)
		}
	}
}
