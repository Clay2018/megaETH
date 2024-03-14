package subcommand

import (
	"fmt"
	"github.com/spf13/cobra"
	"wallet/cmd/subcommand/eth"
)

const ethCmdPrompt = "help info: wallet eth [subcommand]"

var ETHCmd = &cobra.Command{
	Use:   "eth",
	Short: "eth tools",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(ethCmdPrompt)
		return
	},
}

func init() {
	ETHCmd.AddCommand(eth.GenerateETHCmd)
	ETHCmd.AddCommand(eth.BalanceETHCmd)
	ETHCmd.AddCommand(eth.GasPriceETHCmd)
	ETHCmd.AddCommand(eth.CollectETHCmd)
	ETHCmd.AddCommand(eth.FaucetETHCmd)
	ETHCmd.AddCommand(eth.SignETHCmd)
}
