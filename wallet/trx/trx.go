package trx

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/okx/go-wallet-sdk/coins/tron"
	"github.com/okx/go-wallet-sdk/coins/tron/pb"
	"github.com/okx/go-wallet-sdk/coins/tron/token"
	"github.com/pkg/errors"
	"github.com/status-im/keycard-go/hexutils"
	"io"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"
	"wallet/config"
)

const GAS_USDT = int64(30000000)

const TRON_PRO_API_KEY = "TRON-PRO-API-KEY"

var TRXClient = &http.Client{}

func init() {
	var err error

	if err != nil {
		panic(err)
	}
}

func GetLatestBlock() (blockNumber int64, blockHash string, blockTime int64, err error) {
	client := &http.Client{}
	var data = strings.NewReader(`{"detail":false}`)
	url := config.Instance().TRX_URL2 + "/walletsolidity/getblock"
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%s\n", bodyText)

	type TmpS struct {
		BlockId     string `json:"blockID"`
		BlockHeader struct {
			RawData struct {
				Number    int64 `json:"number"`
				Timestamp int64 `json:"timestamp"`
			} `json:"raw_data"`
		} `json:"block_header"`
	}

	var tmp TmpS
	err = json.Unmarshal(bodyText, &tmp)
	if err != nil {
		return 0, "", 0, err
	}

	return tmp.BlockHeader.RawData.Number, tmp.BlockId, tmp.BlockHeader.RawData.Timestamp, nil
}

func SendRawTrans(rawTrans string) (hash string, err error) {
	//without 0x
	client := &http.Client{}
	var data = strings.NewReader(`{"transaction": "` + rawTrans + `"}`)
	url := config.Instance().TRX_URL2 + "/wallet/broadcasthex"
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%s\n", bodyText)

	const CODE_SUCCESS = "SUCCESS"
	type TmpS struct {
		Result bool   `json:"result"`
		Code   string `json:"code"`
		TxId   string `json:"txid"`
	}
	var tmp TmpS
	err = json.Unmarshal(bodyText, &tmp)
	if err != nil {
		return "", err
	}
	if CODE_SUCCESS != tmp.Code {
		return "", errors.New("code is not success")
	}
	return tmp.TxId, nil
}

type TrxBalance struct {
	Balance      string `json:"balance"`
	TokenName    string `json:"tokenName"`
	TokenId      string `json:"tokenId"` //trx = "_"
	TokenDecimal int32  `json:"tokenDecimal"`
}

func GetBalance(addr string) (balances []TrxBalance, isActive bool, err error) {
	//curl "https://apilist.tronscanapi.com/api/accountv2?address=TPfHmimxS2drgvxzDiAbgFahrcsBoZG9EQ"
	url := config.Instance().TRX_URL + "/api/accountv2?address=" + addr

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set(TRON_PRO_API_KEY, config.Instance().TRX_API_KEY)
	resp, err := TRXClient.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false, err
	}

	type TmpS struct {
		ActivePermissions []struct{}   `json:"activePermissions"`
		WithPriceTokens   []TrxBalance `json:"WithPriceTokens"`
	}

	var tmp TmpS
	err = json.Unmarshal(bodyText, &tmp)
	if err != nil {
		return nil, false, err
	}
	if len(tmp.ActivePermissions) > 0 {
		isActive = true
	}
	balances = tmp.WithPriceTokens
	err = nil
	return
}

func NewTrxTransfer(fromAddr string, toAddr string, amount int64) (string, error) {
	blockNumber, blockHash, blockTime, err := GetLatestBlock()
	if err != nil {
		return "", err
	}
	_ = blockTime
	//fmt.Println("debug0:", blockNumber, blockHash, blockTime)
	currentTime := time.Now()
	//fmt.Println("debug1:", currentTime.UnixMilli())
	k1 := make([]byte, 8)
	binary.BigEndian.PutUint64(k1, uint64(blockNumber))
	k2, _ := hex.DecodeString(blockHash)
	refBlockBytes := hex.EncodeToString(k1[6:8])
	refBlockHash := hex.EncodeToString(k2[8:16])
	expiration := currentTime.UnixMilli() + 3600*1000
	timestamp := currentTime.UnixMilli()

	owner, err := tron.GetAddressHash(fromAddr)
	if err != nil {
		return "", err
	}
	to, err := tron.GetAddressHash(toAddr)
	if err != nil {
		return "", err
	}
	transferContract := &pb.TransferContract{OwnerAddress: owner, ToAddress: to, Amount: amount}
	param, err := ptypes.MarshalAny(transferContract)
	if err != nil {
		return "", err
	}
	contract := &pb.Transaction_Contract{Type: pb.Transaction_Contract_TransferContract, Parameter: param}
	raw := new(pb.TransactionRaw)
	refBytes, err := hex.DecodeString(refBlockBytes)
	if err != nil {
		return "", err
	}
	raw.RefBlockBytes = refBytes
	refHash, err := hex.DecodeString(refBlockHash)
	if err != nil {
		return "", err
	}
	raw.RefBlockHash = refHash
	raw.Expiration = expiration
	raw.Timestamp = timestamp
	raw.Contract = []*pb.Transaction_Contract{contract}
	trans := pb.Transaction{RawData: raw}
	data, err := proto.Marshal(&trans)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}
