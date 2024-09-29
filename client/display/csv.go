package display

import (
	"encoding/csv"
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/query"
	"log"
	"os"
	"time"
)

func WriteAccountsCSV(chains *query.Chains) {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	header := []string{"account_name", "account_address", "chain_id", "block_height", "block_time", "token", "balance", "rewards", "staked", "unbonding", "commissions", "original_vesting", "delegated_vesting", "total"}
	if err := w.Write(header); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	for _, chain := range chains.Entries {
		for _, acct := range chain.Accounts {
			entries := acct.Tokens

			// In case there is no token information
			if len(entries) == 0 {
				record := []string{
					acct.Name, acct.Address, "na", "na", "na", "na", "na", "na", "na", "na",
					"na", "na", "na", "na",
				}
				if err := w.Write(record); err != nil {
					log.Fatalln("error writing record", err)
				}
			} else {
				for i := range acct.Tokens {
					total := entries[i].BankBalance + entries[i].Rewards + entries[i].Delegation + entries[i].Unbonding + entries[i].Commission
					record := []string{
						acct.Name,
						acct.Address,
						chain.Id,
						acct.BlockHeight,
						acct.BlockTime.Format(time.RFC3339Nano),
						acct.Tokens[i].DisplayName,
						fmt.Sprintf("%f", acct.Tokens[i].BankBalance),
						fmt.Sprintf("%f", acct.Tokens[i].Rewards),
						fmt.Sprintf("%f", acct.Tokens[i].Delegation),
						fmt.Sprintf("%f", acct.Tokens[i].Unbonding),
						fmt.Sprintf("%f", acct.Tokens[i].Commission),
						fmt.Sprintf("%f", acct.Tokens[i].OriginalVesting),
						fmt.Sprintf("%f", acct.Tokens[i].DelegatedVesting),
						fmt.Sprintf("%f", total),
					}
					if err := w.Write(record); err != nil {
						log.Fatalln("error writing record", err)
					}
				}
			}
		}
	}
}

/*func WriteValidatorCSV(validators *model.ValidatorList) {

	// Outputs to Stdout
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	header := []string{"moniker", "chain_id", "valoper_address", "block_time", "block_height", "voting_power_tokens", "voting_power_percent", "ranking", "commission", "delegators", "unbondings"}
	if err := w.Write(header); err != nil {
		log.Fatalln("error writing record to file", err)
	}
	for idx := range validators.Entries {
		validator := validators.Entries[idx]

		record := []string{
			validator.Moniker,
			validator.Chain.Id,
			validator.ValoperAddress,
			validator.BlockTime.Format(time.RFC822),
			validator.BlockHeight,
			fmt.Sprintf("%d", validator.VotingPower),
			fmt.Sprintf("%.2f", validator.VotingPercent),
			fmt.Sprintf("%d", validator.Ranking),
			fmt.Sprintf("%.2f", validator.Commission),
			validator.NumDelegators,
			fmt.Sprintf("%d", validator.Unbondings),
		}
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record", err)
		}
	}
}
*/
