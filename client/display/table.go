package display

import (
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"os"
	"sort"
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
		account := accounts.Entries[acctIdx].AccountDetails
		// To store the keys in slice in sorted order
		keys := make([]string, 0, len(account.AvailableBalance))
		for k := range account.AvailableBalance {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for idx, coin := range keys {
			total := account.AvailableBalance[coin] +
				account.Rewards[coin] +
				account.Delegations[coin] +
				account.Unbondings[coin] +
				account.Commissions[coin]
			if idx == 0 {
				t.AppendRow([]interface{}{
					accounts.Entries[acctIdx].Name,
					accounts.Entries[acctIdx].Address,
					coin,
					fmt.Sprintf("%f", account.AvailableBalance[coin]),
					fmt.Sprintf("%f", account.Rewards[coin]),
					fmt.Sprintf("%f", account.Delegations[coin]),
					fmt.Sprintf("%f", account.Unbondings[coin]),
					fmt.Sprintf("%f", account.Commissions[coin]),
					fmt.Sprintf("%f", total),
				})
			} else {
				t.AppendRow([]interface{}{
					"",
					"",
					coin,
					fmt.Sprintf("%f", account.AvailableBalance[coin]),
					fmt.Sprintf("%f", account.Rewards[coin]),
					fmt.Sprintf("%f", account.Delegations[coin]),
					fmt.Sprintf("%f", account.Unbondings[coin]),
					fmt.Sprintf("%f", account.Commissions[coin]),
					fmt.Sprintf("%f", total),
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
