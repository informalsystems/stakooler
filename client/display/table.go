package display

import (
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"github.com/informalsystems/stakooler/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"os"
	"strings"
)

func PrintAccountDetailsTable(accounts *model.Accounts) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle(strings.ToUpper("Accounts - Details"))
	t.SetCaption(fmt.Sprintf("Retrieved information for %d accounts", len(accounts.Entries)))
	t.AppendHeader(table.Row{"Name", "Account", "Token", "Price (USD)", "Balance", "Rewards", "Staked", "Un-bonding", "Commissions", "Total (Tokens)", "Total (USD)"})

	for acctIdx := range accounts.Entries {
		entries := accounts.Entries[acctIdx].TokensEntry
		// To store the keys in slice in sorted order

		for i := range accounts.Entries[acctIdx].TokensEntry {
			total := entries[i].Balance + entries[i].Reward + entries[i].Delegation + entries[i].Unbonding + entries[i].Commission
			total_price := total * accounts.Entries[acctIdx].TokensEntry[i].Price
			if i == 0 {
				t.AppendRow([]interface{}{
					accounts.Entries[acctIdx].Name,
					accounts.Entries[acctIdx].Address,
					accounts.Entries[acctIdx].TokensEntry[i].DisplayName,
					FilterZeroValue(accounts.Entries[acctIdx].TokensEntry[i].Price),
					FilterZeroValue(accounts.Entries[acctIdx].TokensEntry[i].Balance),
					FilterZeroValue(accounts.Entries[acctIdx].TokensEntry[i].Reward),
					FilterZeroValue(accounts.Entries[acctIdx].TokensEntry[i].Delegation),
					FilterZeroValue(accounts.Entries[acctIdx].TokensEntry[i].Unbonding),
					FilterZeroValue(accounts.Entries[acctIdx].TokensEntry[i].Commission),
					FilterZeroValue(total),
					FilterZeroValue(total_price),
				})
			} else {
				t.AppendRow([]interface{}{
					"",
					"",
					accounts.Entries[acctIdx].TokensEntry[i].DisplayName,
					FilterZeroValue(accounts.Entries[acctIdx].TokensEntry[i].Price),
					FilterZeroValue(accounts.Entries[acctIdx].TokensEntry[i].Balance),
					FilterZeroValue(accounts.Entries[acctIdx].TokensEntry[i].Reward),
					FilterZeroValue(accounts.Entries[acctIdx].TokensEntry[i].Delegation),
					FilterZeroValue(accounts.Entries[acctIdx].TokensEntry[i].Unbonding),
					FilterZeroValue(accounts.Entries[acctIdx].TokensEntry[i].Commission),
					FilterZeroValue(total),
					FilterZeroValue(total_price),
				})
			}

		}
		t.AppendSeparator()
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Name", Align: text.AlignLeft, AlignHeader: text.AlignCenter},
		{Name: "Account", Align: text.AlignLeft, AlignHeader: text.AlignCenter},
		{Name: "Token", Align: text.AlignLeft, AlignHeader: text.AlignCenter},
		{Name: "Price (USD)", Align: text.AlignLeft, AlignHeader: text.AlignCenter},
		{Name: "Balance", Align: text.AlignRight, AlignHeader: text.AlignCenter},
		{Name: "Rewards", Align: text.AlignRight, AlignHeader: text.AlignCenter},
		{Name: "Staked", Align: text.AlignRight, AlignHeader: text.AlignCenter},
		{Name: "Unbonding", Align: text.AlignRight, AlignHeader: text.AlignCenter},
		{Name: "Commissions", Align: text.AlignRight, AlignHeader: text.AlignCenter},
		{Name: "Total", Align: text.AlignRight, AlignHeader: text.AlignCenter},
	})
	t.Render()
	return
}

func FilterZeroValue(value float64) string {
	if value > utils.ZEROAMOUNT {
		return fmt.Sprintf("%f", value)
	} else {
		return ""
	}
}
