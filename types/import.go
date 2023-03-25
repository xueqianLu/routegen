package types

type DexPair struct {
	Address string `json:"contract"`
	Token0  string `json:"token0"`
	Token1  string `json:"token1"`
}

type DexData struct {
	Name  string    `json:"dex"`
	Fee   string    `json:"fee"`
	Pairs []DexPair `json:"pairs"`
}
