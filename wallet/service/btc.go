package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	bitcoin "wallet/btc"
	"wallet/db"
)

var btcRouter *gin.RouterGroup

func initBTCRouter() {
	btcRouter = app.Group("/btc")
	btcRouter.POST("/sendRawTransaction", BTCSendRawTransaction)
}

type BTCSendRawTransactionReq struct {
	RawTrans string `json:"raw_trans"`
}

type BTCSendRawTransactionResp struct {
	TxHash string `json:"txHash"`
}

func BTCSendRawTransaction(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	bodyText, err := c.GetRawData()
	if err != nil {
		c.String(501, wrapError(err))
		return
	}
	var req BTCSendRawTransactionReq
	err = json.Unmarshal(bodyText, &req)
	if err != nil {
		c.String(501, wrapError(err))
		return
	}
	if req.RawTrans == "" {
		c.String(501, PARAM_IS_INVALID)
		return
	}

	txHash, err := bitcoin.SendRawTrans(req.RawTrans)
	if err != nil {
		c.String(501, wrapError(err))
		return
	}

	var resp BTCSendRawTransactionResp
	resp.TxHash = txHash
	ret, err := json.Marshal(resp)
	if err != nil {
		c.String(501, wrapError(err))
		return
	}
	c.String(200, string(ret))
	return

}

func initBTCSync(startNumber uint64) {

	if db.Instance() == nil {
		return
	}
	fmt.Println("start btc sync service")

	go func() {
		var previousNumber = startNumber
		for {
			time.Sleep(time.Second * 10)

			latestNumber, err := bitcoin.GetLatestBlock()
			if err != nil {
				fmt.Println("error:", err.Error())
				continue
			}
			if previousNumber >= latestNumber {
				continue
			}

			btcTx, err := bitcoin.GetTransaction(previousNumber)
			if err != nil {
				fmt.Println("sync error, blockNumber:", previousNumber)
				continue
			}
			for _, item := range btcTx {
				fmt.Println(item.To, item.TxHash, item.Value, item.Vout)
				insertUserSql := "INSERT INTO btc_deposit(addr, tx_hash, value, vout) VALUES(?, ?, ?, ?)"
				_, err := db.Instance().Exec(insertUserSql, item.To, item.TxHash, item.Value, item.Vout)
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
