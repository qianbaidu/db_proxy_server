package mysql

import (
	"bytes"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/common/log"
	"github.com/qianbaidu/db-proxy/proxy/mysql/provider"
	"github.com/qianbaidu/db-proxy/util"
	. "github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go/hack"
	"time"
)

type ExecResult struct {
	Result    *Result `json:"-"`
	StartTime string  `json:"create_time"`
	EndTime   string  `json:"update_datatime"`
	Status    int     `json:"status"`
	User      string  `json:"username"`
	Sql       string  `json:"sql"`
	DbAlias   string  `json:"name"`
}

func HandleProxy(c *provider.Conn) (err error) {
	if c.Conn == nil {
		return provider.ProxyConnectClosed
	}

	data, err := c.ReadPacket()
	if err != nil {
		log.Error("ReadPacket error : ", err)

		c.Conn = nil
		c.Close()
		return err
	}

	// todo 使用mysql连接池
	go c.H.CheckMysqlConn()

	v := mysqlDispatch(c, data)

	err = c.WriteValue(v)

	if c.Conn != nil {
		c.ResetSequence()
	}

	if err != nil {
		c.Close()
		c.Conn = nil
	}
	return err
}

func mysqlDispatch(c *provider.Conn, data []byte) interface{} {
	cmd := data[0]
	data = data[1:]
	switch cmd {
	case COM_QUIT:
		log.Info(fmt.Sprintf("user:%s cmd:quit ", c.GetUser()))
		c.Close()
		c.Conn = nil
		return provider.NoResponse{}
	case COM_QUERY:
		cmd := util.Trim(hack.String(data))
		log.Info(fmt.Sprintf("[%s] user:%s cmd:%s", time.Now().Format("2006-01-02 15:04:05"),
			c.GetUser(), hack.String(data)))
		if validateReg.MatchString(cmd) == true && selectUpdateReg.MatchString(cmd) == false {
			if r, err := c.H.HandleQuery(cmd); err != nil {
				log.Error(fmt.Sprintf("COM_QUERY[   %s   ]HandleQuery error: %s", hack.String(data), err))
				log.Info("connect : ", c.H)
				return err
			} else {
				return r
			}
		} else {
			return fmt.Errorf("not supported . (cmd: %s)", cmd)
		}
	case COM_PING:
		log.Info(fmt.Sprintf("user:%s cmd:COM_PING", c.GetUser()))
		return nil
	case COM_INIT_DB:
		log.Info(fmt.Sprintf("user:%s cmd:COM_INIT_DB, init db:%s", c.GetUser(), hack.String(data)))
		if _, err := c.H.UseDB(hack.String(data)); err != nil {
			log.Error("COM_INIT_DB useDb error ", err)
			return err
		} else {
			return nil
		}
	case COM_FIELD_LIST:
		log.Info("COM_FIELD_LIST")
		index := bytes.IndexByte(data, 0x00)
		table := hack.String(data[0:index])
		wildcard := hack.String(data[index+1:])

		if fs, err := c.H.HandleFieldList(table, wildcard); err != nil {
			log.Error("HandleFieldList error ", err)
			return err
		} else {
			return fs
		}
	case COM_STMT_PREPARE:
		log.Info(fmt.Sprintf("user:%s cmd:COM_STMT_PREPARE ", c.GetUser()))
		return fmt.Errorf("COM_STMT_PREPARE not supported now")

		c.StmtID++
		st := new(provider.Stmt)
		st.ID = c.StmtID
		st.Query = hack.String(data)

		log.Info(fmt.Sprintf("user:%s cmd:COM_STMT_PREPARE, sql: %s ", c.GetUser(), st.Query))
		var err error
		if st.Params, st.Columns, st.Context, err = c.H.HandleStmtPrepare(st.Query); err != nil {
			return err
		} else {
			st.ResetParams()
			c.Stmts[c.StmtID] = st
			return st
		}
	case COM_STMT_EXECUTE:
		log.Info(fmt.Sprintf("user:%s cmd:COM_STMT_EXECUTE, sql: %s ", c.GetUser(), data))
		return fmt.Errorf("COM_STMT_EXECUTE not supported now")

		if r, err := c.HandleStmtExecute(data); err != nil {
			return err
		} else {
			return r
		}
	case COM_STMT_CLOSE:
		log.Info(fmt.Sprintf("user:%s cmd:COM_STMT_CLOSE", c.GetUser()))
		c.HandleStmtClose(data)
		return provider.NoResponse{}
	case COM_STMT_SEND_LONG_DATA:
		log.Info(fmt.Sprintf("user:%s cmd:COM_STMT_SEND_LONG_DATA,data:%s", c.GetUser(), data))
		c.HandleStmtSendLongData(data)
		return provider.NoResponse{}
	case COM_STMT_RESET:
		log.Info(fmt.Sprintf("user:%s cmd:COM_STMT_RESET,data:%s", c.GetUser(), data))
		if r, err := c.HandleStmtReset(data); err != nil {
			log.Error("HandleStmtReset error ", err)
			return err
		} else {
			return r
		}
	default:
		log.Info(fmt.Sprintf("user:%s cmd:default,data:%s", c.GetUser(), data))
		return fmt.Errorf("COM_STMT_EXECUTE not supported now")

		return c.H.HandleOtherCommand(cmd, data)
	}

	return fmt.Errorf("command %d is not handled correctly", cmd)
}
