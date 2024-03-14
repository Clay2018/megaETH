package trx

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"wallet/config"
	"wallet/db"
	"wallet/trx"
)

const balanceTRXCmdPrompt = "help info: wallet trx balance [addresses.json]"

var BalanceTRXCmd = &cobra.Command{
	Use:   "balance",
	Short: "get addresses balance",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println(balanceTRXCmdPrompt)
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
		//	fmt.Println(balanceTRXCmdPrompt)
		//	return
		//}

		//var addrs = make([]string, 0)
		//err = json.Unmarshal(content, &addrs)
		//if err != nil {
		//	fmt.Println("error:", err.Error())
		//	fmt.Println(balanceTRXCmdPrompt)
		//	return
		//}
		rows, err := db.Instance().Query("select addr, flag from table_trx_balance where flag=1")
		if err != nil {
			fmt.Println("select sb fail, err:", err.Error())
			return
		}
		var sellOrders = make([]db.TableTrxBalance, 0)
		for rows.Next() {
			var sellOrder db.TableTrxBalance
			err := rows.Scan(&sellOrder.Addr, &sellOrder.Flag)
			if err != nil {
				fmt.Println("select sb fail, err:", err.Error())
				return
			}
			sellOrders = append(sellOrders, sellOrder)
		}

		var sqls string
		for _, item := range sellOrders {
			var sql = "insert into table_trx_balance(addr, trx_balance, usdt_balance) values('"
			balances, _, err := trx.GetBalance(item.Addr)
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(balanceTRXCmdPrompt)
				return
			}
			var trxBalance string
			var usdtBalance string
			for _, item := range balances {
				switch item.TokenId {
				case TRX_TOKEN_ID:
					trxBalance = item.Balance
				case config.Instance().TRX_USDT_ADDR:
					usdtBalance = item.Balance
				}
			}

			sql += item.Addr + "', '"
			sql += trxBalance + "','"
			sql += usdtBalance + "')\n"
			sqls += sql

			{
				sellSql := "update table_trx_balance set trx_balance = ?, usdt_balance = ? where addr = ?"
				_, err := db.Instance().Exec(sellSql, trxBalance, usdtBalance, item.Addr)
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
			fmt.Println(balanceTRXCmdPrompt)
			return
		}
		fmt.Println("successfully get balance")
		fmt.Println("sql file path:", sqlPath)
		return
	},
}
