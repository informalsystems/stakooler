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

func PrintTable(accounts *model.Accounts) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle(strings.ToUpper("Accounts Details"))
	t.SetCaption(fmt.Sprintf("Fetched at: %s", time.Now().Format("2006-01-02 15:04:05")))

	t.AppendHeader(table.Row{"Name","Account", "Token", "Balance", "Rewards", "Staked", "Unbonding", "Total"})

	for acctIdx := range accounts.Entries {
		account := accounts.Entries[acctIdx].AccountDetails
		// To store the keys in slice in sorted order
		keys := make([]string, 0, len(account.AvailableBalance))
		for k := range account.AvailableBalance {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for idx, coin := range keys {
			total := account.AvailableBalance[coin] + account.Rewards[coin] + account.Delegations[coin] + account.Unbondings[coin]
			if idx == 0 {
				t.AppendRow([]interface{}{
					accounts.Entries[acctIdx].Name,
					accounts.Entries[acctIdx].Address,
					coin,
					fmt.Sprintf("%f", account.AvailableBalance[coin]),
					fmt.Sprintf("%f", account.Rewards[coin]),
					fmt.Sprintf("%f", account.Delegations[coin]),
					fmt.Sprintf("%f", account.Unbondings[coin]),
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
					fmt.Sprintf("%f", total),
				})
			}
		}
		t.AppendSeparator()
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		// TODO: Align address in the middle
		{Name: "Account", Align: text.AlignCenter, AlignHeader: text.AlignCenter, VAlign: text.VAlignMiddle},

		{Name: "Token", Align: text.AlignLeft, AlignHeader: text.AlignCenter},
		{Name: "Balance", Align: text.AlignRight, AlignHeader: text.AlignCenter},
		{Name: "Rewards", Align: text.AlignRight, AlignHeader: text.AlignCenter},
		{Name: "Staked", Align: text.AlignRight, AlignHeader: text.AlignCenter},
		{Name: "Total", Align: text.AlignRight, AlignHeader: text.AlignCenter},
	})
	t.Render()
	return
}
