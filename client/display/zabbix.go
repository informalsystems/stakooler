package display

import (
	"fmt"
	"log"
	"time"

	sender "github.com/adubkov/go-zabbix"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"github.com/spf13/cast"
)

func ZbxSendValStats(server string, port int, host string, validators *model.Validators) {
	for _, validator := range validators.Entries {
		var metrics []*sender.Metric

		metrics = append(metrics, sender.NewMetric(host, "validator.data.discovery", validator.Chain.ID, validator.BlockTime.Unix()))
		metrics = append(metrics, sender.NewMetric(validator.Chain.ID, "validator.stats.moniker", validator.Moniker, validator.BlockTime.Unix()))
		metrics = append(metrics, sender.NewMetric(validator.Chain.ID, "validator.stats.valoper", validator.ValoperAddress, validator.BlockTime.Unix()))
		metrics = append(metrics, sender.NewMetric(validator.Chain.ID, "validator.stats.block.height", validator.BlockHeight, validator.BlockTime.Unix()))
		metrics = append(metrics, sender.NewMetric(validator.Chain.ID, "validator.stats.voting.power", fmt.Sprintf("%d", validator.VotingPower), validator.BlockTime.Unix()))
		metrics = append(metrics, sender.NewMetric(validator.Chain.ID, "validator.stats.voting.percent", fmt.Sprintf("%.2f", validator.VotingPercent), validator.BlockTime.Unix()))
		metrics = append(metrics, sender.NewMetric(validator.Chain.ID, "validator.stats.ranking", fmt.Sprintf("%d", validator.Ranking), validator.BlockTime.Unix()))
		metrics = append(metrics, sender.NewMetric(validator.Chain.ID, "validator.stats.commission", fmt.Sprintf("%.2f", validator.Commission), validator.BlockTime.Unix()))
		metrics = append(metrics, sender.NewMetric(validator.Chain.ID, "validator.stats.delegators", validator.NumDelegators, validator.BlockTime.Unix()))
		metrics = append(metrics, sender.NewMetric(validator.Chain.ID, "validator.stats.unbondings", fmt.Sprintf("%d", validator.Unbondings), validator.BlockTime.Unix()))

		packet := sender.NewPacket(metrics)
		z := sender.NewSender(server, port)

		resp, err := z.Send(packet)
		if err != nil {
			log.Fatalf("Zabbix send failed: %v", err)
		}

		fmt.Println(cast.ToString(resp))
	}
}

func ZbxSendAcctDetails(server string, port int, host string, accounts *model.Accounts) {
	for _, account := range accounts.Entries {
		var metrics []*sender.Metric

		metrics = append(metrics, sender.NewMetric(host, "validator.data.discovery", account.Chain.ID, time.Now().Unix()))
		metrics = append(metrics, sender.NewMetric(account.Chain.ID, "validator.account.name", account.Name, time.Now().Unix()))
		metrics = append(metrics, sender.NewMetric(account.Chain.ID, "validator.account.address", account.Address, time.Now().Unix()))

		for _, token := range account.TokensEntry {
			metrics = append(metrics, sender.NewMetric(account.Chain.ID, "validator.account.balance", fmt.Sprintf("%.6f", token.Balance), time.Now().Unix()))
			metrics = append(metrics, sender.NewMetric(account.Chain.ID, "validator.account.rewards", fmt.Sprintf("%.6f", token.Reward), time.Now().Unix()))
			metrics = append(metrics, sender.NewMetric(account.Chain.ID, "validator.account.staked", fmt.Sprintf("%.6f", token.Delegation), time.Now().Unix()))
			metrics = append(metrics, sender.NewMetric(account.Chain.ID, "validator.account.unbonding", fmt.Sprintf("%.6f", token.Unbonding), time.Now().Unix()))
			metrics = append(metrics, sender.NewMetric(account.Chain.ID, "validator.account.commission", fmt.Sprintf("%.6f", token.Commission), time.Now().Unix()))

			total := token.Balance + token.Reward + token.Delegation + token.Commission
			metrics = append(metrics, sender.NewMetric(account.Chain.ID, "validator.account.total", fmt.Sprintf("%.6f", total), time.Now().Unix()))
		}

		packet := sender.NewPacket(metrics)
		z := sender.NewSender(server, port)

		resp, err := z.Send(packet)
		if err != nil {
			log.Fatalf("Zabbix send failed: %v", err)
		}

		fmt.Println(cast.ToString(resp))
	}
}
