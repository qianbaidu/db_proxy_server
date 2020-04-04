package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	"github.com/qianbaidu/db-proxy/models"

	"github.com/qianbaidu/db-proxy/service"
	"github.com/qianbaidu/db-proxy/util"
	"github.com/siddontang/go-mysql/client"
	. "github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go/hack"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type MysqlProxyHandler struct {
	User       string
	UserId     int64
	IsSupper   bool
	AuthData   string
	UserDbList service.UserDbList
	CurrentDb  string
	Conn       *client.Conn
}

var (
	showDbReg, _      = regexp.Compile(`show\s+databases`)
	selectSchemata, _ = regexp.Compile(`select(.*?)from\sinformation_schema\.schemata`)
	showTableReg, _   = regexp.Compile(`show(.*?)tables`)
	showSchemaReg, _  = regexp.Compile(`^show`)
	rollbackReg, _    = regexp.Compile(`^rollback`)
	useDbReg          = regexp.MustCompile(`^use\s+([\w-_]+)`)
	setDbReg          = regexp.MustCompile(`^set\s+([\w-_]+)`)
	selectReg         = regexp.MustCompile(`^select`)
	limitReg          = regexp.MustCompile(`\s+limit\s+`)
	limitNumReg       = regexp.MustCompile(`limit\s+(\d+)$`)
	limitOffsetReg    = regexp.MustCompile(`(\s+limit\s+(\d+)(,|\s{0,},\s{0,}|\s+offset\s+)(\d+))`)
	validateReg       = regexp.MustCompile(`^(use|show|select|set|explain|rollback)`)
	selectUpdateReg   = regexp.MustCompile(`^(select(.*?)for\s+update)`)
	sqlErrorReg       = regexp.MustCompile(`SQL syntax|Unknown column|Undeclared variable`)
)

func (h *MysqlProxyHandler) InitDb(dbName string) (conn *client.Conn, err error) {
	if len(dbName) > 0 {
		h.CurrentDb = dbName
	} else {
		for _, db := range h.UserDbList.DBList {
			return h.dbConnect(db.DbList)
		}
	}

	if res, ok := h.UserDbList.DBList[dbName]; ok {
		return h.dbConnect(res.DbList)
	} else {
		return conn, util.ERR_DATABASES_NOT_EXISTS
	}
}

func (h *MysqlProxyHandler) Auth(user string, salt, authData []byte) (bool, error) {
	if userId, err := service.MysqlLoginAuth(user, salt, authData); err != nil || userId < 1 {
		return false, err
	} else {
		h.UserId = userId
		return true, nil
	}
}

func (h *MysqlProxyHandler) UseDB(dbName string) (conn *client.Conn, err error) {
	if _, ok := h.UserDbList.DBList[dbName]; !ok {
		return nil, errors.New(fmt.Sprintf("No database[%s] access", dbName))
	}
	conn, err = h.InitDb(dbName)
	if err != nil {
		log.Error("user db,init db error ", err)
		return
	}
	h.Conn = conn
	return conn, nil
}

func (h *MysqlProxyHandler) dbConnect(res models.DbList) (conn *client.Conn, err error) {
	conn, err = client.Connect(fmt.Sprintf("%s:%d", res.Ip, res.Port), res.Username, res.Password, res.DbName)
	if err != nil {
		log.Error("mysql connect error ", err)
		return conn, err
	}

	// todo 根据数据库配置编码设置
	if err := conn.SetCharset("utf8"); err != nil {
		log.Error("conn set charset utf8 error ", err)
	}

	h.Conn = conn
	h.CurrentDb = res.DbName

	return conn, nil
}

func (h *MysqlProxyHandler) killQuery(threadId uint32) (r sql.Result, err error) {
	if db, ok := h.UserDbList.DBList[h.CurrentDb]; ok {
		dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", db.Username, db.Password, db.Ip, db.Port, db.DbName)
		conn, err := sql.Open("mysql", dns)
		defer conn.Close()
		if err != nil {
			return r, err
		}

		cmd := fmt.Sprintf("kill query %d", threadId)
		r, err = conn.Exec(cmd)
		if err != nil {
			log.Error("exec kill query error : ", err, " cmd : ", cmd)
			return r, err
		}
		return r, nil

	}
	return r, errors.New("kill failed")
}

func (h *MysqlProxyHandler) getBdId() (int64, error) {
	if db, ok := h.UserDbList.DBList[h.CurrentDb]; ok {
		return db.DbList.Id, nil
	} else {
		log.Error("get current db id error ")
		return 0, errors.New("get current db id error")
	}
}

