package eth

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"wallet/config"

	ERC20 "wallet/eth/erc20"
)

var ETHClient *ethclient.Client

const GAS_ETH = uint64(21000)
const GAS_USDT = uint64(100000)

func init() {
	var err error
	ETHClient, err = ethclient.Dial(config.Instance().ETH_URL)
	if err != nil {
		panic(err)
	}
}

func SignTrans(legacyTx *types.LegacyTx, chainId *big.Int, privHex string) (rawTrans string, hash string, err error) {
	privKeyECDSA, err := crypto.HexToECDSA(privHex)
	if err != nil {
		return "", "", err
	}

	signer := types.NewEIP155Signer(chainId)

	tx := types.NewTx(legacyTx)

	signedTx, err := types.SignTx(tx, signer, privKeyECDSA)
	if err != nil {
		return "", "", err
	}
	hash = signedTx.Hash().String()
	rawTransByte, err := signedTx.MarshalBinary()
	if err != nil {
		return "", "", err
	}
	rawTrans = hex.EncodeToString(rawTransByte)
	return rawTrans, hash, nil
}

func SendRawTrans(rawTrans string) (hash string, err error) {
	//without 0x
	rawTxData, err := hex.DecodeString(rawTrans)
	if err != nil {
		return "", err
	}
	var tx types.Transaction
	err = rlp.DecodeBytes(rawTxData, &tx)
	if err != err {
		return "", err
	}

	err = ETHClient.SendTransaction(context.Background(), &tx)
	if err != nil {
		return "", err
	}

	return tx.Hash().String(), nil
}

func GetAddrNonce(addr string) (nonce uint64, err error) {
	addr2 := common.HexToAddress(addr)
	nonce, err = ETHClient.NonceAt(context.Background(), addr2, nil)
	return
}

func GetGasPrice() (gasPrice *big.Int, err error) {
	gasPrice, err = ETHClient.SuggestGasPrice(context.Background())
	return
}

func GetChainId() (chainId *big.Int, err error) {
	chainId, err = ETHClient.ChainID(context.Background())
	return
}

func GetBalance(addr string) (balance *big.Int, err error) {
	balance, err = ETHClient.BalanceAt(context.Background(), common.HexToAddress(addr), nil)
	return
}

func GetERC20Balance(addr string, erc20Addr string) (balance *big.Int, err error) {
	erc20, err := ERC20.NewErc20Caller(common.HexToAddress(erc20Addr), ETHClient)
	if err != nil {
		return nil, err
	}
	balance, err = erc20.BalanceOf(&bind.CallOpts{}, common.HexToAddress(addr))
	if err != nil {
		return nil, err
	}
	return balance, nil
}
