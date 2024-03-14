package config

import (
	"time"
)

type DbConfig struct {
	UserName string `json:"userName"`
	Pwd      string `json:"pwd"`
	Database string `json:"database"`
	URL      string `json:"url"`
}

type Config struct {
	ETH_URL          string   `json:"eth_url"`
	ETH_USDT_ADDR    string   `json:"eth_usdt_addr"`
	ETH_USDC_ADDR    string   `json:"eth_usdc_addr"`
	ETH_ADDRS        []string `json:"eth_addrs"`
	BTC_URL          string   `json:"btc_url"`
	BTC_ADDRS        []string `json:"btc_addrs"`
	TRX_URL          string   `json:"trx_url"`
	TRX_URL2         string   `json:"trx_url2"`
	TRX_API_KEY      string   `json:"trx_api_key"`
	TRX_API_KEY2     string   `json:"trx_api_key2"`
	TRX_ADDRS        []string `json:"trx_addrs"`
	TRX_USDT_ADDR    string   `json:"trx_usdt_addr"`
	DBConfig         DbConfig `json:"db_config"`
	TRX_NUMBER       int      `json:"trx_number"`
	ETH_START_NUMBER uint64   `json:"eth_start_number"`
	BTC_START_NUMBER uint64   `json:"btc_start_number"`
	TRX_START_NUMBER uint64   `json:"trx_start_number"`
	ExpireTime       string   `json:"expire_time"`
}

func Instance() *Config {
	//return &Config{
	//	ETH_URL:       "https://eth-sepolia.g.alchemy.com/v2/Ui6GXseTRBnyEfEA_yCioXRmE2A7H2_A",
	//	ETH_USDT_ADDR: "0xaA8E23Fb1079EA71e0a56F48a2aA51851D8433D0",
	//	ETH_USDC_ADDR: "0x94a9D9AC8a22534E3FaCa9F4e7F2E2cf85d5E4C8",
	//	ETH_ADDRS:     []string{"0xD56b55c7ffea0bd2e34044dbeae3d4bd212094d7", "0x6Ae36900bd70E51EdA10529B239E3EfA6708126F", "0x34EfF639A79f3D1EED43E5da23399350CA027eBC"},
	//	BTC_URL:       "https://mempool.space/api",
	//	BTC_ADDRS:     []string{"bc1pz3sxuh75h4zuhmx0fsj9lchph0lwxdr09e9kz7rrftyqh8kjjhsqtzpqcs", "bc1qq79erms5c9k4rujm7w2xd5pkx7u4qfv0cwl4vu"},
	//	TRX_URL:       "https://apilist.tronscanapi.com",
	//	TRX_API_KEY:   "c1e03b44-6625-4689-b0dc-f7ff1ddd96d6",
	//	TRX_URL2:      "https://nile.trongrid.io",
	//	TRX_API_KEY2:  "c1e03b44-6625-4689-b0dc-f7ff1ddd96d6",
	//	//TRX_USDT_ADDR:      "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
	//	TRX_USDT_ADDR: "TXLAQ63Xg1NAzckPwKHvzw7CSEmLMEqcdj",
	//	TRX_ADDRS:     []string{"TPfHmimxS2drgvxzDiAbgFahrcsBoZG9EQ", "TXw2S8gcnUkVn1ay8Jix77wXeVGnprHsn7"},
	//	DBConfig: DbConfig{
	//		UserName: "",
	//		Pwd:      "",
	//		Database: "",
	//		URL:      "",
	//	},
	//}
	return &Config{
		ETH_URL:       "https://eth-mainnet.g.alchemy.com/v2/OHwgAaRdVqiCJBOPj8MDoZ76HB2-vOU8",
		ETH_USDT_ADDR: "0xdac17f958d2ee523a2206206994597c13d831ec7",
		ETH_USDC_ADDR: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		ETH_ADDRS:     []string{},
		BTC_URL:       "https://mempool.space/api",
		BTC_ADDRS:     []string{},
		TRX_URL:       "https://apilist.tronscanapi.com",
		TRX_API_KEY:   "c1e03b44-6625-4689-b0dc-f7ff1ddd96d6",
		//TRX_URL2:      "https://nile.trongrid.io",
		TRX_URL2:         "https://api.trongrid.io",
		TRX_API_KEY2:     "c1e03b44-6625-4689-b0dc-f7ff1ddd96d6",
		TRX_USDT_ADDR:    "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
		TRX_ADDRS:        []string{},
		TRX_NUMBER:       5,
		ETH_START_NUMBER: 19127516,
		TRX_START_NUMBER: 58692143,
		BTC_START_NUMBER: 828255,
		DBConfig: DbConfig{
			UserName: "",
			Pwd:      "",
			Database: "",
			URL:      "",
		},
		ExpireTime: "2025-01-20",
	}
}

const ETH_ADDRESSES = "eth_addresses.json"
const BTC_ADDRESSES = "btc_addresses.json"
const TRX_ADDRESSES = "trx_addresses.json"

func InitConfig() {
	{
		//pwd, err := os.Getwd()
		//if err != nil {
		//	panic(err)
		//}
		//content, err := os.ReadFile(pwd + "/" + ETH_ADDRESSES)
		//if err != nil {
		//	fmt.Println("error:", err.Error())
		//	return
		//}
		//var req = make([]string, 0)
		//err = json.Unmarshal(content, &req)
		//if err != nil {
		//	fmt.Println("error:", err.Error())
		//	return
		//}

	}
	{
		//pwd, err := os.Getwd()
		//if err != nil {
		//	panic(err)
		//}
		//content, err := os.ReadFile(pwd + "/" + BTC_ADDRESSES)
		//if err != nil {
		//	fmt.Println("error:", err.Error())
		//	return
		//}
		//var req = make([]string, 0)
		//err = json.Unmarshal(content, &req)
		//if err != nil {
		//	fmt.Println("error:", err.Error())
		//	return
		//}

	}
	{
		//pwd, err := os.Getwd()
		//if err != nil {
		//	panic(err)
		//}
		//content, err := os.ReadFile(pwd + "/" + TRX_ADDRESSES)
		//if err != nil {
		//	fmt.Println("error:", err.Error())
		//	return
		//}
		//var req = make([]string, 0)
		//err = json.Unmarshal(content, &req)
		//if err != nil {
		//	fmt.Println("error:", err.Error())
		//	return
		//}
	}
	if Instance().ExpireTime < time.Now().String()[:10] {
		panic("Program is already expire, please connect related service staff")
	}
}
