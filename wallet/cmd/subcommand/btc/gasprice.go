package btc

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	bitcoin "wallet/btc"
	"wallet/db"
)

var gasPriceBTCPrompt = "help info: wallet btc gasprice"

var GasPriceBTCCmd = &cobra.Command{
	Use:   "gasprice",
	Short: "get gas price",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			fmt.Println(gasPriceBTCPrompt)
			return
		}

		gasPrice, err := bitcoin.GetFeeRate()
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(gasPriceBTCPrompt)
			return
		}
		{
			coin_name := "btc"
			sellSql := "update gas_price set gas_price = ? where coin_name = ?"
			_, err := db.Instance().Exec(sellSql, strconv.FormatInt(int64(gasPrice), 10), coin_name)
			if err != nil {
				fmt.Println("update fail, sql:", sellSql)
				return
			}
		}

		fmt.Println("btc gasPrice:", gasPrice)
		return
	},
}
