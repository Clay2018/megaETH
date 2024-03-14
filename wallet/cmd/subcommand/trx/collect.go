package trx

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"math/big"
	"os"
	"strconv"
	"wallet/config"
	"wallet/trx"
)

const collectTRXCmdPrompt = "help info: wallet trx collect [token] [addresses.json] [gasPrice] [toAddr]"

type UnsignedTransItem struct {
	Trans string `json:"trans"`
	From  string `json:"from"`
}

var CollectTRXCmd = &cobra.Command{
	Use:   "collect",
	Short: "collect trx from addresses",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 4 {
			fmt.Println(collectTRXCmdPrompt)
			return
		}

		filename := args[1]
		gasPriceStr := args[2]
		to := args[3]

		token := args[0]
		switch token {
		case "trx":
			collectTRX(filename, gasPriceStr, to)
		case "usdt":
			collectTRC20(filename, gasPriceStr, to, config.Instance().TRX_USDT_ADDR)
		default:
			fmt.Println("not support this token")
			return
		}

	},
}

func collectTRX(filename string, gasPriceStr string, to string) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectTRXCmdPrompt)
		return
	}
	var req = make([]string, 0)
	err = json.Unmarshal(content, &req)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectTRXCmdPrompt)
		return
	}

	var resp = make([]UnsignedTransItem, 0)

	var nonAddrs = make([]string, 0)
	for _, addr := range req {
		balance, err := trx.GetTrxBalance(addr)
		if err != nil {
			fmt.Println("err:", err.Error())
			nonAddrs = append(nonAddrs, addr)
			return
		}
		if balance <= TRX_THERSHOLD_ON_COLLECT_TRX {
			fmt.Println("balance < minimum balance,", balance, TRX_THERSHOLD_ON_COLLECT_TRC20)
			nonAddrs = append(nonAddrs, addr)
			continue
		}

		amount := balance - TRX_THERSHOLD_ON_COLLECT_TRX
		trans, err := trx.NewTrxTransfer(addr, to, amount)
		if err != nil {
			fmt.Println("err:", err.Error())
			nonAddrs = append(nonAddrs, addr)
			continue
		}

		resp = append(resp, UnsignedTransItem{
			Trans: trans,
			From:  addr,
		})
	}

	ret, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectTRXCmdPrompt)
		return
	}

	unsignedFile := UNSIGNED_FILE_NAME
	unsignedFilePath := pwd + "/" + unsignedFile
	nonAddrsFile := UNSIGNED_ERROR_FILE_NAME
	nonAddrsFilePath := pwd + "/" + nonAddrsFile

	err = os.WriteFile(unsignedFilePath, ret, 0644)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectTRXCmdPrompt)
		return
	}
	fmt.Println("unsigned file path:", unsignedFilePath)

	if len(nonAddrs) != 0 {
		ret, err = json.Marshal(nonAddrs)
		err = os.WriteFile(nonAddrsFilePath, ret, 0644)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(collectTRXCmdPrompt)
			return
		}
		fmt.Println("non addresses file path:", nonAddrsFilePath)
	}
	return
}

func collectTRC20(filename string, gasPriceStr string, to string, trc20Addr string) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectTRXCmdPrompt)
		return
	}
	var req = make([]string, 0)
	err = json.Unmarshal(content, &req)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectTRXCmdPrompt)
		return
	}

	var resp = make([]UnsignedTransItem, 0)
	var nonAddrs = make([]string, 0)

	var thershold = TRX_THERSHOLD_ON_COLLECT_TRC20
	if gasPriceStr != "" {
		tmp, err := strconv.ParseInt(gasPriceStr, 10, 64)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(collectTRXCmdPrompt)
			return
		}
		thershold = tmp
	}

	for _, addr := range req {
		balance, err := trx.GetTrxBalance(addr)
		if err != nil {
			fmt.Println("err:", err.Error(), ", address:", addr)
		}
		if err != nil || balance < thershold {
			fmt.Println("address:", addr, ",trx balance less than ", thershold, ", balance:", balance)
			nonAddrs = append(nonAddrs, addr)
			continue
		}
		var balanceUSDT = big.NewInt(1)
		var balanceStr string
		{
			balances, isActive, err := trx.GetBalance(addr)
			if err != nil || !isActive {
				fmt.Println("addr is non active or network error")
				nonAddrs = append(nonAddrs, addr)
				continue
			}
			for _, item := range balances {
				if item.TokenId == trc20Addr {
					balanceStr = item.Balance
					break
				}
			}
		}
		balanceUSDT, flag := balanceUSDT.SetString(balanceStr, 10)
		if !flag || balanceUSDT.Cmp(big.NewInt(0)) == 0 {
			fmt.Println("balanceStr is invalid,", balanceStr)
			nonAddrs = append(nonAddrs, addr)
			continue
		}

		trans, err := trx.NewTrx20Transfer(addr, to, balanceUSDT, trc20Addr)
		if err != nil {
			fmt.Println("err:", err.Error())
			nonAddrs = append(nonAddrs, addr)
			continue
		}
		resp = append(resp, UnsignedTransItem{
			Trans: trans,
			From:  addr,
		})
	}

	ret, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectTRXCmdPrompt)
		return
	}

	unsignedFile := UNSIGNED_FILE_NAME
	unsignedFilePath := pwd + "/" + unsignedFile
	nonAddrsFile := UNSIGNED_ERROR_FILE_NAME
	nonAddrsFilePath := pwd + "/" + nonAddrsFile

	err = os.WriteFile(unsignedFilePath, ret, 0644)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectTRXCmdPrompt)
		return
	}
	fmt.Println("unsigned file path:", unsignedFilePath)

	if len(nonAddrs) > 0 {
		ret, err := json.Marshal(nonAddrs)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(collectTRXCmdPrompt)
			return
		}
		err = os.WriteFile(nonAddrsFilePath, ret, 0644)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(collectTRXCmdPrompt)
			return
		}
		fmt.Println("non addresses file path:", nonAddrsFilePath)
	}

	return
}
