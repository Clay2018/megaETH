package eth

import (
	"fmt"
	"github.com/spf13/cobra"
	"wallet/db"
	"wallet/eth"
)

var gasPriceETHPrompt = "help info: wallet eth gasprice"

var GasPriceETHCmd = &cobra.Command{
	Use:   "gasprice",
	Short: "get gas price",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			fmt.Println(gasPriceETHPrompt)
			return
		}

		gasPrice, err := eth.GetGasPrice()
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(gasPriceETHPrompt)
			return
		}
		{
			coin_name := "eth"
			sellSql := "update gas_price set gas_price = ? where coin_name = ?"
			_, err := db.Instance().Exec(sellSql, gasPrice.String(), coin_name)
			if err != nil {
				fmt.Println("update fail, sql:", sellSql)
				return
			}
		}

		fmt.Println("eth gasPrice:", gasPrice)
		return
	},
}
