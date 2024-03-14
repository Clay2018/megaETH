package bitcoin

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

func CreateAddress(network string) (privateKeyHex string, publicKeyHex string, address string) {
	var netParams *chaincfg.Params

	switch network {
	case MAINNET:
		netParams = &chaincfg.MainNetParams
	case TESTNET:
		netParams = &chaincfg.TestNet3Params
	}

	privateKey, _ := btcec.NewPrivateKey()
	privateKeyHex = hex.EncodeToString(privateKey.Serialize())
	privateKeyHex = HexToWif(network, privateKeyHex)

	publicKeyHex = hex.EncodeToString(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(privateKey.PubKey())))
	taprootAddress, _ := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(privateKey.PubKey())), netParams)

	address = taprootAddress.EncodeAddress()

	return
}

func HexToWif(network string, privHex string) string {
	var netParams *chaincfg.Params

	switch network {
	case MAINNET:
		netParams = &chaincfg.MainNetParams
	case TESTNET:
		netParams = &chaincfg.TestNet3Params
	}

	privateKeyBytes, _ := hex.DecodeString(privHex)
	privateKey, _ := btcec.PrivKeyFromBytes(privateKeyBytes)

	wif, _ := btcutil.NewWIF(privateKey, netParams, true)
	return wif.String()
}

func WifToHex(wif string) (privHex string) {
	wifS, _ := btcutil.DecodeWIF(wif)
	return hex.EncodeToString(wifS.PrivKey.Serialize())
}
