package service

import (
	"encoding/json"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	bitcoin "wallet/btc"
	"wallet/config"
	"wallet/eth"
	"wallet/trx"
)

var app = gin.Default()

func initCors() {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}                                        // 允许什么域名访问，支持多个域名
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}  // 允许的 HTTP 方法
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"} // 允许的 HTTP 头
	config.AllowCredentials = true
	// 设置cors中间件
	app.Use(cors.New(config))
}

func wrapError(err error) string {
	var tmp struct {
		Error string `json:"error"`
	}
	tmp.Error = err.Error()

	ret, _ := json.Marshal(tmp)
	return string(ret)
}

func StartServer() {
	initCors()
	initETHRouter()
	initBTCRouter()
	initTRXRouter()

	config.InitConfig()
	eth.ETHSyncInit()
	bitcoin.BTCSyncInit()
	trx.TRXInit()

	initETHSync(config.Instance().ETH_START_NUMBER)
	initBTCSync(config.Instance().BTC_START_NUMBER)
	initTRXSync(config.Instance().TRX_START_NUMBER)
	app.Run("0.0.0.0:10000")
}
