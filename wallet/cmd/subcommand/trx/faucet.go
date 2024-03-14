package trx

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"wallet/trx"
)

const faucetTRXCmdPrompt = "help info: wallet trx faucet [token] [addresses.json] [gasPrice] [fromAddr]"

var FaucetETHCmd = &cobra.Command{
	Use:   "faucet",
	Short: "faucet trx to addresses",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 4 {
			fmt.Println(faucetTRXCmdPrompt)
			return
		}

		token := args[0]
		filename := args[1]
		gasPriceStr := args[2]
		from := args[3]

		switch token {
		case "trx":
			faucetTRX(filename, gasPriceStr, from)
		default:
			fmt.Println("not support this token")
			return
		}
	},
}

func faucetTRX(filename string, gasPriceStr string, from string) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(faucetTRXCmdPrompt)
		return
	}
	var req = make([]string, 0)
	err = json.Unmarshal(content, &req)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(faucetTRXCmdPrompt)
		return
	}

	gasPrice, err := strconv.ParseInt(gasPriceStr, 10, 64)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(faucetTRXCmdPrompt)
		return
	}

	var resp = make([]UnsignedTransItem, 0)
	var nonAddrs = make([]string, 0)
	for _, addr := range req {
		balance, err := trx.GetTrxBalance(addr)
		if err != nil {
			fmt.Println("err:", err.Error())
			nonAddrs = append(nonAddrs, addr)
			continue
		}
		if balance >= gasPrice {
			fmt.Println("balance is invalid", balance, gasPrice)
			nonAddrs = append(nonAddrs, addr)
			continue
		}

		amount := gasPrice - balance
		trans, err := trx.NewTrxTransfer(from, addr, amount)
		if err != nil {
			fmt.Println("err:", err.Error())
			nonAddrs = append(nonAddrs, addr)
			continue
		}
		resp = append(resp, UnsignedTransItem{
			Trans: trans,
			From:  from,
		})
	}

	ret, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(faucetTRXCmdPrompt)
		return
	}

	unsignedFile := UNSIGNED_FILE_NAME
	unsignedFilePath := pwd + "/" + unsignedFile
	nonAddrsFile := UNSIGNED_ERROR_FILE_NAME
	nonAddrsFilePath := pwd + "/" + nonAddrsFile

	err = os.WriteFile(unsignedFilePath, ret, 0644)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(faucetTRXCmdPrompt)
		return
	}
	fmt.Println("unsigned file path:", unsignedFilePath)

	if len(nonAddrs) != 0 {
		ret, err = json.Marshal(nonAddrs)
		err = os.WriteFile(nonAddrsFilePath, ret, 0644)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(faucetTRXCmdPrompt)
			return
		}
		fmt.Println("non addresses file path:", nonAddrsFilePath)
	}
	return
}
