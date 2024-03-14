package trx

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	"wallet/db"
)

var gasPriceTRXPrompt = "help info: wallet trx gasprice"

var GasPriceETHCmd = &cobra.Command{
	Use:   "gasprice",
	Short: "get gas price",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			fmt.Println(gasPriceTRXPrompt)
			return
		}

		gasPrice := TRX_THERSHOLD_ON_COLLECT_TRC20
		{
			coin_name := "trx"
			sellSql := "update gas_price set gas_price = ? where coin_name = ?"
			_, err := db.Instance().Exec(sellSql, strconv.FormatInt(gasPrice, 10), coin_name)
			if err != nil {
				fmt.Println("update fail, sql:", sellSql)
				return
			}
		}

		fmt.Println("trx gasPrice:", gasPrice)
		return
	},
}
