package db

type TableETHDeposit struct {
	Addr         string `json:"addr"`
	TxHash       string `json:"tx_hash"`
	Value        string `json:"value"`
	ContractAddr string `json:"contract_addr"`
}

type TableBalance struct {
	Addr        string `json:"addr"`
	EthBalance  string `json:"eth_balance"`
	UsdtBalance string `json:"usdt_balance"`
	UsdcBalance string `json:"usdc_balance"`
	Flag        int    `json:"flag"`
}
