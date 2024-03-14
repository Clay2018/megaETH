package btc

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	bitcoin "wallet/btc"
	"wallet/db"
)

const balanceBTCCmdPrompt = "help info: wallet btc balance [addresses.json]"

var BalanceETHCmd = &cobra.Command{
	Use:   "balance",
	Short: "get addresses balance",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println(balanceBTCCmdPrompt)
			return
		}
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		//filename := args[0]
		//content, err := os.ReadFile(filename)
		//if err != nil {
		//	fmt.Println("error:", err.Error())
		//	fmt.Println(balanceBTCCmdPrompt)
		//	return
		//}

		//var addrs = make([]generateBTCCmdOutputItem, 0)
		//err = json.Unmarshal(content, &addrs)
		//if err != nil {
		//	fmt.Println("error:", err.Error())
		//	fmt.Println(balanceBTCCmdPrompt)
		//	return
		//}
		rows, err := db.Instance().Query("select addr, flag from table_btc_balance where flag=1")
		if err != nil {
			fmt.Println("select sb fail, err:", err.Error())
			return
		}
		var sellOrders = make([]db.TableBTCBalance, 0)
		for rows.Next() {
			var sellOrder db.TableBTCBalance
			err := rows.Scan(&sellOrder.Addr, &sellOrder.Flag)
			if err != nil {
				fmt.Println("select sb fail, err:", err.Error())
				return
			}
			sellOrders = append(sellOrders, sellOrder)
		}

		var sqls string
		for _, item := range sellOrders {
			var sql = "insert into table_btc_balance(addr, btc_balance) values('"
			balance, err := bitcoin.GetBalance(item.Addr)
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(balanceBTCCmdPrompt)
				return
			}
			sql += item.Addr + "', '"
			sql += strconv.FormatInt(balance, 10) + "')\n"
			sqls += sql
			{
				sellSql := "update table_btc_balance set btc_balance = ? where addr = ?"
				_, err := db.Instance().Exec(sellSql, strconv.FormatInt(balance, 10), item.Addr)
				if err != nil {
					fmt.Println("update fail, sql:", sellSql)
					return
				}
			}
		}

		sqlFilename := INSERT_BALANCE_SQL_FILE_NAME
		sqlPath := pwd + "/" + sqlFilename
		err = os.WriteFile(sqlPath, []byte(sqls), 0644)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(balanceBTCCmdPrompt)
			return
		}
		fmt.Println("successfully get balance")
		fmt.Println("sql file path:", sqlPath)
		return
	},
}
