package display

import (
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"os"
	"strings"
	"time"
)

func PrintAccountDetailsTable(accounts *model.Accounts) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle(strings.ToUpper("Accounts - Details"))
	t.SetCaption(fmt.Sprintf("Fetched: %s", time.Now().Format("2006-01-02 15:04:05")))
	t.AppendHeader(table.Row{"Name", "Account", "Token", "Balance", "Rewards", "Staked", "Unbonding", "Commissions", "Total"})

	for acctIdx := range accounts.Entries {
		entries := accounts.Entries[acctIdx].TokensEntry
		// To store the keys in slice in sorted order

		for i := range accounts.Entries[acctIdx].TokensEntry {
			total := entries[i].Balance + entries[i].Reward + entries[i].Delegation + entries[i].Unbonding + entries[i].Commission
			if i == 0 {
				t.AppendRow([]interface{}{
					accounts.Entries[acctIdx].Name,
					accounts.Entries[acctIdx].Address,
					accounts.Entries[acctIdx].TokensEntry[i].DisplayName,
					FilterZeroValue(accounts.Entries[acctIdx].TokensEntry[i].Balance),
					FilterZeroValue(accounts.Entries[acctIdx].TokensEntry[i].Reward),
					FilterZeroValue(accounts.Entries[acctIdx].TokensEntry[i].Delegation),
					FilterZeroValue(accounts.Entries[acctIdx].TokensEntry[i].Unbonding),
					FilterZeroValue(accounts.Entries[acctIdx].TokensEntry[i].Commission),
					FilterZeroValue(total),
				})
			} else {
				t.AppendRow([]interface{}{
					"",
					"",
					accounts.Entries[acctIdx].TokensEntry[i].DisplayName,
					FilterZeroValue(accounts.Entries[acctIdx].TokensEntry[i].Balance),
					FilterZeroValue(accounts.Entries[acctIdx].TokensEntry[i].Reward),
					FilterZeroValue(accounts.Entries[acctIdx].TokensEntry[i].Delegation),
					FilterZeroValue(accounts.Entries[acctIdx].TokensEntry[i].Unbonding),
					FilterZeroValue(accounts.Entries[acctIdx].TokensEntry[i].Commission),
					FilterZeroValue(total),
				})
			}

		}
		t.AppendSeparator()
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Name", Align: text.AlignLeft, AlignHeader: text.AlignCenter},
		{Name: "Account", Align: text.AlignLeft, AlignHeader: text.AlignCenter},
		{Name: "Token", Align: text.AlignLeft, AlignHeader: text.AlignCenter},
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
	if value > 0.00000 {
		return fmt.Sprintf("%f", value)
	} else {
		return ""
	}
}
