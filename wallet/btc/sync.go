package bitcoin

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"wallet/config"
	"wallet/db"
	"wallet/pkg/btcapi/mempool"
)

var BTCClient *http.Client = &http.Client{}

var addrs = make(map[string]struct{}, 0)

func BTCSyncInit() {
	var req = make([]string, 0)
	rows, err := db.Instance().Query("select addr, flag from table_btc_balance where flag=1")
	if err != nil {
		fmt.Println("select sb fail, err:", err.Error())
		return
	}
	var sellOrders = make([]db.TableBTCBalance, 0)
	for rows.Next() {
		var sellOrder db.TableBTCBalance
		err := rows.Scan(&sellOrder.Addr, &sellOrder.Flag)
		if err != nil {
			fmt.Println("select sb fail, err:", err.Error())
			return
		}
		sellOrders = append(sellOrders, sellOrder)
		req = append(req, sellOrder.Addr)
	}
	config.Instance().BTC_ADDRS = append(config.Instance().BTC_ADDRS, req...)

	for _, item := range req {
		addrs[item] = struct{}{}
	}
}

func GetLatestBlock() (blockNumber uint64, err error) {
	url := config.Instance().BTC_URL + "/blocks/tip/height"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	resp, err := BTCClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var blockNumber2 int
	blockNumber2, err = strconv.Atoi(string(bodyText))
	blockNumber = uint64(blockNumber2)
	return
}

func GetBlockHash(blockNumber uint64) (string, error) {
	url := config.Instance().BTC_URL + "/block-height/" + strconv.Itoa(int(blockNumber))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	resp, err := BTCClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	blockHash := string(bodyText)
	return blockHash, nil
}

type BTCTx struct {
	To     string
	Value  int64
	TxHash string
	Vout   int32
}

func GetTransaction(blockNumber uint64) ([]*BTCTx, error) {
	//"https://mempool.space/api/block/000000000000000015dc777b3ff2611091336355d3f0ee9766a2cf3be8e4b1ce/txs"
	blockHash, err := GetBlockHash(blockNumber)
	if err != nil {
		return nil, err
	}
	url := config.Instance().BTC_URL + "/block/" + blockHash + "/txs"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := BTCClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("%s\n", string(bodyText))
	type TmpS struct {
		Txid string `json:"txid"`
		Vout []struct {
			ScriptPubkeyAddress string `json:"scriptpubkey_address"`
			Value               int64  `json:"value"`
		} `json:"vout"`
		Status struct {
			Confirmed bool `json:"confirmed"`
		} `json:"status"`
	}

	var tmp = make([]TmpS, 0)
	err = json.Unmarshal(bodyText, &tmp)
	if err != nil {
		return nil, err
	}

	var btcTxs = make([]*BTCTx, 0)
	for _, item := range tmp {
		for i, item2 := range item.Vout {
			_, exist := addrs[item2.ScriptPubkeyAddress]
			if item2.ScriptPubkeyAddress == "bc1pctg36zxy5r3h0j2w7ke8ssnqdrlkt7nf908quh9z3zxrjwxatr9snk7gmh" {
				fmt.Println(item2, exist)
			}
			if !exist {
				continue
			}
			btcTxs = append(btcTxs, &BTCTx{
				To:     item2.ScriptPubkeyAddress,
				Value:  item2.Value,
				Vout:   int32(i),
				TxHash: item.Txid,
			})
		}
	}

	return btcTxs, nil
}

func SendRawTrans(rawTrans string) (hash string, err error) {
	//curl -X POST -sSLd "0200000001fd5b5fcd1cb066c27cfc9fda5428b9be850b81ac440ea51f1ddba2f987189ac1010000008a4730440220686a40e9d2dbffeab4ca1ff66341d06a17806767f12a1fc4f55740a7af24c6b5022049dd3c9a85ac6c51fecd5f4baff7782a518781bbdd94453c8383755e24ba755c01410436d554adf4a3eb03a317c77aa4020a7bba62999df633bba0ea8f83f48b9e01b0861d3b3c796840f982ee6b14c3c4b7ad04fcfcc3774f81bff9aaf52a15751fedfdffffff02416c00000000000017a914bc791b2afdfe1e1b5650864a9297b20d74c61f4787d71d0000000000001976a9140a59837ccd4df25adc31cdad39be6a8d97557ed688ac00000000" "https://mempool.space/api/tx"
	url := config.Instance().BTC_URL + "/tx"
	var data = strings.NewReader(rawTrans)
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := BTCClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(bodyText), nil
}

func GetFeeRate() (feeRate int32, err error) {
	//{"fastestFee":34,"halfHourFee":32,"hourFee":30,"economyFee":30,"minimumFee":23}
	url := config.Instance().BTC_URL + "/v1/fees/recommended"
	//fmt.Println("url:", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	resp, err := BTCClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	type TmpS struct {
		EconomyFee int32 `json:"fastestFee"`
	}
	//fmt.Printf("%s\n", bodyText)

	var tmp TmpS
	err = json.Unmarshal(bodyText, &tmp)
	if err != nil {
		return 0, err
	}
	return tmp.EconomyFee, nil
}

func GetBalance(address string) (balance int64, err error) {
	url := config.Instance().BTC_URL + "/address/" + address
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	resp, err := BTCClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	type TmpS struct {
		Address    string `json:""`
		ChainStats struct {
			FundedTxoSum int64 `json:"funded_txo_sum"`
			SpentTxoSum  int64 `json:"spent_txo_sum"`
		} `json:"chain_stats"`
	}

	var tmp TmpS
	err = json.Unmarshal(bodyText, &tmp)
	if err != nil {
		return 0, err
	}

	return tmp.ChainStats.FundedTxoSum - tmp.ChainStats.SpentTxoSum, nil
}

func ListUTXO(network string, utxoPublicKeyHex string) (utxoArray []UTXO, err error) {
	var netParams *chaincfg.Params
	switch network {
	case MAINNET:
		netParams = &chaincfg.MainNetParams
	case TESTNET:
		netParams = &chaincfg.TestNet3Params
	}
	btcApiClient := mempool.NewClient(netParams)

	pubKey, _ := hex.DecodeString(utxoPublicKeyHex)
	utxoTaprootAddress, err := btcutil.NewAddressTaproot(pubKey, netParams)
	if err != nil {
		return nil, err
	}

	unspentList, err := btcApiClient.ListUnspent(utxoTaprootAddress)

	if err != nil {
		return nil, err
	}
	utxoArray = make([]UTXO, 0)

	for i := 0; i < len(unspentList); i++ {
		tmp := UTXO{
			Hash:        unspentList[i].Outpoint.Hash.String(),
			Index:       unspentList[i].Outpoint.Index,
			Value:       unspentList[i].Output.Value,
			PkScriptHex: "5120" + hex.EncodeToString(unspentList[i].Output.PkScript),
		}
		utxoArray = append(utxoArray, tmp)
	}
	return
}
