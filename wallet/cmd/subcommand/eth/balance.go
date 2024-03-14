package eth

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"wallet/config"
	"wallet/db"
	"wallet/eth"
)

const balanceETHCmdPrompt = "help info: wallet eth balance [addresses.json]"

var BalanceETHCmd = &cobra.Command{
	Use:   "balance",
	Short: "get addresses balance",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println(balanceETHCmdPrompt)
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
		//	fmt.Println(balanceETHCmdPrompt)
		//	return
		//}

		//var addrs = make([]string, 0)
		//err = json.Unmarshal(content, &addrs)
		//if err != nil {
		//	fmt.Println("error:", err.Error())
		//	fmt.Println(balanceETHCmdPrompt)
		//	return
		//}

		rows, err := db.Instance().Query("select addr, flag from table_eth_balance where flag=1")
		if err != nil {
			fmt.Println("select sb fail, err:", err.Error())
			return
		}
		var sellOrders = make([]db.TableBalance, 0)
		for rows.Next() {
			var sellOrder db.TableBalance
			err := rows.Scan(&sellOrder.Addr, &sellOrder.Flag)
			if err != nil {
				fmt.Println("select sb fail, err:", err.Error())
				return
			}
			sellOrders = append(sellOrders, sellOrder)
		}

		var sqls string
		for _, item := range sellOrders {
			var sql = "insert into table_eth_balance(addr, eth_balance, usdt_balance, usdc_balance) values('"
			balance, err := eth.GetBalance(item.Addr)
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(balanceETHCmdPrompt)
				return
			}
			balanceUSDT, err := eth.GetERC20Balance(item.Addr, config.Instance().ETH_USDT_ADDR)
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(balanceETHCmdPrompt)
				return
			}
			balanceUSDC, err := eth.GetERC20Balance(item.Addr, config.Instance().ETH_USDC_ADDR)
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(balanceETHCmdPrompt)
				return
			}
			sql += item.Addr + "', '"
			sql += balance.String() + "','"
			sql += balanceUSDT.String() + "','"
			sql += balanceUSDC.String() + "')\n"
			sqls += sql
			{
				sellSql := "update table_eth_balance set eth_balance = ?, usdt_balance = ?, usdc_balance = ? where addr = ?"
				_, err := db.Instance().Exec(sellSql, balance.String(), balanceUSDT.String(), balanceUSDC.String(), item.Addr)
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
			fmt.Println(balanceETHCmdPrompt)
			return
		}
		fmt.Println("successfully get balance")
		fmt.Println("sql file path:", sqlPath)
		return
	},
}
