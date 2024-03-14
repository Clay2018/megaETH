package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"wallet/config"
	"wallet/db"
	"wallet/trx"
)

var trxRouter *gin.RouterGroup

func initTRXRouter() {
	ethRouter = app.Group("/tron")
	ethRouter.POST("/sendRawTransaction", TRXSendRawTransaction)
}

type TRXSendRawTransactionReq struct {
	RawTrans string `json:"raw_trans"`
}

type TRXSendRawTransactionResp struct {
	TxHash string `json:"txHash"`
}

func TRXSendRawTransaction(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	bodyText, err := c.GetRawData()
	if err != nil {
		c.String(501, wrapError(err))
		return
	}
	var req TRXSendRawTransactionReq
	err = json.Unmarshal(bodyText, &req)
	if err != nil {
		c.String(501, wrapError(err))
		return
	}
	if req.RawTrans == "" {
		c.String(501, PARAM_IS_INVALID)
		return
	}

	txHash, err := trx.SendRawTrans(req.RawTrans)
	if err != nil {
		c.String(501, wrapError(err))
		return
	}

	var resp TRXSendRawTransactionResp
	resp.TxHash = txHash
	ret, err := json.Marshal(resp)
	if err != nil {
		c.String(501, wrapError(err))
		return
	}
	c.String(200, string(ret))
	return
}

func initTRXSync(startNumber uint64) {

	if db.Instance() == nil {
		return
	}
	fmt.Println("start trx sync service")

	for i := 0; i < config.Instance().TRX_NUMBER; i++ {
		go func(j uint64) {
			var previousNumber = startNumber
			for {
				time.Sleep(time.Second * 1)

				blockNumber, _, blockTime, err := trx.GetLatestBlock()
				if err != nil {
					fmt.Println("error:", err.Error())
					continue
				}
				if previousNumber >= uint64(blockNumber) {
					continue
				}

				trxs, err := trx.GetTransaction(int64(previousNumber))
				if err != nil {
					fmt.Println("sync error, blockNumber:", previousNumber)
					continue
				}

				blockTime, err = trx.GetBlockTime(int64(previousNumber))
				contractTx, err := trx.GetTransferLog(blockTime-1, blockTime)
				if err != nil {
					fmt.Println("sync error, blockNumber:", previousNumber)
					continue
				}

				for _, item := range trxs {
					insertUserSql := "INSERT INTO trx_deposit(addr, tx_hash, value, contract_addr) VALUES(?, ?, ?, ?)"
					_, err := db.Instance().Exec(insertUserSql, item.To, item.Hash, strconv.FormatInt(item.Value, 10), item.ContractAddr)
					if err != nil {
						fmt.Println("insert fail")
						return
					}
				}
				for _, item := range contractTx {
					insertUserSql := "INSERT INTO trx_deposit(addr, tx_hash, value, contract_addr) VALUES(?, ?, ?, ?)"
					_, err := db.Instance().Exec(insertUserSql, item.To, item.Hash, strconv.FormatInt(item.Value, 10), item.ContractAddr)
					if err != nil {
						fmt.Println("insert fail")
						return
					}
				}

				fmt.Println("success, blockNumber:", previousNumber)
				previousNumber += j
			}
		}(uint64(i))

	}
}
