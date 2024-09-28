package display

import (
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"os"
	"strings"
	"time"
)

func PrintAccountDetailsTable(chains *model.Chains) {
	for _, chain := range chains.Entries {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetTitle(strings.ToUpper(fmt.Sprintf("Details for %s", chain.Name)))
		t.SetCaption(fmt.Sprintf("Retrieved information for %d accounts", len(chain.Accounts)))
		t.AppendHeader(table.Row{"Name", "Account", "Token", "Balance", "Rewards", "Staked", "Unbonding", "Commissions", "Original Vesting", "Delegated Vesting", "Total"})

		for _, account := range chain.Accounts {
			for idx, e := range account.TokensEntry {
				total := e.Vesting - e.DelegatedVesting + e.Balance + e.Reward + e.Delegation + e.Unbonding + e.Commission
				if idx == 0 {
					t.AppendRow([]interface{}{
						account.Name,
						account.Address,
						e.DisplayName,
						FilterZeroValue(e.Balance),
						FilterZeroValue(e.Reward),
						FilterZeroValue(e.Delegation),
						FilterZeroValue(e.Unbonding),
						FilterZeroValue(e.Commission),
						FilterZeroValue(e.Vesting),
						FilterZeroValue(e.DelegatedVesting),
						FilterZeroValue(total),
					})
				} else {
					t.AppendRow([]interface{}{
						"",
						"",
						FilterZeroValue(e.Balance),
						FilterZeroValue(e.Reward),
						FilterZeroValue(e.Delegation),
						FilterZeroValue(e.Unbonding),
						FilterZeroValue(e.Commission),
						FilterZeroValue(e.Vesting),
						FilterZeroValue(e.DelegatedVesting),
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
			{Name: "Original Vesting", Align: text.AlignRight, AlignHeader: text.AlignCenter},
			{Name: "Delegated Vesting", Align: text.AlignRight, AlignHeader: text.AlignCenter},
			{Name: "Total", Align: text.AlignRight, AlignHeader: text.AlignCenter},
		})
		t.Render()
	}
	return
}

func PrintValidatorStasTable(validators *model.ValidatorList) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle(strings.ToUpper("Validator - Statistics"))
	t.SetCaption(fmt.Sprintf("Retrieved information for %d validators", len(validators.Entries)))
	t.AppendHeader(table.Row{"Moniker", "Chain", "Validator Address", "Block Time", "Block Height", "Voting Power (VP)", "VP (%)", "Ranking", "Commission", "# Validators", "Delegators", "Unbondings"})

	for idx := range validators.Entries {
		validator := validators.Entries[idx]
		p := message.NewPrinter(language.English)

		t.AppendRow([]interface{}{
			validator.Moniker,
			validator.Chain.Id,
			validator.ValoperAddress,
			validator.BlockTime.Format(time.RFC822),
			validator.BlockHeight,
			p.Sprintf("%d (%s)", validator.VotingPower, validator.Chain.Denom),
			p.Sprintf("%.2f", validator.VotingPercent),
			p.Sprintf("%d", validator.Ranking),
			p.Sprintf("%.2f", validator.Commission),
			validator.NumValidators,
			validator.NumDelegators,
			p.Sprintf("%d (%s)", validator.Unbondings, validator.Chain.Denom),
		})
		t.AppendSeparator()
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Moniker", Align: text.AlignLeft, AlignHeader: text.AlignCenter},
		{Name: "Chain", Align: text.AlignLeft, AlignHeader: text.AlignCenter},
		{Name: "Validator Address", Align: text.AlignLeft, AlignHeader: text.AlignCenter},
		{Name: "Block Time", Align: text.AlignLeft, AlignHeader: text.AlignCenter},
		{Name: "Block Height", Align: text.AlignRight, AlignHeader: text.AlignCenter},
		{Name: "Voting Power (VP)", Align: text.AlignRight, AlignHeader: text.AlignCenter},
		{Name: "VP (%)", Align: text.AlignRight, AlignHeader: text.AlignCenter},
		{Name: "Ranking", Align: text.AlignRight, AlignHeader: text.AlignCenter},
		{Name: "Commission", Align: text.AlignRight, AlignHeader: text.AlignCenter},
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
