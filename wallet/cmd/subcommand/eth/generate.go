package eth

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"wallet/db"
	"wallet/eth"
)

const generateETHCmdPrompt = "help info: wallet eth generate [number]"

type generateETHCmdOutputItem struct {
	Addr string `json:"address"`
	Priv string `json:"priv"`
}

var GenerateETHCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate [number] address",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println(generateETHCmdPrompt)
			return
		}
		numberStr := args[0]
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(generateETHCmdPrompt)
			return
		}

		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		var resp = make([]generateETHCmdOutputItem, 0)

		var sqls string
		for i := 0; i < number; i++ {
			var sql = "insert into table_addr(addr, priv) values('"
			addr, priv, err := eth.GenerateAddrAndPriv()
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(generateETHCmdPrompt)
				return
			}
			resp = append(resp, generateETHCmdOutputItem{
				Addr: addr,
				Priv: priv,
			})
			sql += addr + "', '"
			sql += priv + "')\n"
			sqls += sql
			{
				sellSql := "insert into table_eth_balance(addr, eth_balance, usdt_balance, usdc_balance, flag) values(?, ?, ?, ?, ?)"
				_, err := db.Instance().Exec(sellSql, addr, "0", "0", "0", 1)
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
			fmt.Println(generateETHCmdPrompt)
			return
		}

		err = os.WriteFile(keyPath, ret, 0644)
		err = os.WriteFile(sqlPath, []byte(sqls), 0644)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(generateETHCmdPrompt)
			return
		}
		fmt.Println("successfully generate addresses")
		fmt.Println("key file path:", keyPath)
		fmt.Println("sql file path:", sqlPath)
		return
	},
}
