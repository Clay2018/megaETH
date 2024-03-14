package eth

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"math/big"
	"os"
	"strconv"
	"wallet/config"
	"wallet/eth"
)

const collectETHCmdPrompt = "help info: wallet eth collect [token] [addresses.json] [gasPrice] [toAddr]"

type CollectETHCmdResp struct {
	ChainId string                  `json:"chainId"`
	Items   []CollectETHCmdRespItem `json:"items"`
}

type CollectETHCmdRespItem struct {
	From     string `json:"from"`
	Nonce    string `json:"nonce"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	To       string `json:"to"`
	Value    string `json:"value"`
	Data     string `json:"data"`
}

var CollectETHCmd = &cobra.Command{
	Use:   "collect",
	Short: "collect eth from addresses",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 4 {
			fmt.Println(collectETHCmdPrompt)
			return
		}

		filename := args[1]
		gasPriceStr := args[2]
		to := args[3]

		token := args[0]
		switch token {
		case "eth":
			collectETH(filename, gasPriceStr, to)
		case "usdt":
			collectERC20(filename, gasPriceStr, to, config.Instance().ETH_USDT_ADDR)
		case "usdc":
			collectERC20(filename, gasPriceStr, to, config.Instance().ETH_USDC_ADDR)
		default:
			fmt.Println("not support this token")
			return
		}

	},
}

func collectETH(filename string, gasPriceStr string, to string) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectETHCmdPrompt)
		return
	}
	var req = make([]string, 0)
	err = json.Unmarshal(content, &req)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectETHCmdPrompt)
		return
	}

	var resp CollectETHCmdResp
	resp.Items = make([]CollectETHCmdRespItem, 0)

	chainId, err := eth.GetChainId()
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectETHCmdPrompt)
		return
	}
	resp.ChainId = chainId.String()

	var nonAddrs = make([]string, 0)
	fee, flag := big.NewInt(1).SetString(gasPriceStr, 10)
	if !flag {
		fmt.Println("error:", err.Error())
		fmt.Println(collectETHCmdPrompt)
		return
	}
	fee = fee.Mul(fee, big.NewInt(int64(eth.GAS_ETH)))
	for _, addr := range req {
		balance, err := eth.GetBalance(addr)
		if err != nil {
			fmt.Println("err:", err.Error())
			nonAddrs = append(nonAddrs, addr)
			continue
		}
		if balance.Cmp(big.NewInt(0)) == 0 {
			fmt.Println("balance is zero")
			nonAddrs = append(nonAddrs, addr)
			continue
		}
		nonce, err := eth.GetAddrNonce(addr)
		if err != nil {
			fmt.Println("err:", err.Error())
			nonAddrs = append(nonAddrs, addr)
			continue
		}

		resp.Items = append(resp.Items, CollectETHCmdRespItem{
			From:     addr,
			Nonce:    strconv.FormatInt(int64(nonce), 10),
			Gas:      strconv.FormatInt(int64(eth.GAS_ETH), 10),
			GasPrice: gasPriceStr,
			To:       to,
			Value:    balance.Sub(balance, fee).String(),
		})
	}

	ret, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectETHCmdPrompt)
		return
	}

	unsignedFile := UNSIGNED_FILE_NAME
	unsignedFilePath := pwd + "/" + unsignedFile
	nonAddrsFile := UNSIGNED_ERROR_FILE_NAME
	nonAddrsFilePath := pwd + "/" + nonAddrsFile

	err = os.WriteFile(unsignedFilePath, ret, 0644)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectETHCmdPrompt)
		return
	}
	fmt.Println("unsigned file path:", unsignedFilePath)

	if len(nonAddrs) != 0 {
		ret, err = json.Marshal(nonAddrs)
		err = os.WriteFile(nonAddrsFilePath, ret, 0644)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(collectETHCmdPrompt)
			return
		}
		fmt.Println("non addresses file path:", nonAddrsFilePath)
	}
	return
}

func collectERC20(filename string, gasPriceStr string, to string, erc20Addr string) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectETHCmdPrompt)
		return
	}
	var req = make([]string, 0)
	err = json.Unmarshal(content, &req)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectETHCmdPrompt)
		return
	}

	var resp CollectETHCmdResp
	resp.Items = make([]CollectETHCmdRespItem, 0)

	chainId, err := eth.GetChainId()
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectETHCmdPrompt)
		return
	}
	resp.ChainId = chainId.String()

	var nonAddrs = make([]string, 0)
	fee, flag := big.NewInt(1).SetString(gasPriceStr, 10)
	if !flag {
		fmt.Println("error:", err.Error())
		fmt.Println(collectETHCmdPrompt)
		return
	}
	fee = fee.Mul(fee, big.NewInt(int64(eth.GAS_USDT)))
	for _, addr := range req {
		balance, err := eth.GetBalance(addr)
		if err != nil {
			fmt.Println("err:", err.Error())
			nonAddrs = append(nonAddrs, addr)
			continue
		}
		if fee.Cmp(balance) > 0 {
			fmt.Println("fee > address's balance", fee, balance)
			nonAddrs = append(nonAddrs, addr)
			continue
		}
		balance20, err := eth.GetERC20Balance(addr, erc20Addr)
		if err != nil {
			fmt.Println("err:", err.Error())
			nonAddrs = append(nonAddrs, addr)
			continue
		}
		if balance20.Cmp(big.NewInt(0)) == 0 {
			fmt.Println("balance is zero.")
			nonAddrs = append(nonAddrs, addr)
			continue
		}
		nonce, err := eth.GetAddrNonce(addr)
		if err != nil {
			fmt.Println("err:", err.Error())
			nonAddrs = append(nonAddrs, addr)
			continue
		}
		data := ""
		{
			USDTBytes, err := os.ReadFile("../eth/abi/ERC20.json")
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(collectETHCmdPrompt)
				return
			}
			contractAbi, err := abi.JSON(bytes.NewReader(USDTBytes))
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(collectETHCmdPrompt)
				return
			}
			dataBytes, err := contractAbi.Pack("transfer", common.HexToAddress(to), balance20)
			if err != nil {
				fmt.Println("error:", err.Error())
				fmt.Println(collectETHCmdPrompt)
				return
			}
			data = hex.EncodeToString(dataBytes)
		}

		resp.Items = append(resp.Items, CollectETHCmdRespItem{
			From:     addr,
			Nonce:    strconv.FormatInt(int64(nonce), 10),
			Gas:      strconv.FormatInt(int64(eth.GAS_USDT), 10),
			GasPrice: gasPriceStr,
			To:       erc20Addr,
			Value:    big.NewInt(0).String(),
			Data:     data,
		})
	}

	ret, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectETHCmdPrompt)
		return
	}

	unsignedFile := UNSIGNED_FILE_NAME
	unsignedFilePath := pwd + "/" + unsignedFile
	nonAddrsFile := UNSIGNED_ERROR_FILE_NAME
	nonAddrsFilePath := pwd + "/" + nonAddrsFile

	err = os.WriteFile(unsignedFilePath, ret, 0644)
	if err != nil {
		fmt.Println("error:", err.Error())
		fmt.Println(collectETHCmdPrompt)
		return
	}
	fmt.Println("unsigned file path:", unsignedFilePath)

	if len(nonAddrs) > 0 {
		ret, err := json.Marshal(nonAddrs)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(collectETHCmdPrompt)
			return
		}
		err = os.WriteFile(nonAddrsFilePath, ret, 0644)
		if err != nil {
			fmt.Println("error:", err.Error())
			fmt.Println(collectETHCmdPrompt)
			return
		}
		fmt.Println("non addresses file path:", nonAddrsFilePath)
	}

	return
}
