package display

import (
	"encoding/csv"
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"log"
	"os"
	"time"
)

func WriteCSV(accounts *model.Accounts) {

	// Outputs to Stdout
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	header := []string{"account_name", "account_address", "chain_id", "block_height", "block_time", "token", "token_usd", "balance", "rewards", "staked", "unbonding", "commissions", "total", "total_usd"}
	if err := w.Write(header); err != nil {
		log.Fatalln("error writing record to file", err)
	}
	for acctIdx := range accounts.Entries {
		entries := accounts.Entries[acctIdx].TokensEntry

		// In case there is no token information
		if len(entries) == 0 {
			record := []string{
				accounts.Entries[acctIdx].Name, accounts.Entries[acctIdx].Address, "na", "na", "na", "na", "na", "na", "na",
				"na", "na", "na", "na",
			}
			if err := w.Write(record); err != nil {
				log.Fatalln("error writing record", err)
			}
		} else {
			for i := range accounts.Entries[acctIdx].TokensEntry {
				record := []string{
					accounts.Entries[acctIdx].Name,
					accounts.Entries[acctIdx].Address,
					accounts.Entries[acctIdx].Chain.ID,
					accounts.Entries[acctIdx].TokensEntry[i].BlockHeight,
					accounts.Entries[acctIdx].TokensEntry[i].BlockTime.Format(time.RFC3339Nano),
					accounts.Entries[acctIdx].TokensEntry[i].DisplayName,
					fmt.Sprintf("%f", accounts.Entries[acctIdx].TokensEntry[i].Price),
					fmt.Sprintf("%f", accounts.Entries[acctIdx].TokensEntry[i].Balance),
					fmt.Sprintf("%f", accounts.Entries[acctIdx].TokensEntry[i].Reward),
					fmt.Sprintf("%f", accounts.Entries[acctIdx].TokensEntry[i].Delegation),
					fmt.Sprintf("%f", accounts.Entries[acctIdx].TokensEntry[i].Unbonding),
					fmt.Sprintf("%f", accounts.Entries[acctIdx].TokensEntry[i].Commission),
					fmt.Sprintf("%f", accounts.Entries[acctIdx].TokensEntry[i].Total),
					fmt.Sprintf("%f", accounts.Entries[acctIdx].TokensEntry[i].TotalPrice),
				}
				if err := w.Write(record); err != nil {
					log.Fatalln("error writing record", err)
				}
			}
		}
	}
}
