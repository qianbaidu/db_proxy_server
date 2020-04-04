package provider

import (
	"github.com/siddontang/go-mysql/client"
	. "github.com/siddontang/go-mysql/mysql"
)

type Handler interface {
	// mysql login auth
	Auth(userName string, salt, authData []byte) (bool, error)
	// init user db connect
	UseDB(dbName string) (*client.Conn, error)
	//handle COM_INIT_DB command, you can check whether the dbName is valid, or other.
	//UseDB(dbName string) error
	//handle COM_QUERY command, like SELECT, INSERT, UPDATE, etc...
	//If Result has a Resultset (SELECT, SHOW, etc...), we will send this as the response, otherwise, we will send Result
	HandleQuery(query string) (*Result, error)
	//handle COM_FILED_LIST command
	HandleFieldList(table string, fieldWildcard string) ([]*Field, error)
	//handle COM_STMT_PREPARE, params is the param number for this statement, columns is the column number
	//context will be used later for statement execute
	HandleStmtPrepare(query string) (params int, columns int, context interface{}, err error)
	//handle COM_STMT_EXECUTE, context is the previous one set in prepare
	//query is the statement prepare query, and args is the params for this statement
	HandleStmtExecute(context interface{}, query string, args []interface{}) (*Result, error)
	//handle COM_STMT_CLOSE, context is the previous one set in prepare
	//this handler has no response
	HandleStmtClose(context interface{}) error
	//handle any other command that is not currently handled by the library,
	//default implementation for this method will return an ER_UNKNOWN_ERROR
	HandleOtherCommand(cmd byte, data []byte) error

	InitGetUserDbList() error

	InitUserDbList(userDblist interface{}, user string) error

	GetUserDbList(user interface{}) bool

	CheckMysqlConn()

	Destory()
}
