package subcommand

import (
	"fmt"
	"github.com/spf13/cobra"
	"wallet/cmd/subcommand/btc"
)

const btcCmdPrompt = "help info: wallet btc [subcommand]"

var BTCCmd = &cobra.Command{
	Use:   "btc",
	Short: "btc tools",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(btcCmdPrompt)
		return
	},
}

func init() {
	BTCCmd.AddCommand(btc.GenerateBTCCmd)
	BTCCmd.AddCommand(btc.BalanceETHCmd)
	BTCCmd.AddCommand(btc.GasPriceBTCCmd)
	BTCCmd.AddCommand(btc.CollectBTCCmd)
	BTCCmd.AddCommand(btc.SignBTCCmd)
}
