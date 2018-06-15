package tests

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
)

import (
	. "github.com/ubrabbit/go-common/common"
	config "github.com/ubrabbit/go-common/config"
	lib "github.com/ubrabbit/go-common/lib"
)

//测试事物成功
func TransactionCallBack(tx *sql.Tx, args ...interface{}) error {
	var1 := args[0].(int)
	var2 := args[1].(int)
	var3 := args[2].(int)
	var4 := args[3].(int)
	LogInfo("TransactionCallBack: %d %d", var1, var2, var3, var4)
	sql := "UPDATE tbl_Account SET sm_autoLoad=? WHERE id=?"
	if _, err := tx.Exec(sql, var1, var2); err != nil {
		return err
	}
	if _, err := tx.Exec(sql, var3, var4); err != nil {
		return err
	}
	return nil
}

//测试事物失败回滚
func TransactionCallBackError(tx *sql.Tx, args ...interface{}) error {
	sql := "UPDATE tbl_Account SET sm_autoLoad=32"
	if _, err := tx.Exec(sql); err != nil {
		return err
	}
	return errors.New("Test Transaction Error")
}

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

	lib.MysqlTransaction(TransactionCallBack, 3, 1, 4, 2)
	lib.MysqlTransaction(TransactionCallBackError)

	LogInfo(">>>>>>>>>>>>>>>>>>> Test SeekDB")
	ch := lib.MysqlSeekDB("SELECT * FROM tbl_Account")
	for {
		record := <-ch
		if record == nil {
			break
		}
		LogInfo("tbl_Account entry: %v", record)
	}

	lib.MysqlClose()
}
