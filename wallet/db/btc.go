package db

type TableBTCDeposit struct {
	Addr   string `json:"addr"`
	TxHash string `json:"tx_hash"`
	Value  int64  `json:"value"`
	Vout   int32  `json:"vout"`
}

type TableBTCBalance struct {
	Addr       string `json:"addr"`
	BtcBalance string `json:"btc_balance"`
	Flag       int    `json:"flag"`
}
