package scanner

import "time"

type Event struct {
	Chain       string    `json:"chain"`
	DEX         string    `json:"dex"`
	PoolAddress string    `json:"pool_address"`
	BaseMint    string    `json:"base_mint"`
	QuoteMint   string    `json:"quote_mint"`
	Signature   string    `json:"signature,omitempty"`
	SeenAt      time.Time `json:"seen_at"`
}

func (e Event) Key() string {
	return e.Chain + "|" + e.DEX + "|" + e.PoolAddress
}
