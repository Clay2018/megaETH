package trx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"wallet/config"
	"wallet/db"
)

var addrs = make(map[string]struct{}, 0)
var contractAddrs = make([]string, 0)

func TRXInit() {
	var req = make([]string, 0)
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
		req = append(req, sellOrder.Addr)
	}
	config.Instance().TRX_ADDRS = append(config.Instance().TRX_ADDRS, req...)

	for _, item := range req {
		addrs[item] = struct{}{}
	}

	contractAddrs = append(contractAddrs, config.Instance().TRX_USDT_ADDR)
}

// deposit
type TrxTx struct {
	To           string
	Value        int64
	Hash         string
	ContractAddr string
	ContractName string
}

func GetBlockTime(blockNUmber int64) (blockTime int64, err error) {
	client := &http.Client{}
	var data = strings.NewReader(`{"num":` + strconv.FormatInt(blockNUmber, 10) + `}`)
	url := config.Instance().TRX_URL2 + "/wallet/getblockbynum"
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		return 0, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	//fmt.Printf("%s\n", bodyText)

	type TmpS struct {
		RawData struct {
			Timestamp int64 `json:"timestamp"`
		} `json:"raw_data"`
	}
	var tmp TmpS
	err = json.Unmarshal(bodyText, &tmp)
	if err != nil {
		return 0, err
	}
	return tmp.RawData.Timestamp, nil
}

func GetTransaction(blockNumber int64) (transactions []*TrxTx, err error) {
	//https://apilist.tronscanapi.com/api/transfer?sort=-timestamp&count=true&limit=20&start=0&address=TLa2f6VPqDgRE67v1736s7bJ8Ray5wYjU7&filterTokenValue=1

	urlPrefix := config.Instance().TRX_URL + "/api/transfer?sort=-timestamp&count=100&block=" + strconv.FormatInt(blockNumber, 10) + "&toAddress="
	transactions = make([]*TrxTx, 0)
	for addr, _ := range addrs {
		url := urlPrefix + addr

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}
		req.Header.Set(TRON_PRO_API_KEY, config.Instance().TRX_API_KEY)
		resp, err := TRXClient.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		bodyText, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		//fmt.Printf("%s\n", bodyText)

		const ContractRetSuccess = "SUCCESS"
		type TmpS struct {
			Data []struct {
				TransactionHash   string `json:"transactionHash"`
				TransferToAddress string `json:"transferToAddress"`
				Amount            int64  `json:"amount"`
				Confirmed         bool   `json:"confirmed"`
				ContractRet       string `json:"contractRet"`
				TokenInfo         struct {
					TokenId   string `json:"tokenId"`
					TokenName string `json:"tokenName"`
				} `json:"tokenInfo"`
			} `json:"data"`
		}

		var tmp TmpS
		err = json.Unmarshal(bodyText, &tmp)
		if err != nil {
			fmt.Println("err:", err.Error())
			continue
		}

		for _, item := range tmp.Data {
			if !(item.ContractRet == ContractRetSuccess && item.Confirmed) {
				continue
			}
			if item.TransferToAddress != addr {
				continue
			}
			transactions = append(transactions, &TrxTx{
				To:           item.TransferToAddress,
				Value:        item.Amount,
				ContractName: item.TokenInfo.TokenName,
				ContractAddr: item.TokenInfo.TokenId,
				Hash:         item.TransactionHash,
			})
		}
	}
	err = nil
	return
}

func GetTransferLog(fromBlockTime int64, toBlockTime int64) (transactions []*TrxTx, err error) {
	// curl "https://apilist.tronscanapi.com/api/token_trc20/transfers?limit=100&contract_address=TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t&start_timestamp=1704946055000&end_timestamp=1704946065000&relatedAddress=TPfHmimxS2drgvxzDiAbgFahrcsBoZG9EQ"
	urlPrefix := config.Instance().TRX_URL + "/api/token_trc20/transfers?limit=100&start_timestamp=" + strconv.FormatInt(fromBlockTime, 10) + "&end_timestamp=" +
		strconv.FormatInt(toBlockTime, 10) + "&contract_address=" + config.Instance().TRX_USDT_ADDR + "&relatedAddress="
	transactions = make([]*TrxTx, 0)
	for addr, _ := range addrs {
		url := urlPrefix + addr

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}
		req.Header.Set(TRON_PRO_API_KEY, config.Instance().TRX_API_KEY)
		resp, err := TRXClient.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		bodyText, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		//fmt.Printf("%s\n", bodyText)

		const EventTypeTransfer = "Transfer"
		const contractRetSuccess = "SUCCESS"
		type TmpS struct {
			TokenTransfers []struct {
				TransactionId   string `json:"transaction_id"`
				ToAddress       string `json:"to_address"`
				Quant           string `json:"quant"`
				EventType       string `json:"event_type"`
				Confirmed       bool   `json:"confirmed"`
				ContractRet     string `json:"contractRet"`
				FinalResult     string `json:"finalResult"`
				ContractAddress string `json:"contract_address"`
				Revert          bool   `json:"revert"`
			} `json:"token_transfers"`
		}
		var tmp TmpS
		err = json.Unmarshal(bodyText, &tmp)
		if err != nil {
			continue
		}

		for _, item := range tmp.TokenTransfers {
			if !(item.EventType == EventTypeTransfer && !item.Revert && item.Confirmed &&
				item.ContractRet == contractRetSuccess && item.FinalResult == contractRetSuccess) {
				continue
			}
			if item.ToAddress != addr {
				continue
			}
			value, err := strconv.ParseInt(item.Quant, 10, 64)
			if err != nil {
				continue
			}
			transactions = append(transactions, &TrxTx{
				To:           item.ToAddress,
				Value:        value,
				Hash:         item.TransactionId,
				ContractAddr: item.ContractAddress,
			})
		}
	}
	err = nil
	return
}
