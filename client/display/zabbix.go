package display

import (
	"fmt"
	"log"

	sender "github.com/adubkov/go-zabbix"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"github.com/spf13/cast"
)

func ZbxSend(host string, port int, validators *model.Validators) {
	for idx := range validators.Entries {
		var metrics []*sender.Metric
		validator := validators.Entries[idx]

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
		z := sender.NewSender(host, port)

		resp, err := z.Send(packet)
		if err != nil {
			log.Fatalf("Zabbix send failed: %v", err)
		}

		println(cast.ToString(resp))
	}
}
