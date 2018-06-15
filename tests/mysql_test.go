package tests

import (
	"fmt"
	"testing"
)

import (
	. "github.com/ubrabbit/go-common/common"
	config "github.com/ubrabbit/go-common/config"
	lib "github.com/ubrabbit/go-common/lib"
)

func TestMysql(t *testing.T) {
	fmt.Printf("\n\n=====================  TestMysql  =====================\n")

	config.InitConfig("config_test.conf")
	cfg := config.GetMysqlConfig()
	lib.InitMysql(cfg.IP, cfg.Port, cfg.Database, cfg.Account, cfg.Password)
	LogInfo("Mysql Show Tables: ", lib.MysqlQuery("show tables"))

	result := lib.MysqlQuery("SELECT * FROM tbl_Account")
	LogInfo("SELECT * FROM tbl_Account: ", result)
	sql := "UPDATE tbl_Account SET sm_autoLoad=? WHERE id=1"
	data, err := lib.MysqlTransactionExec(sql, "1")
	if err != nil {
		LogFatal("Transaction Error:  %v", err)
	}
	LogInfo("Transaction Result: %v", data)
	lib.MysqlClose()
}
