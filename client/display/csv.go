package display

import (
	"encoding/csv"
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"log"
	"os"
)

func WriteCSV(accounts *model.Accounts) {
	f, err := os.Create("accounts.csv")
	defer f.Close()

	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()

	header := []string{"Name", "Account", "Token", "Balance", "Rewards", "Staked", "Unbonding", "Commissions", "Total"}
	if err := w.Write(header); err != nil {
		log.Fatalln("error writing record to file", err)
	}
	for acctIdx := range accounts.Entries {
		entries := accounts.Entries[acctIdx].TokensEntry
		for i := range accounts.Entries[acctIdx].TokensEntry {
			total := entries[i].Balance + entries[i].Reward + entries[i].Delegation + entries[i].Unbonding + entries[i].Commission
			record := []string{
				accounts.Entries[acctIdx].Name,
				accounts.Entries[acctIdx].Address,
				accounts.Entries[acctIdx].TokensEntry[i].DisplayName,
				fmt.Sprintf("%f", accounts.Entries[acctIdx].TokensEntry[i].Balance),
				fmt.Sprintf("%f", accounts.Entries[acctIdx].TokensEntry[i].Reward),
				fmt.Sprintf("%f", accounts.Entries[acctIdx].TokensEntry[i].Delegation),
				fmt.Sprintf("%f", accounts.Entries[acctIdx].TokensEntry[i].Unbonding),
				fmt.Sprintf("%f", accounts.Entries[acctIdx].TokensEntry[i].Commission),
				fmt.Sprintf("%f", total),
			}
			if err := w.Write(record); err != nil {
				log.Fatalln("error writing record to file", err)
			}
		}
	}
}
