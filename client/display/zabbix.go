package display

/*
func ZbxValidatorStats(config *model.Config) {
	for _, validator := range config.Validators.Entries {
		var metrics []*sender.Metric

		metrics = append(metrics, sender.NewMetric(validator.Chain.Id, "validator.stats.moniker", validator.Moniker, validator.BlockTime.Unix()))
		metrics = append(metrics, sender.NewMetric(validator.Chain.Id, "validator.stats.valoper", validator.ValoperAddress, validator.BlockTime.Unix()))
		metrics = append(metrics, sender.NewMetric(validator.Chain.Id, "validator.stats.block.height", validator.BlockHeight, validator.BlockTime.Unix()))
		metrics = append(metrics, sender.NewMetric(validator.Chain.Id, "validator.stats.voting.power", fmt.Sprintf("%d", validator.VotingPower), validator.BlockTime.Unix()))
		metrics = append(metrics, sender.NewMetric(validator.Chain.Id, "validator.stats.voting.percent", fmt.Sprintf("%.2f", validator.VotingPercent), validator.BlockTime.Unix()))
		metrics = append(metrics, sender.NewMetric(validator.Chain.Id, "validator.stats.ranking", fmt.Sprintf("%d", validator.Ranking), validator.BlockTime.Unix()))
		metrics = append(metrics, sender.NewMetric(validator.Chain.Id, "validator.stats.commission", fmt.Sprintf("%.2f", validator.Commission), validator.BlockTime.Unix()))
		metrics = append(metrics, sender.NewMetric(validator.Chain.Id, "validator.stats.delegators", validator.NumDelegators, validator.BlockTime.Unix()))
		metrics = append(metrics, sender.NewMetric(validator.Chain.Id, "validator.stats.unbondings", fmt.Sprintf("%d", validator.Unbondings), validator.BlockTime.Unix()))

		fmt.Println(fmt.Sprintf("Validator stats response: %s", SendPacket(metrics, config)))
	}
}*/
/*
func ZbxAccountsDetails(config *model.Config) {
	for _, chain := range config.Chains.Entries {

		for _, account := range config.Accounts.Entries {
			if chain.Id == account.Chain.Id {
				var metrics []*sender.Metric
				for _, token := range account.Tokens {
					metrics = append(metrics, sender.NewMetric(account.Chain.Id, "account.address.["+account.Address+"]", account.Address, account.BlockTime.Unix()))
					metrics = append(metrics, sender.NewMetric(account.Chain.Id, "account.balance.["+account.Address+"]", fmt.Sprintf("%.2f", token.Balance), account.BlockTime.Unix()))
					metrics = append(metrics, sender.NewMetric(account.Chain.Id, "account.height.["+account.Address+"]", account.BlockHeight, account.BlockTime.Unix()))
					metrics = append(metrics, sender.NewMetric(account.Chain.Id, "account.commission.["+account.Address+"]", fmt.Sprintf("%.2f", token.Commission), account.BlockTime.Unix()))
					metrics = append(metrics, sender.NewMetric(account.Chain.Id, "account.denom.["+account.Address+"]", token.DisplayName, account.BlockTime.Unix()))
					metrics = append(metrics, sender.NewMetric(account.Chain.Id, "account.name.["+account.Address+"]", account.Name, account.BlockTime.Unix()))
					metrics = append(metrics, sender.NewMetric(account.Chain.Id, "account.rewards.["+account.Address+"]", fmt.Sprintf("%.2f", token.Reward), account.BlockTime.Unix()))
					metrics = append(metrics, sender.NewMetric(account.Chain.Id, "account.staked.["+account.Address+"]", fmt.Sprintf("%.2f", token.Delegation), account.BlockTime.Unix()))
					metrics = append(metrics, sender.NewMetric(account.Chain.Id, "account.unbonding.["+account.Address+"]", fmt.Sprintf("%.2f", token.Unbonding), account.BlockTime.Unix()))

				}

				if metrics != nil {
					fmt.Println(fmt.Sprintf("Accounts details response: %s", SendPacket(metrics, config)))
				}
			}
		}
	}
}*/
/*
func ZbxSendAccountsDiscovery(config *model.Config) {
	for _, chain := range config.Chains.Entries {
		var message []*sender.Metric

		data := []string{"["}
		for _, account := range config.Accounts.Entries {
			if chain.Id == account.Chain.Id {
				data = append(data, fmt.Sprintf("{\"{#ACCT}\":\"%s\",\"{#ADDR}\":\"%s\"}", account.Name, account.Address))
				data = append(data, ",")
			}
		}
		data = data[:len(data)-1] // Removing the last comma from the list as it will break Zabbix processing
		data = append(data, "]")

		if len(data) > 2 { // Only send data for chains that have accounts configured
			message = append(message, sender.NewMetric(chain.Id, "account.discovery", BuildString(data), time.Now().Unix()))
			fmt.Println(fmt.Sprintf("Account discovery response: %s", SendPacket(message, config)))
		}
	}
}

func ZbxSendChainDiscovery(config *model.Config) {
	var message []*sender.Metric

	data := []string{"["}
	for idx, chain := range config.Chains.Entries {
		data = append(data, fmt.Sprintf("{\"{#CHAIN}\":\"%s\"}", chain.Id))
		if idx != len(config.Chains.Entries)-1 {
			data = append(data, ",")
		}
	}
	data = append(data, "]")

	message = append(message, sender.NewMetric(config.Zabbix.Host, "chain.data.discovery", BuildString(data), time.Now().Unix()))
	fmt.Println(fmt.Sprintf("Chain discovery response: %s", SendPacket(message, config)))
}

func BuildString(data []string) string {
	var builder strings.Builder
	for _, s := range data {
		_, err := builder.WriteString(s)
		if err != nil {
			log.Fatal(err)
		}
	}

	return builder.String()
}

func SendPacket(message []*sender.Metric, config *model.Config) string {
	packet := sender.NewPacket(message)
	z := sender.NewSender(config.Zabbix.Server, config.Zabbix.Port)

	resp, err := z.Send(packet)
	if err != nil {
		log.Fatalf("Zabbix send failed: %v", err)
	}

	return fmt.Sprintf("%s", string(resp))
}
*/