func (h *MysqlProxyHandler) NewEventLog(sql string) (e *models.EventLog) {
	dbId, _ := h.getBdId()
	e = &models.EventLog{
		Sql:            sql,
		Status:         0,
		CreateTime:     time.Now().Format("2006-01-02 15:04:05"),
		DbId:           dbId,
		UserId:         h.UserId,
		UpdateDatatime: time.Now().Format("2006-01-02 15:04:05"),
	}
	return e
}

func (h *MysqlProxyHandler) buildSimpleShowResultset(values []interface{}, name string, sql string) (result *Result, err error) {
	evnetLog := h.NewEventLog(sql)
	r := new(Resultset)
	field := &Field{}
	field.Name = hack.Slice(name)
	field.Charset = 33
	field.Type = MYSQL_TYPE_VAR_STRING

	r.Fields = []*Field{field}

	var row []byte

	for _, value := range values {
		row, err = util.FormatValue(value)
		if err != nil {
			return nil, err
		}
		r.RowDatas = append(r.RowDatas,
			PutLengthEncodedString(row))
	}

	result = &Result{Resultset: r,}
	evnetLog.Status = util.SQL_EXEC_FINISHED
	service.RecordEvent(evnetLog)

	return result, nil
}

func (h *MysqlProxyHandler) HandleSetSql(sql string) (res *Result, err error) {
	return h.Execute(sql)
}

func filterDbName(sql string) (newSql string) {
	newSql = sql
	s := strings.Split(sql, "from")
	if len(s) >= 3 {
		newSql = fmt.Sprintf("%s from %s ", s[0], s[1])
	}
	return newSql
}

func (h *MysqlProxyHandler) execRollback(sql string) (res *Result, err error) {
	return res, nil
}

func (h *MysqlProxyHandler) execShowSql(sql string) (*Result, error) {
	sql = filterDbName(sql)
	return h.Execute(sql)
}

func (h *MysqlProxyHandler) HandleQuery(sql string) (res *Result, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error("panic error ", err, " sql : ", sql)
		}
	}()

	// use db
	if useDbParms := useDbReg.FindStringSubmatch(strings.Replace(sql, "`", "", -1)); len(useDbParms) >= 2 {
		log.Info("HandleQuery start, user db start ", useDbParms)
		if _, err := h.UseDB(useDbParms[1]); err != nil {
			log.Error("MatchUseDb useDbParms,use db errr ", err)
			return res, err
		} else {
			return res, nil
		}
	}

	if h.Conn == nil {
		log.Info(fmt.Sprintf("reinit db connect,db:%s, ", h.CurrentDb))
		_, err := h.InitDb(h.CurrentDb)
		if err != nil {
			log.Error("InitDb error ", err)
			return nil, err
		}
	}

	// set
	if setDbReg.MatchString(sql) {
		return h.HandleSetSql(sql)
	}

	// show databases
	if showDbReg.MatchString(sql) {
		return h.showDabases(sql)
	}
	// 容错navicat客户端拉取数据库雷彪
	if selectSchemata.MatchString(sql) {
		return h.selectSchemata(sql)
	}

	// show tables
	if showTableReg.MatchString(sql) {
		return h.showTables(sql)
	}

	// 兼容客户端连接查询 show参数给予权限
	if showSchemaReg.MatchString(sql) {
		return h.execShowSql(sql)
	}

	if rollbackReg.MatchString(sql) {
		return h.execRollback(sql)
	}

	if h.IsSupper == false {
		// 表查询权限
		if r, err := h.checkTablePermission(sql); err != nil {
			log.Error("checkTablePermission error : ", err, " sql : ", sql)
			return res, NewError(ER_UNKNOWN_ERROR, err.Error())
		} else {
			//todo 使用配置参数确定是否开启limit限制
			//if selectReg.MatchString(sql) {
			//	sql = parserLimit(sql)
			//}
			//log.Info(fmt.Sprintf("[%s] parserLimit sql: %s", time.Now().Format("2006-01-02 15:04:05"), sql))
			res, err := h.Execute(sql)
			if err != nil {
				return res, err
			}
			if err := h.resultFilter(res, r); err != nil {
				return nil, err
			}
			return res, nil
		}
	}

	return h.Execute(sql)
}

var connectBadError = errors.New("io.ReadFull(header) failed. err EOF: connection was bad")

