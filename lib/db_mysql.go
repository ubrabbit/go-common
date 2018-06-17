package lib

/*
https://godoc.org/github.com/go-sql-driver/mysql
https://github.com/go-sql-driver/mysql#usage
https://github.com/go-sql-driver/mysql/wiki/Examples

user@unix(/path/to/socket)/dbname?charset=utf8
user:password@tcp(localhost:5555)/dbname?charset=utf8
user:password@/dbname
user:password@tcp([de:ad:be:ef::ca:fe]:80)/dbname
*/

import (
	"fmt"
	"sync"
)

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	. "github.com/ubrabbit/go-common/common"
)

const (
	MAX_OPEN_MYSQL_CONNECTIONS = 32
	MAX_IDLE_MYSQL_CONNECTIONS = 8
)

var (
	g_MysqlDB *MysqlConnect = nil
)

type MysqlConnect struct {
	sync.Mutex
	db *sql.DB
}

func NewMysqlConn(host string, port int, dbname string, username string, password string) *MysqlConnect {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", username, password, host, port, dbname)
	LogInfo("Connect Mysql: %s %d", host, port)
	db, err := sql.Open("mysql", dataSourceName)
	CheckFatal(err)

	db.SetMaxOpenConns(MAX_OPEN_MYSQL_CONNECTIONS)
	db.SetMaxIdleConns(MAX_IDLE_MYSQL_CONNECTIONS)
	err = db.Ping()
	CheckFatal(err)

	LogInfo(fmt.Sprintf("Connect Mysql %s Succ", host))
	return &MysqlConnect{db: db}
}

func InitMysql(host string, port int, dbname string, username string, password string) {
	LogInfo("InitMysql %s:%d %s", host, port, dbname)
	conn := NewMysqlConn(host, port, dbname, username, password)
	g_MysqlDB = conn
	LogInfo("InitMysql Success")
}

func (self *MysqlConnect) Close() {
	self.Lock()
	defer func() {
		err := recover()
		if err != nil {
			LogError("Mysql Close Error: %v", err)
		}
		self.Unlock()
	}()
	self.db.Close()
}

func (self *MysqlConnect) Transaction(txFunc func(*sql.Tx, ...interface{}) error, args ...interface{}) (err error) {
	tx, err := self.db.Begin()
	CheckFatal(err)
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			LogPanic("Transaction Fatal: %v", p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = txFunc(tx, args...)
	if err != nil {
		LogError("Transaction Error: %v", err)
	}
	return err
}

func (self *MysqlConnect) TransactionExec(sql string, args ...interface{}) (result interface{}, err error) {
	tx, err := self.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			LogPanic("TransactionExec Fatal: %v", p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	result, err = tx.Exec(sql, args...)
	if err != nil {
		LogError("TransactionExec Error: %v", err)
	}
	return result, err
}

func (self *MysqlConnect) PrepareStmt(sql string) *sql.Stmt {
	sql_stmt, err := self.db.Prepare(sql)
	CheckFatal(err)
	return sql_stmt
}

func (self *MysqlConnect) execStmt(sql string, arg ...interface{}) (sql.Result, error) {
	sql_stmt, err := self.db.Prepare(sql)
	defer sql_stmt.Close()
	CheckFatal(err)

	result, err := sql_stmt.Exec(arg...)
	return result, err
}

func (self *MysqlConnect) Query(sql string, arg ...interface{}) []map[string]interface{} {
	sql_stmt, err := self.db.Prepare(sql)
	CheckFatal(err)

	rows, err := sql_stmt.Query(arg...)
	CheckFatal(err)

	columns, err := rows.Columns()
	CheckFatal(err)

	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	data := make([]map[string]interface{}, 0)
	for rows.Next() {
		//将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		CheckFatal(err)

		record := make(map[string]interface{})
		for i, element := range values {
			switch value := element.(type) {
			case []byte:
				record[columns[i]] = string(value)
			default:
				record[columns[i]] = value
			}
		}
		data = append(data, record)
	}
	return data
}

func (self *MysqlConnect) SeekDB(sql string) chan map[string]interface{} {
	sql = fmt.Sprintf("%s LIMIT ?,?", sql)
	LogDebug("sql_stmt: %s", sql)

	cache_size := 10000
	ch := make(chan map[string]interface{}, cache_size)
	go func() {
		start, seek_cnt := 0, cache_size
		for {
			res := self.Query(sql, start, seek_cnt)
			if len(res) <= 0 {
				break
			}
			for _, record := range res {
				ch <- record
			}
			start += seek_cnt
		}
		ch <- nil
	}()
	return ch
}

func checkDBConn() {
	if g_MysqlDB == nil {
		LogFatal("DB Conn is not inited!!!!!")
	}
}

func MysqlClose() {
	if g_MysqlDB != nil {
		g_MysqlDB.Close()
		g_MysqlDB = nil
	}
}

func MysqlQuery(sql string, arg ...interface{}) []map[string]interface{} {
	checkDBConn()
	return g_MysqlDB.Query(sql, arg...)
}

func MysqlUpdate(sql string, arg ...interface{}) int64 {
	checkDBConn()

	result, err := g_MysqlDB.execStmt(sql, arg...)
	CheckFatal(err)

	num, err := result.RowsAffected()
	CheckFatal(err)
	return num
}

func MysqlInsert(sql string, arg ...interface{}) int64 {
	checkDBConn()

	result, err := g_MysqlDB.execStmt(sql, arg...)
	CheckFatal(err)

	lastid, err := result.LastInsertId()
	CheckFatal(err)
	return lastid
}

func MysqlDelete(sql string, arg ...interface{}) int64 {
	checkDBConn()

	result, err := g_MysqlDB.execStmt(sql, arg...)
	CheckFatal(err)

	num, err := result.RowsAffected()
	CheckFatal(err)
	return num
}

func MysqlTransaction(f func(*sql.Tx, ...interface{}) error, args ...interface{}) error {
	checkDBConn()
	return g_MysqlDB.Transaction(f, args...)
}

func MysqlTransactionExec(sql string, args ...interface{}) (interface{}, error) {
	checkDBConn()
	return g_MysqlDB.TransactionExec(sql, args...)
}

func MysqlSeekDB(sql string) chan map[string]interface{} {
	checkDBConn()
	return g_MysqlDB.SeekDB(sql)
}
