package eth

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
	"math/big"
	"os"
	"strconv"
	"wallet/eth"
)

const signETHCmdPrompt = "help info: wallet eth sign [unsigned.json] [key.json]"

var SignETHCmd = &cobra.Command{
	Use:   "sign",
	Short: "sign eth trans",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println(signETHCmdPrompt)
			return
		}
		unsignedFile := args[0]
		keyJson := args[1]
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		contentUnsigned, err := os.ReadFile(unsignedFile)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(signETHCmdPrompt)
			return
		}

		var req CollectETHCmdResp
		req.Items = make([]CollectETHCmdRespItem, 0)
		err = json.Unmarshal(contentUnsigned, &req)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(signETHCmdPrompt)
			return
		}

		var addr2priv = make(map[string]string, 0)
		{
			contentKey, err := os.ReadFile(keyJson)
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(signETHCmdPrompt)
				return
			}
			var req = make([]generateETHCmdOutputItem, 0)
			err = json.Unmarshal(contentKey, &req)
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(signETHCmdPrompt)
				return
			}

			for _, item := range req {
				addr2priv[item.Addr] = item.Priv
			}
		}

		var nonAddrs = make([]string, 0)
		var rawTransactions = make([]string, 0)
		chainId, flag := big.NewInt(1).SetString(req.ChainId, 10)
		if !flag {
			fmt.Println("chain id is invalid")
			fmt.Println(signETHCmdPrompt)
			return
		}
		for _, item := range req.Items {
			from := common.HexToAddress(item.From)
			_, exist := addr2priv[from.String()]
			if !exist {
				fmt.Println("addr's private-key is not exist, addr:", from.String())
				nonAddrs = append(nonAddrs, item.From)
				return
			}
			nonce, err := strconv.ParseInt(item.Nonce, 10, 64)
			if err != nil {
				fmt.Println("err:", err.Error())
				nonAddrs = append(nonAddrs, item.From)
				continue
			}
			gas, err := strconv.ParseInt(item.Gas, 10, 64)
			if err != nil {
				fmt.Println("err:", err.Error())
				nonAddrs = append(nonAddrs, item.From)
				continue
			}
			gasPrice, flag := big.NewInt(1).SetString(item.GasPrice, 10)
			if !flag {
				fmt.Println("err: gasprice is invalid", item.GasPrice)
				nonAddrs = append(nonAddrs, item.From)
				continue
			}
			value, flag := big.NewInt(1).SetString(item.Value, 10)
			if !flag {
				fmt.Println("err: value is invalid", item.Value)
				nonAddrs = append(nonAddrs, item.From)
				continue
			}
			to := common.HexToAddress(item.To)
			tx := &types.LegacyTx{
				Nonce:    uint64(nonce),
				Gas:      uint64(gas),
				GasPrice: gasPrice,
				To:       &to,
				Value:    value,
				Data:     common.Hex2Bytes(item.Data),
			}
			rawTrans, _, err := eth.SignTrans(tx, chainId, addr2priv[from.String()])
			if err != nil {
				fmt.Println("err:", err.Error())
				nonAddrs = append(nonAddrs, item.From)
				continue
			}
			rawTransactions = append(rawTransactions, rawTrans)
		}

		ret, err := json.Marshal(rawTransactions)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(signETHCmdPrompt)
			return
		}

		signedFile := SIGNED_FILE_NAME
		signedFilePath := pwd + "/" + signedFile
		nonAddrsFile := SIGNED_ERROR_FILE_NAME
		nonAddrsFilePath := pwd + "/" + nonAddrsFile

		err = os.WriteFile(signedFilePath, ret, 0644)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(signETHCmdPrompt)
			return
		}
		fmt.Println("signed file path:", signedFilePath)

		if len(nonAddrs) != 0 {
			ret, err := json.Marshal(nonAddrs)
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(signETHCmdPrompt)
				return
			}
			err = os.WriteFile(nonAddrsFilePath, ret, 0644)
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(signETHCmdPrompt)
				return
			}
			fmt.Println("non addresses file path:", nonAddrsFilePath)
		}
		return
	},
}
