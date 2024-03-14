package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"wallet/db"
	"wallet/eth"
)

var ethRouter *gin.RouterGroup

func initETHRouter() {
	ethRouter = app.Group("/eth")
	ethRouter.GET("/hello", Hello)
	ethRouter.POST("/sendRawTransaction", ETHSendRawTransaction)
}

func Hello(c *gin.Context) {
	c.String(200, "helllo")
	return
}

type ETHSendRawTransactionReq struct {
	RawTrans string `json:"raw_trans"`
}

type ETHSendRawTransactionResp struct {
	TxHash string `json:"txHash"`
}

func ETHSendRawTransaction(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	bodyText, err := c.GetRawData()
	if err != nil {
		c.String(501, wrapError(err))
		return
	}
	var req ETHSendRawTransactionReq
	err = json.Unmarshal(bodyText, &req)
	if err != nil {
		c.String(501, wrapError(err))
		return
	}
	if req.RawTrans == "" {
		c.String(501, PARAM_IS_INVALID)
		return
	}

	txHash, err := eth.SendRawTrans(req.RawTrans)
	if err != nil {
		c.String(501, wrapError(err))
		return
	}

	var resp ETHSendRawTransactionResp
	resp.TxHash = txHash
	ret, err := json.Marshal(resp)
	if err != nil {
		c.String(501, wrapError(err))
		return
	}
	c.String(200, string(ret))
	return

}

func initETHSync(startNumber uint64) {

	if db.Instance() == nil {
		return
	}
	fmt.Println("start eth sync service")

	go func() {
		var previousNumber = startNumber
		for {
			time.Sleep(time.Second * 10)

			latestNumber, err := eth.GetLatestBlock()
			if err != nil {
				fmt.Println("error:", err.Error())
				continue
			}
			if previousNumber >= latestNumber {
				continue
			}

			ethTx, err := eth.GetTransaction(previousNumber)
			contractTx, err := eth.GetTransferLog(previousNumber, previousNumber)
			if err != nil {
				fmt.Println("sync error, blockNumber:", previousNumber)
				continue
			}

			for _, item := range ethTx {
				insertUserSql := "INSERT INTO eth_deposit(addr, tx_hash, value, contract_addr) VALUES(?, ?, ?, ?)"
				_, err := db.Instance().Exec(insertUserSql, item.To.String(), item.Hash.String(), item.Value.String(), "")
				if err != nil {
					fmt.Println("insert fail")
					return
				}
			}
			for _, item := range contractTx {
				insertUserSql := "INSERT INTO eth_deposit(addr, tx_hash, value, contract_addr) VALUES(?, ?, ?, ?)"
				_, err := db.Instance().Exec(insertUserSql, item.To.String(), item.Hash.String(), item.Value.String(), item.ContractAddr.String())
				if err != nil {
					fmt.Println("insert fail")
					return
				}
			}

			fmt.Println("success, blockNumber:", previousNumber)
			previousNumber += 1
		}
	}()
}
