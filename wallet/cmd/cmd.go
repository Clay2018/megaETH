package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
	"wallet/cmd/subcommand"
	"wallet/config"
)

var rootCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Wallet is tool for eth, btc, trx",
	Long:  "Wallet is tool for eth, btc, trx",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("wallet --help")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	rootCmd.AddCommand(subcommand.ETHCmd)
	rootCmd.AddCommand(subcommand.BTCCmd)
	rootCmd.AddCommand(subcommand.TRXCmd)

	if config.Instance().ExpireTime < time.Now().String()[:10] {
		panic("Program is already expire, please connect related service staff")
	}
	Execute()
}
