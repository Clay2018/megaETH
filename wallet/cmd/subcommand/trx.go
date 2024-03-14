package subcommand

import (
	"fmt"
	"github.com/spf13/cobra"
	"wallet/cmd/subcommand/trx"
)

const trxCmdPrompt = "help info: wallet trx [subcommand]"

var TRXCmd = &cobra.Command{
	Use:   "trx",
	Short: "trx tools",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(trxCmdPrompt)
		return
	},
}

func init() {
	TRXCmd.AddCommand(trx.GenerateTRXCmd)

	TRXCmd.AddCommand(trx.BalanceTRXCmd)
	TRXCmd.AddCommand(trx.GasPriceETHCmd)
	TRXCmd.AddCommand(trx.CollectTRXCmd)
	TRXCmd.AddCommand(trx.FaucetETHCmd)
	TRXCmd.AddCommand(trx.SignTRXCmd)
}
