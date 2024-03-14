package btc

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	bitcoin "wallet/btc"
)

const collectBTCCmdPrompt = "help info: wallet btc collect [token] [addresses.json] [gasPrice] [toAddr]"

var CollectBTCCmd = &cobra.Command{
	Use:   "collect",
	Short: "collect btc from addresses",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 4 {
			fmt.Println(collectBTCCmdPrompt)
			return
		}

		filename := args[1]
		gasPriceStr := args[2]
		to := args[3]

		token := args[0]
		switch token {
		case "btc":
			collectBTC(filename, gasPriceStr, to)
		default:
			fmt.Println("not support this token")
			return
		}

	},
}

type AddressFileItem struct {
	Addr string `json:"address"`
	Pub  string `json:"pub"`
}

type UTXO struct {
	Hash        string `json:"hash"`
	Index       uint32 `json:"index"`
	Value       int64  `json:"value"`
	PkScriptHex string `json:"pkScriptHex"`
	Address     string `json:"address"`
}

type collectBTCResp struct {
	FeeRate int32  `json:"feeRate"`
	To      string `json:"to"`
	Utxos   []UTXO `json:"utxos"`
}

func collectBTC(filename string, gasPriceStr string, to string) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectBTCCmdPrompt)
		return
	}
	var addrs = make([]AddressFileItem, 0)
	err = json.Unmarshal(content, &addrs)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectBTCCmdPrompt)
		return
	}

	var resp collectBTCResp
	resp.Utxos = make([]UTXO, 0)
	var nonAddrs = make([]string, 0)

	feeRate, err := strconv.Atoi(gasPriceStr)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectBTCCmdPrompt)
		return
	}
	resp.FeeRate = int32(feeRate)
	resp.To = to

	for _, item := range addrs {
		utxos, err := bitcoin.ListUTXO(bitcoin.MAINNET, item.Pub)
		if err != nil {
			fmt.Println("err:", err.Error())
			nonAddrs = append(nonAddrs, item.Addr)
			continue
		}
		if len(utxos) == 0 {
			fmt.Println("len(utxo) == 0")
			nonAddrs = append(nonAddrs, item.Addr)
			continue
		}
		for _, item2 := range utxos {
			resp.Utxos = append(resp.Utxos, UTXO{
				Hash:        item2.Hash,
				Index:       item2.Index,
				PkScriptHex: item2.PkScriptHex,
				Value:       item2.Value,
				Address:     item.Addr,
			})
		}
	}

	ret, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectBTCCmdPrompt)
		return
	}

	unsignedFile := UNSIGNED_FILE_NAME
	unsignedFilePath := pwd + "/" + unsignedFile
	nonAddrsFile := UNSIGNED_ERROR_FILE_NAME
	nonAddrsFilePath := pwd + "/" + nonAddrsFile

	err = os.WriteFile(unsignedFilePath, ret, 0644)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectBTCCmdPrompt)
		return
	}
	fmt.Println("unsigned file path:", unsignedFilePath)

	if len(nonAddrs) != 0 {
		ret, err = json.Marshal(nonAddrs)
		err = os.WriteFile(nonAddrsFilePath, ret, 0644)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(collectBTCCmdPrompt)
			return
		}
		fmt.Println("non addresses file path:", nonAddrsFilePath)
	}
	return
}
