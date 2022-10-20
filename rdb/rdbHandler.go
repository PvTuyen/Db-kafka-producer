package rdb

import (
	"db-producer/common"
	"db-producer/kafkaProducer"
	"db-producer/logger"
	"db-producer/schema"
	"time"

	//"db-producer/system"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
	"strings"
)

type Rdb struct {
	Host         string
	Port         string
	User         string
	Password     string
	DatabaseName string
	DriverName   string
	Connect      *sqlx.DB
}

type DbAdapter interface {
	Select() string
}

func InitRdbInfo() *Rdb {
	rootDb := Rdb{
		Host:         getEnv("Host"),
		Port:         getEnv("Port"),
		User:         getEnv("User"),
		Password:     getEnv("Password"),
		DatabaseName: getEnv("DatabaseName"),
		DriverName:   getEnv("DriverName"),
	}
	rootDb.Connect = InitRdb(rootDb)

	return &rootDb
}

//InitRdb init connect sql
func InitRdb(rdb Rdb) *sqlx.DB {
	info := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable",
		rdb.Host, rdb.User, rdb.Password, rdb.Port, rdb.DatabaseName)
	connect, err := sqlx.Open(rdb.DriverName, info)
	if err != nil {
		logger.LogError(common.ERROR, common.CONNECT_SQL_FAIL)
	}
	connect.SetMaxOpenConns(100)
	connect.SetMaxIdleConns(5)

	return connect
}

// getEnv get value in file env
func getEnv(name string) string {
	err := godotenv.Load(".env")
	if err != nil {
		logger.LogError(common.ERROR, err.Error())

		os.Exit(1)
	}
	val := os.Getenv(name)
	if val == "" {
		logger.LogError(common.ERROR, common.GET_ENV_FAIL+name)

		os.Exit(1)
	}

	return val
}

//GetData get data root in database
func (rdb *Rdb) GetData(tableName string, columns []string) []map[string]interface{} {
	sql := getSql(rdb, tableName, columns)
	rows, err := rdb.Connect.Queryx(sql)
	if err != nil {
		logger.LogError(common.QUERY_FAIL, fmt.Sprintf("%v \n%v", sql, err.Error()))

		return nil
	}
	defer rows.Close()
	results := []map[string]interface{}{}
	for rows.Next() {
		resMap := make(map[string]interface{})
		resMap["table"] = strings.ReplaceAll(tableName, common.SNAPSHOT, "")
		resValue, err := rows.SliceScan()
		if err != nil {
			logger.LogError(common.ROWS_CONVERT_FAIL, fmt.Sprintf("%v \n%v", sql, err.Error()))

			return nil
		}
		if len(resValue) != len(columns) {
			logger.LogError(common.LEN_COLUMNS_NOT_EQUALS_LEN_ROWS, fmt.Sprintf("%v", sql))

			return nil
		}
		for i, col := range columns {
			resMap[col] = resValue[i]
		}

		results = append(results, resMap)
	}
	return results
}

func getSql(rdb *Rdb, tableName string, columns []string) string {
	sql := fmt.Sprintf("select %v from %v; ", strings.Join(columns, ", "),tableName )
	return sql
}

func genSqlWithPostgres(name string, columns []string) string{

	return ""
}

func ProcessGetData(pro *kafkaProducer.Producer, schema *schema.Schema, rdb *Rdb) {
	processDeleteDataSnapshot(pro, schema, rdb)
	processGetDataRoot(pro, schema, rdb)
	processGetDataSnapshot(pro, schema, rdb)
}

func processGetDataRoot(pro *kafkaProducer.Producer, schema *schema.Schema, rdb *Rdb) {
	for k, v :=  range schema.MapInfo{
		go func(tableName string, cols []string) {
			listMsg := rdb.GetData(tableName, cols)
			for _, msg := range listMsg {
				(*pro.MapTopic)[tableName] <- msg
			}
		}(k, v)
	}
}

func processGetDataSnapshot(pro *kafkaProducer.Producer, schema *schema.Schema, rdb *Rdb) {
	for tableName, cols :=  range schema.MapInfo{
		go getDataSnapshot(pro, tableName, cols, rdb)
	}
}
func getDataSnapshot(pro *kafkaProducer.Producer, tableName string, columns []string, rdb *Rdb) {
	columns = append(columns, common.STATUS)
	for  {
		listMsg := rdb.GetData(tableName + common.SNAPSHOT, columns)
		for _, msg := range listMsg {
			(*pro.MapTopic)[tableName] <- msg
		}
		time.Sleep(1 * time.Second)
	}
}

func processDeleteDataSnapshot(pro *kafkaProducer.Producer, s *schema.Schema, rdb *Rdb) {
	go func() {
		for  {
			select {
			case msg := <- *pro.Snapshot:
				msgMap := msg.(map[string]interface{})
				rdb.deleteSnapshot(msgMap, s)
			}
		}
	}()
}

func (rdb *Rdb) deleteSnapshot(msgMap map[string]interface{}, s *schema.Schema) {
	tableName := msgMap[common.TABLE]
	columns := s.MapInfo[fmt.Sprintf("%v", tableName)]
	columnsValue := []string{}
	for _, column := range columns {
		value, _ := msgMap[column]
		columnsValue = append(columnsValue, fmt.Sprintf("%v = '%v'", column, value))
	}
	sql := fmt.Sprintf("DELETE FROM %v%v WHERE  %v ;", tableName, common.SNAPSHOT, strings.Join(columnsValue, " AND "))

	_, err := rdb.Connect.Exec(sql)
	if err != nil {
		logger.LogError(common.DELETE_SNAPSHOT_FAIL, fmt.Sprintf("%v\n%v",sql,err.Error()))
		os.Exit(1)
		return
	}
}
