package eth

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"wallet/config"
	"wallet/db"
)

var addrs = make(map[string]struct{}, 0)
var contractAddrs = make([]string, 0)
var topics = make([][]common.Hash, 0)

func ETHSyncInit() {
	var req = make([]string, 0)
	rows, err := db.Instance().Query("select addr, flag from table_eth_balance where flag=1")
	if err != nil {
		fmt.Println("select sb fail, err:", err.Error())
		return
	}
	var sellOrders = make([]db.TableBalance, 0)
	for rows.Next() {
		var sellOrder db.TableBalance
		err := rows.Scan(&sellOrder.Addr, &sellOrder.Flag)
		if err != nil {
			fmt.Println("select sb fail, err:", err.Error())
			return
		}
		sellOrders = append(sellOrders, sellOrder)
		req = append(req, sellOrder.Addr)
	}

	config.Instance().ETH_ADDRS = append(config.Instance().ETH_ADDRS, req...)

	for _, item := range req {
		addr := common.HexToAddress(item)
		addrs[addr.String()] = struct{}{}
	}

	usdt := common.HexToAddress(config.Instance().ETH_USDT_ADDR)
	usdc := common.HexToAddress(config.Instance().ETH_USDC_ADDR)
	contractAddrs = append(contractAddrs, usdt.String())
	contractAddrs = append(contractAddrs, usdc.String())

	topic := make([]common.Hash, 0)
	topic = append(topic, common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"))
	topics = append(topics, topic)
}

func GetLatestBlock() (blockNumber uint64, err error) {
	blockNumber, err = ETHClient.BlockNumber(context.Background())
	return
}

// deposit
type ETHTx struct {
	To    common.Address
	Value big.Int
	Hash  common.Hash
}

func GetTransaction(blockNumber uint64) (transactions []*ETHTx, err error) {
	block, err := ETHClient.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		return nil, err
	}
	transactions = make([]*ETHTx, 0)

	trans := block.Transactions()
	for _, item := range trans {
		if item.To() == nil {
			continue
		}
		_, exist := addrs[item.To().String()]
		if !exist {
			continue
		}
		transactions = append(transactions, &ETHTx{
			To:    *item.To(),
			Value: *item.Value(),
			Hash:  item.Hash(),
		})
	}
	return transactions, nil
}

type ETHContractTx struct {
	To           common.Address
	Value        big.Int
	ContractAddr common.Address
	Hash         common.Hash
}

func GetTransferLog(fromBlock uint64, toBlock uint64) (contractTrans []*ETHContractTx, err error) {
	contractAddrs2 := make([]common.Address, 0)
	for _, item := range contractAddrs {
		contractAddrs2 = append(contractAddrs2, common.HexToAddress(item))
	}
	logs, err := ETHClient.FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: contractAddrs2,
		Topics:    topics,
	})
	if err != nil {
		return nil, err
	}

	contractTrans = make([]*ETHContractTx, 0)
	for _, item := range logs {
		if item.Removed {
			continue
		}
		if len(item.Topics) != 3 {
			continue
		}
		to := item.Topics[2]
		value := item.Data
		to2 := common.HexToAddress(to.String())
		_, exist := addrs[to2.String()]
		if !exist {
			continue
		}
		value2 := big.NewInt(0).SetBytes(value)
		contractTrans = append(contractTrans, &ETHContractTx{
			To:           to2,
			Value:        *value2,
			ContractAddr: item.Address,
			Hash:         item.TxHash,
		})
	}
	return
}
