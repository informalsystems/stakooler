package model

type RawAccountData struct {
	Accounts []struct {
		Name      string `json:"name"`
		Address   string `json:"address"`
		Reporting bool   `json:"reporting"`
	} `json:"accounts"`

	Chains []struct {
		Name     string   `json:"name"`
		Id       string   `json:"id"`
		Rest     string   `json:"rest"`
		Accounts []string `json:"accounts"`
	} `json:"chains"`
}
