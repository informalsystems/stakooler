package display

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/informalsystems/stakooler/client/cosmos/model"
)

func WriteDollarValueReport(chains []*model.Chain) {
	file, err := os.Create("dollar_value_report.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer func(File *os.File) {
		err = File.Close()
		if err != nil {
			return
		}
	}(file)

	w := csv.NewWriter(file)
	defer w.Flush()

	header := []string{"account_name", "token", "rewards", "commissions", "total USD value", "total CAD value"}
	if err = w.Write(header); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	accounts := make(map[string][]*model.Token)
	for _, chain := range chains {
		for _, account := range chain.Accounts {
			if _, ok := accounts[account.Name]; !ok {
				accounts[account.Name] = make([]*model.Token, 0)
			}
			for _, token := range account.Tokens {
				accounts[account.Name] = append(accounts[account.Name], token)
			}
		}
	}

	for name, account := range accounts {
		for _, token := range account {
			record := []string{
				name,
				token.DisplayName,
				fmt.Sprintf("%f", token.Balances.Rewards),
				fmt.Sprintf("%f", token.Balances.Commission),
				fmt.Sprintf("%f", (token.Balances.Rewards+token.Balances.Commission)*token.PriceUSD),
				fmt.Sprintf("%f", (token.Balances.Rewards+token.Balances.Commission)*token.PriceCAD),
			}
			if err = w.Write(record); err != nil {
				log.Fatalln("error writing record", err)
			}
		}
	}
}

func WriteAccountsCSV(chains []*model.Chain) {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	header := []string{"account_name", "account_address", "chain_id", "block_height", "block_time", "token", "balance", "rewards", "staked", "unbonding", "commissions", "original_vesting", "delegated_vesting", "total"}
	if err := w.Write(header); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	for _, chain := range chains {
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
					total := entries[i].Balances.Bank +
						entries[i].Balances.Rewards +
						entries[i].Balances.Delegated +
						entries[i].Balances.Unbonding +
						entries[i].Balances.Commission
					record := []string{
						acct.Name,
						acct.Address,
						chain.Id,
						acct.BlockHeight,
						acct.BlockTime.Format(time.DateTime),
						acct.Tokens[i].DisplayName,
						fmt.Sprintf("%f", acct.Tokens[i].Balances.Bank),
						fmt.Sprintf("%f", acct.Tokens[i].Balances.Rewards),
						fmt.Sprintf("%f", acct.Tokens[i].Balances.Delegated),
						fmt.Sprintf("%f", acct.Tokens[i].Balances.Unbonding),
						fmt.Sprintf("%f", acct.Tokens[i].Balances.Commission),
						fmt.Sprintf("%f", acct.Tokens[i].Balances.OriginalVesting),
						fmt.Sprintf("%f", acct.Tokens[i].Balances.DelegatedVesting),
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
