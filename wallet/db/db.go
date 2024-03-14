package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"wallet/config"
)

var mysqlInstance *sql.DB

func init() {
	dataSourceName := config.Instance().DBConfig.UserName + ":" + config.Instance().DBConfig.Pwd +
		"@tcp(" + config.Instance().DBConfig.URL + `)/` + config.Instance().DBConfig.Database
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("can't connect mysql db: err", err.Error())
	}
	mysqlInstance = db
}

func Instance() *sql.DB {
	return mysqlInstance
}