func (h *MysqlProxyHandler) tryExecute(sql string) (result *Result, err error) {
	var (
		tryTimes int
		conn     *client.Conn
	)
	for tryTimes < 3 {
		conn, err = h.InitDb(h.CurrentDb)
		if err != nil {
			log.Error(fmt.Sprintf("exec sql error, times : %d ,retry connec error : ", tryTimes, err))
			time.Sleep(3 * time.Second)
			tryTimes++
		} else {
			result, err = conn.Execute(sql)
			log.Info("retry connect and execte sql , error : ", err)
			break
		}
	}
	return
}

func (h *MysqlProxyHandler) Execute(sql string) (result *Result, err error) {
	e := h.NewEventLog(sql)

	ch := make(chan bool, 1)
	go func(c chan bool) {
	err = h.Conn.SetAutoCommit()
	if err != nil {
		log.Error("SetAutoCommit error : ", err)
	}

	result, err = h.Conn.Execute(sql)
	if err != nil {
		errStr := fmt.Sprintf("%s", err)
		log.Error(fmt.Sprintf("error : %v  ; sql :[%s] ; ", err, sql))
		if ok := sqlErrorReg.MatchString(errStr); !ok {
			result, err = h.tryExecute(sql)
		}
	}

	c <- true
	}(ch)

	select {
	case <-ch:
		e.Status = util.SQL_EXEC_FINISHED
	case <-time.After(util.DEFAULT_EXEC_SQL_TIME_OUT * time.Second):
		e.Status = util.SQL_EXEC_TIMEOUT
		err = NewError(ER_LOCK_WAIT_TIMEOUT, "sql查询超时")
		log.Error(fmt.Sprintf("执行sql [%s] 超时，KILL QUERY %d", sql, h.Conn.GetConnectionID()))
		_, err := h.killQuery(h.Conn.GetConnectionID())
		if err != nil {
			log.Error("kill query failed ", err)
		}
	}
	service.RecordEvent(e)
	return result, err
}

func (h *MysqlProxyHandler) resultFilter(res *Result, resp service.SqlParserResp) (err error) {
	tableColumn := make(map[string]string, 0)
	for _, v := range resp.Result.Table {
		v = strings.Replace(v, "`", "", -1)
		if _, ok := util.InformationSchema[v]; ok {
			continue
		}
		if db, ok := h.UserDbList.DBList[h.CurrentDb]; ok {
			if db.DbPermission.HaveAll == util.DB_PERMISSION_STATUS_ALLOW {
				continue
			}
			if table, ok := db.Tables[v]; ok {
				for _, colV := range table.Columns {
					tableColumn[colV.ColumnName] = "encode"
				}
			} else {
				log.Error(fmt.Sprintf("resultFilter: table：%s 无操作权限", v))
				return errors.New(fmt.Sprintf("[%s] no access", v))
			}
		}
	}
	fieldsMap := make(map[int]*Field, 0)
	for k, v := range res.Fields {
		fieldsMap[int(k)] = v
	}

	aliasMap := getColumnAliasMap(resp, tableColumn)
	columnLen := PutLengthEncodedInt(uint64(len(res.Fields)))
	data := make([]byte, 4, 1024)
	data = append(data, columnLen...)
	rowDatas := make([]RowData, 0)

	for _, v := range res.RowDatas {
		row, err := v.ParseText(res.Fields)
		if err != nil {
			log.Error("ParseText error ", err, " v : ", string(v))
		}
		rowdata := make(RowData, 0)
		for k, vv := range row {
			var name string
			var aliasName string
			if r, ok := fieldsMap[k]; ok {
				name = string(r.Name)
				aliasName = string(r.OrgName)

				if _, ok := tableColumn[name]; ok {
					rowdata = append(rowdata, PutLengthEncodedString([]byte(util.ResetSecretFieldValue(r.Type)))...)
				} else if _, ok := aliasMap[aliasName]; ok {
					rowdata = append(rowdata, PutLengthEncodedString([]byte(util.ResetSecretFieldValue(r.Type)))...)
				} else {
					var row []byte
					if vv == nil {
						row, _ = util.FormatDefaultValue(r.Type)
						rowdata = append(rowdata, PutLengthEncodedString([]byte(row))...)
					} else {
						row, _ = util.FormatValue(vv)
						rowdata = append(rowdata, PutLengthEncodedString([]byte(row))...)
					}
				}
			}
		}

		rowDatas = append(rowDatas, rowdata)
	}
	res.RowDatas = rowDatas

	return
}

