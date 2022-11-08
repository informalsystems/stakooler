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
	t.SetCaption(fmt.Sprintf("Retrieved information for %d accounts", len(accounts.Entries)))
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

func PrintValidatorStasTable(validators *model.Validators) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle(strings.ToUpper("Validator - Statistics"))
	t.SetCaption(fmt.Sprintf("Retrieved information for %d validators", len(validators.Entries)))
	t.AppendHeader(table.Row{"Moniker", "Chain", "Validator Address", "Block Time", "Block Height", "Voting Power (VP)", "VP (%)", "Ranking", "# Validators", "Delegators", "Unbondings"})

	for idx := range validators.Entries {
		validator := validators.Entries[idx]
		t.AppendRow([]interface{}{
			validator.Name,
			validator.Chain.ID,
			validator.ValoperAddress,
			validator.BlockTime.Format(time.RFC822),
			validator.BlockHeight,
			fmt.Sprintf("%d (%s)", validator.VotingPower, validator.Chain.Denom),
			fmt.Sprintf("%.2f", validator.VotingPercent),
			fmt.Sprintf("%d", validator.Ranking),
			validator.NumValidators,
			validator.NumDelegators,
			fmt.Sprintf("%d (%s)", validator.Unbondings, validator.Chain.Denom),
		})
		t.AppendSeparator()
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Moniker", Align: text.AlignLeft, AlignHeader: text.AlignCenter},
		{Name: "Chain", Align: text.AlignLeft, AlignHeader: text.AlignCenter},
		{Name: "Validator Address", Align: text.AlignLeft, AlignHeader: text.AlignCenter},
		{Name: "Block Time", Align: text.AlignLeft, AlignHeader: text.AlignCenter},
		{Name: "Block Height", Align: text.AlignRight, AlignHeader: text.AlignCenter},
		{Name: "Voting Power", Align: text.AlignRight, AlignHeader: text.AlignCenter},
		{Name: "Voting Power (%)", Align: text.AlignRight, AlignHeader: text.AlignCenter},
		{Name: "Ranking", Align: text.AlignRight, AlignHeader: text.AlignCenter},
		{Name: "# Validators", Align: text.AlignRight, AlignHeader: text.AlignCenter},
		{Name: "Delegators", Align: text.AlignRight, AlignHeader: text.AlignCenter},
		{Name: "Unbondings", Align: text.AlignRight, AlignHeader: text.AlignCenter},
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