func NewTrx20Transfer(fromAddress string, toAddress string, amount *big.Int, contractAddress string) (string, error) {
	blockNumber, blockHash, blockTime, err := GetLatestBlock()
	if err != nil {
		return "", err
	}
	_ = blockTime
	currentTime := time.Now()
	k1 := make([]byte, 8)
	binary.BigEndian.PutUint64(k1, uint64(blockNumber))
	k2, _ := hex.DecodeString(blockHash)
	refBlockBytes := hex.EncodeToString(k1[6:8])
	refBlockHash := hex.EncodeToString(k2[8:16])
	expiration := currentTime.UnixMilli() + 3600*1000
	timestamp := currentTime.UnixMilli()

	raw := new(pb.TransactionRaw)
	refBytes, err := hex.DecodeString(refBlockBytes)
	if err != nil {
		return "", err
	}
	raw.RefBlockBytes = refBytes
	refHash, err := hex.DecodeString(refBlockHash)
	if err != nil {
		return "", err
	}
	raw.RefBlockHash = refHash
	raw.Expiration = expiration
	raw.Timestamp = timestamp

	fromAddressHash, err := tron.GetAddressHash(fromAddress)
	if err != nil {
		return "", err
	}
	toAddressHash, err := tron.GetAddressHash(toAddress)
	if err != nil {
		return "", err
	}
	contractAddressHash, err := tron.GetAddressHash(contractAddress)
	if err != nil {
		return "", err
	}
	input, err := token.Transfer(hex.EncodeToString(toAddressHash), amount)
	if err != nil {
		return "", err
	}
	transferContract := &pb.TriggerSmartContract{OwnerAddress: fromAddressHash, ContractAddress: contractAddressHash, CallValue: 0, CallTokenValue: 0, Data: input}
	param, err := ptypes.MarshalAny(transferContract)
	if err != nil {
		return "", err
	}

	contract := &pb.Transaction_Contract{Type: pb.Transaction_Contract_TriggerSmartContract, Parameter: param}
	raw.FeeLimit = GAS_USDT

	raw.Contract = []*pb.Transaction_Contract{contract}
	trans := pb.Transaction{RawData: raw}
	data, err := proto.Marshal(&trans)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}

// without 0x
func Sign(rawTrans, priv string) (trans string, err error) {
	//fmt.Println("debug0")
	rawTrans2, err := tron.SignStart(rawTrans)
	if err != nil {
		//fmt.Println("debug1")
		return "", err
	}
	//fmt.Println("debug2")
	privKey, _ := btcec.PrivKeyFromBytes(hexutils.HexToBytes(priv))
	if err != nil {
		return "", err
	}
	//fmt.Println("debug4")
	signature, err := tron.Sign(rawTrans2, privKey)
	if err != nil {
		//fmt.Println("debug5")
		return "", err
	}
	//fmt.Println("debug6")
	trans, err = tron.SignEnd(rawTrans, signature)
	return
}

func GetTrxBalance(addr string) (balance int64, err error) {
	client := &http.Client{}
	var data = strings.NewReader(`
{
  "address": "` + addr + `",
  "visible": true
}
`)
	req, err := http.NewRequest("POST", "https://api.trongrid.io/walletsolidity/getaccount", data)
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
		Balance int64 `json:"balance"`
	}

	var tmp TmpS
	err = json.Unmarshal(bodyText, &tmp)
	if err != nil {
		return 0, err
	}

	return tmp.Balance, nil
}
