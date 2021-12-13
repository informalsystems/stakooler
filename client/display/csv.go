package display

import (
	"encoding/csv"
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

	//for acctIdx := range accounts.Entries {
	//	account := accounts.Entries[acctIdx]
	//	accountDetails := accounts.Entries[acctIdx].TokensEntry
	//	// TODO: refactor account details, don't use a map, use a struct
	//	record := []string{account.Name, account.Address, accountDetails}
	//
	//
	//	if err := w.Write(record); err != nil {
	//		log.Fatalln("error writing record to file", err)
	//	}
	//}
}
