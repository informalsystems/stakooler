package display

import (
	"fmt"
	"log"
	"strings"
	"time"

	sender "github.com/adubkov/go-zabbix"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"github.com/spf13/cast"
)

func ZbxSendValStats(server string, port int, host string, validators *model.Validators) {
	for _, validator := range validators.Entries {
		var metrics []*sender.Metric

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

func ZbxSendChainDiscovery(config *model.Config) {
	var message []*sender.Metric

	data := []string{"["}
	for idx, chain := range config.Chains.Entries {
		if idx != len(config.Chains.Entries)-1 {
			data = append(data, fmt.Sprintf("{\"{#CHAIN}\":\"%s\"},", chain.ID))
		} else {
			data = append(data, fmt.Sprintf("{\"{#CHAIN}\":\"%s\"}", chain.ID))
		}
	}
	data = append(data, "]")

	var builder strings.Builder
	for _, s := range data {
		_, err := builder.WriteString(s)
		if err != nil {
			log.Fatal(err)
		}
	}

	message = append(message, sender.NewMetric(config.Zabbix.Host, "validator.data.discovery", builder.String(), time.Now().Unix()))
	packet := sender.NewPacket(message)
	z := sender.NewSender(config.Zabbix.Server, config.Zabbix.Port)

	resp, err := z.Send(packet)
	if err != nil {
		log.Fatalf("Zabbix send failed: %v", err)
	}

	fmt.Println(fmt.Sprintf("Chain discovery response: %s", cast.ToString(resp)))
}