func (h *MysqlProxyHandler) checkTablePermission(sql string) (r service.SqlParserResp, err error) {
	r, err = service.ParserSql(sql, h.User)
	if err != nil {
		log.Error(fmt.Sprintf("sql[%s]解析失败", sql))
		return r, err
	}
	for _, v := range r.Result.Table {
		v = strings.Replace(v, "`", "", -1)
		if _, ok := util.InformationSchema[v]; ok {
			continue
		}
	}
	return r, nil
}

func (h *MysqlProxyHandler) showDabases(sql string) (res *Result, err error) {
	dbs := make([]interface{}, 0, 0)
	if len(h.UserDbList.DBList) < 1 {
		return res, NewError(ER_UNKNOWN_ERROR, "无数据库操作权限")
	} else {
		for k, _ := range h.UserDbList.DBList {
			dbs = append(dbs, k)
		}
	}

	return h.buildSimpleShowResultset(dbs, "Database", sql)
}

func buildSchemataData(db string, encode string, defaultEncode string) (rowDatas RowData) {
	rowdata := make(RowData, 0)

	row, _ := util.FormatValue(db)
	rowdata = append(rowdata, PutLengthEncodedString([]byte(row))...)

	row, _ = util.FormatValue(encode)
	rowdata = append(rowdata, PutLengthEncodedString([]byte(row))...)

	row, _ = util.FormatValue("utf8_general_ci")
	rowdata = append(rowdata, PutLengthEncodedString([]byte(row))...)

	return rowdata
}

func (h *MysqlProxyHandler) selectSchemata(sql string) (res *Result, err error) {
	rowDatas := make([]RowData, 0)

	if len(h.UserDbList.DBList) < 1 {
		return res, NewError(ER_UNKNOWN_ERROR, "无数据库操作权限")
	} else {
		for k, _ := range h.UserDbList.DBList {
			// todo 根据数据库配置设置编码类型
			dbType := "utf8"
			rowDatas = append(rowDatas, buildSchemataData(k, dbType, "utf8_general_ci"))
		}
	}

	r := new(Resultset)
	r.RowDatas = rowDatas
	field := &Field{}
	field.Name = hack.Slice("Database")
	field.Charset = 33
	field.Type = MYSQL_TYPE_VAR_STRING

	r.Fields = []*Field{field}

	res = &Result{Resultset: r,}
	return res, nil
}

func (h *MysqlProxyHandler) showTables(sql string) (res *Result, err error) {
	res, err = h.Execute(sql)
	if len(h.UserDbList.DBList) < 1 {
		return res, NewError(ER_UNKNOWN_ERROR, "无数据库操作权限")
	} else {
		if t, err := h.UserDbList.DBList[h.CurrentDb]; err {
			if t.DbPermission.HaveAll == util.DB_PERMISSION_STATUS_ALLOW || h.UserDbList.Super {
				return res, nil
			}
			tables := make([]interface{}, 0, 0)
			for k, _ := range t.Tables {
				tables = append(tables, k)
			}
			return h.buildSimpleShowResultset(tables, fmt.Sprintf("Tables_in_%s", h.CurrentDb), sql)
		} else {
			return res, NewError(ER_UNKNOWN_ERROR, "无当前数据库操作权限")
		}
	}
}

func (h *MysqlProxyHandler) HandleFieldList(table string, fieldWildcard string) ([]*Field, error) {

	return nil, fmt.Errorf("not supported HandleFieldList now")
}

func (h *MysqlProxyHandler) HandleStmtPrepare(query string) (int, int, interface{}, error) {

	return 0, 0, nil, fmt.Errorf("not supported HandleFieldList now")
}

func (h *MysqlProxyHandler) HandleStmtExecute(context interface{}, query string, args []interface{}) (*Result, error) {

	return nil, fmt.Errorf("HandleStmtExecute not supported now")
}

func (h *MysqlProxyHandler) HandleStmtClose(context interface{}) error {

	return nil
}

func (h *MysqlProxyHandler) HandleOtherCommand(cmd byte, data []byte) error {
	return NewError(
		ER_UNKNOWN_ERROR,
		fmt.Sprintf("HandleOtherCommand command %d is not supported now", cmd),
	)
}

func (h *MysqlProxyHandler) InitGetUserDbList() error {
	userDbList, err := service.GetUserDbList(h.UserId, util.DB_TYPE_MYSQL)
	if err != nil {
		log.Error("service InitGetUserDbList error", err)
		return err
	}
	h.UserDbList = userDbList
	return nil
}

func (h *MysqlProxyHandler) GetUserDbList(userDbList interface{}) bool {
	return h.UserDbList.DBList == nil
}

