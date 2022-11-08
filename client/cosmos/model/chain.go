package model

type Chains struct {
	Entries []Chain
}

type Chain struct {
	ID       string
	LCD      string
	Denom    string
	Exponent int
}
