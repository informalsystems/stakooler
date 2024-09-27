package model

type Chains struct {
	Entries []Chain
}

type Chain struct {
	Name         string
	Id           string
	RestEndpoint string
	Bech32Prefix string
	Accounts     []Account
	Denom        string
	Exponent     int
	AssetList    *AssetList
}
