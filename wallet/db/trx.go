package db

type TableTRXDeposit struct {
	Addr         string `json:"addr"`
	TxHash       string `json:"tx_hash"`
	Value        string `json:"value"`
	ContractAddr string `json:"contract_addr"`
}

type TableTrxBalance struct {
	Addr        string `json:"addr"`
	TrxBalance  string `json:"trx_balance"`
	UsdtBalance string `json:"usdt_balance"`
	Flag        int    `json:"flag"`
}
