package btc

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	bitcoin "wallet/btc"
	"wallet/db"
)

const generateBTCCmdPrompt = "help info: wallet btc generate [number]"

type generateBTCCmdOutputItem struct {
	Addr string `json:"address"`
	Pub  string `json:"pub"`
	Priv string `json:"priv"`
}

var GenerateBTCCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate [number] address",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println(generateBTCCmdPrompt)
			return
		}
		numberStr := args[0]
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(generateBTCCmdPrompt)
			return
		}

		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		var resp = make([]generateBTCCmdOutputItem, 0)

		var sqls string
		for i := 0; i < number; i++ {
			var sql = "insert into table_btc_addr(addr, pub, priv) values('"
			priv, pub, addr := bitcoin.CreateAddress(bitcoin.MAINNET)
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(generateBTCCmdPrompt)
				return
			}
			resp = append(resp, generateBTCCmdOutputItem{
				Addr: addr,
				Pub:  pub,
				Priv: priv,
			})
			sql += addr + "', '"
			sql += pub + "', '"
			sql += priv + "')\n"
			sqls += sql
			{
				sellSql := "insert into table_btc_balance(addr, btc_balance, flag) values(?, ?, ?)"
				_, err := db.Instance().Exec(sellSql, addr, "0", 1)
				if err != nil {
					fmt.Println("insert fail, err:", err.Error())
					continue
				}
			}
		}

		keyFilename := KEY_JSON_FILE_NAME
		keyPath := pwd + "/" + keyFilename

		sqlFilename := INSERT_KEY_SQL_FILE_NAME
		sqlPath := pwd + "/" + sqlFilename

		ret, err := json.Marshal(resp)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(generateBTCCmdPrompt)
			return
		}

		err = os.WriteFile(keyPath, ret, 0644)
		err = os.WriteFile(sqlPath, []byte(sqls), 0644)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(generateBTCCmdPrompt)
			return
		}
		fmt.Println("successfully generate addresses")
		fmt.Println("key file path:", keyPath)
		fmt.Println("sql file path:", sqlPath)
		return
	},
}
