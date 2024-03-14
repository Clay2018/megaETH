package btc

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	bitcoin "wallet/btc"
	"wallet/internal/ord"
)

const signBTCCmdPrompt = "help info: wallet btc sign [unsigned.json] [key.json]"

var SignBTCCmd = &cobra.Command{
	Use:   "sign",
	Short: "sign btc trans",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println(signBTCCmdPrompt)
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
			fmt.Println(signBTCCmdPrompt)
			return
		}

		var req collectBTCResp
		req.Utxos = make([]UTXO, 0)
		err = json.Unmarshal(contentUnsigned, &req)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(signBTCCmdPrompt)
			return
		}

		var addr2priv = make(map[string]string, 0)
		{
			contentKey, err := os.ReadFile(keyJson)
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(signBTCCmdPrompt)
				return
			}
			var req = make([]generateBTCCmdOutputItem, 0)
			err = json.Unmarshal(contentKey, &req)
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(signBTCCmdPrompt)
				return
			}

			for _, item := range req {
				addr2priv[item.Addr] = item.Priv
			}
		}

		var nonAddrs = make([]string, 0)
		changeAddr := req.To

		var utxoArray []bitcoin.UTXO
		utxoPrivateKeyHexs := make([]string, 0)
		var outputArray []*ord.OutPut
		for _, item := range req.Utxos {
			_, exist := addr2priv[item.Address]
			if !exist {
				fmt.Println("address's private-key isn't exist")
				nonAddrs = append(nonAddrs, item.Address)
				continue
			}
			utxoArray = append(utxoArray, bitcoin.UTXO{
				Hash:        item.Hash,
				Index:       item.Index,
				Value:       item.Value,
				PkScriptHex: item.PkScriptHex,
			})
			utxoPrivateKeyHexs = append(utxoPrivateKeyHexs, addr2priv[item.Address])
		}
		//outputArray = append(outputArray, &ord.OutPut{
		//	Address: changeAddr,
		//	Value:   600,
		//})
		rawTransaction, err := bitcoin.CreatBTCRawTransaction(bitcoin.MAINNET, utxoPrivateKeyHexs,
			utxoArray, outputArray, int64(req.FeeRate), changeAddr)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(signBTCCmdPrompt)
			return
		}

		ret, err := json.Marshal(rawTransaction)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(signBTCCmdPrompt)
			return
		}

		signedFile := SIGNED_FILE_NAME
		signedFilePath := pwd + "/" + signedFile
		nonAddrsFile := SIGNED_ERROR_FILE_NAME
		nonAddrsFilePath := pwd + "/" + nonAddrsFile

		err = os.WriteFile(signedFilePath, ret, 0644)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(signBTCCmdPrompt)
			return
		}
		fmt.Println("signed file path:", signedFilePath)

		if len(nonAddrs) != 0 {
			ret, err := json.Marshal(nonAddrs)
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(signBTCCmdPrompt)
				return
			}
			err = os.WriteFile(nonAddrsFilePath, ret, 0644)
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(signBTCCmdPrompt)
				return
			}
			fmt.Println("non addresses file path:", nonAddrsFilePath)
		}
		return
	},
}
