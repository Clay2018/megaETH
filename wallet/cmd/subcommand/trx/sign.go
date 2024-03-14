package trx

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"wallet/trx"
)

const signTRXCmdPrompt = "help info: wallet trx sign [unsigned.json] [key.json]"

var SignTRXCmd = &cobra.Command{
	Use:   "sign",
	Short: "sign trx trans",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println(signTRXCmdPrompt)
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
			fmt.Println(signTRXCmdPrompt)
			return
		}

		var req = make([]UnsignedTransItem, 0)
		err = json.Unmarshal(contentUnsigned, &req)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(signTRXCmdPrompt)
			return
		}

		var addr2priv = make(map[string]string, 0)
		{
			contentKey, err := os.ReadFile(keyJson)
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(signTRXCmdPrompt)
				return
			}
			var req = make([]generateTRXCmdOutputItem, 0)
			err = json.Unmarshal(contentKey, &req)
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(signTRXCmdPrompt)
				return
			}

			for _, item := range req {
				addr2priv[item.Addr] = item.Priv
			}
		}

		var nonAddrs = make([]string, 0)
		var rawTransactions = make([]string, 0)
		for _, item := range req {

			privHex := addr2priv[item.From]
			signedTrans, err := trx.Sign(item.Trans, privHex)
			if err != nil {
				fmt.Println("err:", err.Error())
				nonAddrs = append(nonAddrs, item.From)
				continue
			}

			rawTransactions = append(rawTransactions, signedTrans)
		}

		ret, err := json.Marshal(rawTransactions)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(signTRXCmdPrompt)
			return
		}

		signedFile := SIGNED_FILE_NAME
		signedFilePath := pwd + "/" + signedFile
		nonAddrsFile := SIGNED_ERROR_FILE_NAME
		nonAddrsFilePath := pwd + "/" + nonAddrsFile

		err = os.WriteFile(signedFilePath, ret, 0644)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(signTRXCmdPrompt)
			return
		}
		fmt.Println("signed file path:", signedFilePath)

		if len(nonAddrs) != 0 {
			ret, err := json.Marshal(nonAddrs)
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(signTRXCmdPrompt)
				return
			}
			err = os.WriteFile(nonAddrsFilePath, ret, 0644)
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(signTRXCmdPrompt)
				return
			}
			fmt.Println("non addresses file path:", nonAddrsFilePath)
		}
		return
	},
}
