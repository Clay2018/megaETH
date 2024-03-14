package trx

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"wallet/db"
	"wallet/trx"
)

const generateTRXCmdPrompt = "help info: wallet trx generate [number]"

type generateTRXCmdOutputItem struct {
	Addr string `json:"address"`
	Pub  string `json:"pub"`
	Priv string `json:"priv"`
}

var GenerateTRXCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate [number] address",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println(generateTRXCmdPrompt)
			return
		}
		numberStr := args[0]
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(generateTRXCmdPrompt)
			return
		}

		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		var resp = make([]generateTRXCmdOutputItem, 0)

		var sqls string
		for i := 0; i < number; i++ {
			var sql = "insert into table_addr(addr, pub, priv) values('"
			addr, pub, priv, err := trx.GenerateAddrAndPriv2()
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(generateTRXCmdPrompt)
				return
			}
			resp = append(resp, generateTRXCmdOutputItem{
				Addr: addr,
				Priv: priv,
				Pub:  pub,
			})
			sql += addr + "', '"
			sql += pub + "', '"
			sql += priv + "')\n"
			sqls += sql
			{
				sellSql := "insert into table_trx_balance(addr, trx_balance, usdt_balance, flag) values(?, ?, ?, ?)"
				_, err := db.Instance().Exec(sellSql, addr, "0", "0", 1)
				if err != nil {
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
			fmt.Println(generateTRXCmdPrompt)
			return
		}

		err = os.WriteFile(keyPath, ret, 0644)
		err = os.WriteFile(sqlPath, []byte(sqls), 0644)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(generateTRXCmdPrompt)
			return
		}
		fmt.Println("successfully generate addresses")
		fmt.Println("key file path:", keyPath)
		fmt.Println("sql file path:", sqlPath)
		return
	},
}