func (h *MysqlProxyHandler) InitUserDbList(userDbList interface{}, user string) error {
	if h.UserDbList.DBList == nil {
		if dbList, ok := userDbList.(service.UserDbList); ok {
			h.UserDbList = dbList
			if h.Conn == nil {
				h.InitDb("")
			}
		} else {
			log.Error("InitUserDbList error ", userDbList)
			return NewError(ER_UNKNOWN_ERROR, "数据库信息不存在")
		}
	}
	return nil
}

func getColumnAliasMap(resp service.SqlParserResp, tableColumn map[string]string) (resultMap map[string]string) {
	aliasMap := make(map[string]string, 0)
	for k, v := range resp.Result.Alias {
		if k == "" {
			continue
		}
		aliasMap[filterAliasTableName(k)] = filterAliasTableName(v)
	}
	// Get alias field for filter
	resultMap = make(map[string]string, 0)
	for k, _ := range aliasMap {
		if r, ok := tableColumn[getColumnByAlias(k, aliasMap)]; ok {
			resultMap[k] = r
		}
	}
	return
}

func getColumnByAlias(name string, aliasMap map[string]string) (res string) {
	key := filterAliasTableName(name)
	if r, ok := aliasMap[key]; ok {
		return getColumnByAlias(r, aliasMap)
	} else {
		return key
	}
}

func filterAliasTableName(name string) (table string) {
	key := strings.TrimRight(strings.Replace(name, "'", "", -1), ")")
	keys := strings.Split(key, ".")

	if len(keys) >= 2 {
		return keys[1]
	} else {
		return keys[0]
	}
}

func parserLimit(sql string) (newSql string) {
	var (
		limit  = util.DEFAULT_SQL_LIMIT
		offset = 0
		page   int
		l      int
		o      int
		lE     error
		oE     error
	)

	par := strings.Split(sql, "limit")
	if limitReg.MatchString(sql) {
		if p := limitNumReg.FindStringSubmatch(sql); len(p) >= 2 {
			limitNum, err := strconv.Atoi(strings.Trim(par[1], ""))
			if err == nil && limitNum > limit {
				sql = fmt.Sprintf("%s limit %d ", par[0], limit)
			}
		} else if p := limitOffsetReg.FindStringSubmatch(sql); len(p) >= 5 {
			if strings.Trim(p[3], " ") == "offset" {
				l, lE = strconv.Atoi(strings.Trim(p[2], ""))
				o, oE = strconv.Atoi(strings.Trim(p[4], ""))
			} else {
				l, lE = strconv.Atoi(strings.Trim(p[4], ""))
				o, oE = strconv.Atoi(strings.Trim(p[2], ""))
			}
			if l > limit {
				if o != 0 && lE == nil && oE == nil {
					page = o / l
				}
				offset = page * limit
				sql = fmt.Sprintf("%s limit %d,%d ", par[0], offset, limit)
			}
		}
	} else {
		sql = fmt.Sprintf("%s limit %d ", par[0], limit)
	}
	newSql = sql
	return newSql
}

var checkConnIsRunning bool
var checkConnStopChan chan bool
var checkConnStop bool

func (h *MysqlProxyHandler) CheckMysqlConn() {
	if checkConnIsRunning == true {
		return
	}
	checkConnStop = false
	sleepTime := util.DEFAULT_PING_INTERVAL * time.Second
	go func() {
		checkConnIsRunning = true
		for checkConnStop == false {
			//log.Info("check db Connect, db pint test")

			if h.Conn == nil {
				if _, err := h.InitDb(""); err != nil {
					checkConnIsRunning = false
					return
				}
			}
			if err := h.Conn.Ping(); err != nil {
				log.Error(fmt.Sprintf("user: %s , db Conn error : %v ,try reconnect . ", h.User, err))
				if _, err := h.InitDb(h.CurrentDb); err != nil {
					log.Error(fmt.Sprintf("user: %s , db Conn error : %v ,try reconnect failed .", h.User, err))
					sleepTime = 1 * time.Second
				} else {
					sleepTime = util.DEFAULT_PING_INTERVAL * time.Second
				}
			}
			time.Sleep(sleepTime)
		}
	}()

	checkConnStopChan = make(chan bool, 1)
	switch {
	case <-checkConnStopChan:
		checkConnStop = true
		fmt.Println("closet ping check")
	}
}

func (h *MysqlProxyHandler) Destory() {
	log.Info("Destory , user : ", h.User)
	checkConnStopChan <- true
	h.Conn.Close()
	h = nil
}
