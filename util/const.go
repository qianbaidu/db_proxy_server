package util

import (
	"fmt"
	. "github.com/siddontang/go-mysql/mysql"
	"github.com/sirupsen/logrus"
)

const (
	DB_TYPE_MYSQL = 1
	DB_TYPE_MONGO = 2
	DB_TYPE_REIDS = 3
)

const (
	DB_PERMISSION_STATUS_ALLOW     = 1
	DB_PERMISSION_STATUS_NOT_ALLOW = 0
)

var DB_TYPE_LIST = map[int]string{1: "mysql", 2: "mongodb", 3: "redis"}

const (
	API_RESPONSE_STATUS_SUCCESS       = 0
	REST_API_RESPONSE_STATSUS_SUCCESS = 200
)

const (
	DEFAULT_ENCODE = "utf8"
)

const (
	// 解析sql语句
	API_MYSQL_PARSER = "%s/parser/mysql"

)

const (
	DEFAULT_EVENT_LOG_QUEUE_SIZE = 200
	DEFAULT_EVENT_LOG_PUSH_TIMER = 10
	DEFAULT_EXEC_SQL_TIME_OUT    = 10
	DEFAULT_PING_INTERVAL        = 5
	DEFAULT_SQL_LIMIT            = 50
	DEFAULT_RELOAD_DB_TIME       = 30
	DEFAULT_MONGO_CONN_TIME_OUT  = 10
	DEFAULT_MONGO_CONN_POOL_SIZE = 100
)
const (
	SQL_EXEC_START    = 1
	SQL_EXEC_FINISHED = 2
	SQL_EXEC_TIMEOUT  = 3
)

const SALT = "encode_mypass_salt01"

type Database struct {
	Addr     string
	User     string
	Password string
	DbName   string
	Charset  string
}

func GetDbType(t int) string {
	if res, ok := DB_TYPE_LIST[t]; ok {
		return res
	} else {
		logrus.Error("未知的数据库类型", t)
	}
	return ""
}

type ApiResponse struct {
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Status int         `json:"status"`
}

func ResetSecretFieldValue(v uint8) (string) {
	switch v {
	case MYSQL_TYPE_TINY, MYSQL_TYPE_SHORT, MYSQL_TYPE_FLOAT, MYSQL_TYPE_INT24, MYSQL_TYPE_TINY_BLOB, MYSQL_TYPE_LONG,
		MYSQL_TYPE_LONG_BLOB, MYSQL_TYPE_LONGLONG:
		return "0"
	case MYSQL_TYPE_DOUBLE, MYSQL_TYPE_NEWDECIMAL, MYSQL_TYPE_DECIMAL:
		return "0.0"
	case MYSQL_TYPE_TIMESTAMP, MYSQL_TYPE_DATE, MYSQL_TYPE_DATETIME, MYSQL_TYPE_NEWDATE, MYSQL_TYPE_TIMESTAMP2,
		MYSQL_TYPE_DATETIME2, MYSQL_TYPE_TIME2:
		return "00-00-00"
	case MYSQL_TYPE_TIME:
		return "00:00:00"
	case MYSQL_TYPE_YEAR:
		return "0000"
	case MYSQL_TYPE_VARCHAR, MYSQL_TYPE_VAR_STRING, MYSQL_TYPE_STRING, MYSQL_TYPE_NULL,
		MYSQL_TYPE_MEDIUM_BLOB, MYSQL_TYPE_BLOB, MYSQL_TYPE_ENUM, MYSQL_TYPE_SET:
		return "Secret Field"
	case MYSQL_TYPE_JSON:
		return `{"Secret Field":"true"}`
	case MYSQL_TYPE_BIT:
		return fmt.Sprintf("%b", "Secret Field")
	default:
		return ""
	}
}

var InformationSchema map[string]bool

func initInformationSchema() {
	InformationSchema = make(map[string]bool, 0)
	InformationSchema["information_schema.collations"] = true
	InformationSchema["information_schema"] = true
	InformationSchema["information_schema.schemata"] = true
	InformationSchema["information_schema.routines"] = true
	InformationSchema["information_schema.profiling"] = true
	InformationSchema["information_schema.collations"] = true
	InformationSchema["information_schema.character_sets"] = true
}

func init() {
	initInformationSchema()
}
