package bitcoin

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"wallet/internal/ord"
	"wallet/pkg/btcapi"
)

func CreatBTCRawTransaction(network string, utxoPrivateKeyWifs []string,
	inputUtxoArray []UTXO, outputArray []*ord.OutPut, feeRate int64, changeAddress string) (rawTransaction string, err error) {
	var netParams *chaincfg.Params
	switch network {
	case MAINNET:
		netParams = &chaincfg.MainNetParams
	case TESTNET:
		netParams = &chaincfg.TestNet3Params
	default:
		return "", errors.New("invalid networ")
	}
	if len(utxoPrivateKeyWifs) != len(inputUtxoArray) {
		return "", errors.New("len(utxoPrivateKeyWifs) != len(inputUtxoArray)")
	}
	if feeRate <= 0 {
		return "", errors.New("invalid feeRate")
	}

	checkBalance := func(inputUtxoArray []UTXO, outputArray []*ord.OutPut) bool {
		var totalInput int64
		var totalOutput int64
		for _, Utxo := range inputUtxoArray {
			totalInput += Utxo.Value
		}
		for _, Output := range outputArray {
			totalOutput += Output.Value
		}

		return totalOutput <= totalInput
	}
	if !checkBalance(inputUtxoArray, outputArray) {
		return "", errors.New("not enough balance")
	}

	utxoPrivateKeys := make([]*btcec.PrivateKey, 0, len(utxoPrivateKeyWifs))

	for _, utxoPrivateKeyWif := range utxoPrivateKeyWifs {
		utxoPrivateKeyHex := WifToHex(utxoPrivateKeyWif)
		if len(utxoPrivateKeyHex) != PRIVATE_KEY_LENGTH {
			return "", errors.New("invalid privateKey")
		}
		utxoPrivateKeyBytes, err := hex.DecodeString(utxoPrivateKeyHex)
		if err != nil {
			return "", err
		}
		utxoPrivateKey, _ := btcec.PrivKeyFromBytes(utxoPrivateKeyBytes)
		utxoPrivateKeys = append(utxoPrivateKeys, utxoPrivateKey)
	}

	commitTxOutPointList := make([]*wire.OutPoint, 0)
	commitTxOutPointList2 := make([]*btcapi.UnspentOutput, 0)
	commitTxPrivateKeyList := make([]*btcec.PrivateKey, 0)

	length := len(inputUtxoArray)
	for i := 0; i < length; i++ {
		utxoHash, _ := chainhash.NewHashFromStr(inputUtxoArray[i].Hash)
		outPoint := wire.NewOutPoint(utxoHash, inputUtxoArray[i].Index)

		pkScriptBytes, _ := hex.DecodeString(inputUtxoArray[i].PkScriptHex)
		output := wire.NewTxOut(inputUtxoArray[i].Value, pkScriptBytes)

		unspent := btcapi.UnspentOutput{
			Outpoint: outPoint,
			Output:   output,
		}

		commitTxOutPointList = append(commitTxOutPointList, unspent.Outpoint)
		commitTxOutPointList2 = append(commitTxOutPointList2, &unspent)
		commitTxPrivateKeyList = append(commitTxPrivateKeyList, utxoPrivateKeys[i])
	}

	request := ord.BTCRequest{
		CommitTxOutPointList2:  commitTxOutPointList2,
		CommitTxOutPointList:   commitTxOutPointList,
		CommitTxPrivateKeyList: commitTxPrivateKeyList,
		FeeRate:                feeRate,
		OutPuts:                outputArray,
		ChangeAddress:          changeAddress,
	}

	tool, err := ord.NewBTCToolWithBtcApiClient(netParams, &request)
	if err != nil {
		return "", err
	}

	rawTransaction, err = tool.GetCommitTxHex()

	return
}
