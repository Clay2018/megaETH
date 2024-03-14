package eth

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"math/big"
	"os"
	"strconv"
	"wallet/eth"
)

const faucetETHCmdPrompt = "help info: wallet eth faucet [token] [addresses.json] [gasPrice] [fromAddr]"

type FaucetETHCmdResp struct {
	ChainId string                  `json:"chainId"`
	Items   []CollectETHCmdRespItem `json:"items"`
}

type FaucetETHCmdRespItem struct {
	From     string `json:"from"`
	Nonce    string `json:"nonce"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	To       string `json:"to"`
	Value    string `json:"value"`
}

var FaucetETHCmd = &cobra.Command{
	Use:   "faucet",
	Short: "faucet eth to addresses",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 4 {
			fmt.Println(faucetETHCmdPrompt)
			return
		}

		token := args[0]
		filename := args[1]
		gasPriceStr := args[2]
		from := args[3]

		switch token {
		case "eth":
			faucetETH(filename, gasPriceStr, from)
		default:
			fmt.Println("not support this token")
			return
		}
	},
}

func faucetETH(filename string, gasPriceStr string, from string) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(faucetETHCmdPrompt)
		return
	}
	var req = make([]string, 0)
	err = json.Unmarshal(content, &req)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(faucetETHCmdPrompt)
		return
	}

	var resp CollectETHCmdResp
	resp.Items = make([]CollectETHCmdRespItem, 0)

	chainId, err := eth.GetChainId()
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(faucetETHCmdPrompt)
		return
	}
	resp.ChainId = chainId.String()

	var nonAddrs = make([]string, 0)
	fee, flag := big.NewInt(1).SetString(gasPriceStr, 10)
	if !flag {
		fmt.Println("error:", err.Error())
		fmt.Println(faucetETHCmdPrompt)
		return
	}
	fee = fee.Mul(fee, big.NewInt(int64(eth.GAS_ETH)))
	value := fee
	nonce, err := eth.GetAddrNonce(from)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(faucetETHCmdPrompt)
		return
	}
	for i, addr := range req {
		resp.Items = append(resp.Items, CollectETHCmdRespItem{
			From:     from,
			Nonce:    strconv.FormatInt(int64(nonce+uint64(i)), 10),
			Gas:      strconv.FormatInt(int64(eth.GAS_ETH), 10),
			GasPrice: gasPriceStr,
			To:       addr,
			Value:    value.String(),
		})
	}

	ret, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(faucetETHCmdPrompt)
		return
	}

	unsignedFile := UNSIGNED_FILE_NAME
	unsignedFilePath := pwd + "/" + unsignedFile
	nonAddrsFile := UNSIGNED_ERROR_FILE_NAME
	nonAddrsFilePath := pwd + "/" + nonAddrsFile

	err = os.WriteFile(unsignedFilePath, ret, 0644)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(faucetETHCmdPrompt)
		return
	}
	fmt.Println("unsigned file path:", unsignedFilePath)

	if len(nonAddrs) != 0 {
		ret, err = json.Marshal(nonAddrs)
		err = os.WriteFile(nonAddrsFilePath, ret, 0644)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(faucetETHCmdPrompt)
			return
		}
		fmt.Println("non addresses file path:", nonAddrsFilePath)
	}
	return
}
