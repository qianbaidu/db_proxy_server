package mongo

import (
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	"github.com/qianbaidu/db-proxy/models"
	p "github.com/qianbaidu/db-proxy/proxy/mongo/provider"
	"github.com/qianbaidu/db-proxy/service"
	"github.com/qianbaidu/db-proxy/util"
	"reflect"
	"strings"
	"time"
)

type MongoProxyHandler struct {
	User       string
	UserId     int64
	IsSupper   bool
	AuthData   string
	UserDbList service.UserDbList
	CurrentDb  string
	Conn       *mgo.Session
}

func (m *MongoProxyHandler) isMaster(db string) (doc p.SimpleBSON, err error) {
	return simpleBSONConvert(p.NewIsMaster())
}

func (m *MongoProxyHandler) getLastError() (doc p.SimpleBSON, err error) {
	return simpleBSONConvert(p.NewLastError())
}

func simpleBSONConvert(d interface{}) (p.SimpleBSON, error) {
	doc, err := p.SimpleBSONConvert(d)
	if err != nil {
		log.Error("SimpleBSONConvert error, ", err, " doc : ", d)
	}
	return doc, nil
}

func (m *MongoProxyHandler) getNonce() (doc p.SimpleBSON, err error) {
	// todo 自动生成
	msg := p.Nonce{"6d32d13b13436425"}
	return simpleBSONConvert(msg)
}

func (m *MongoProxyHandler) getBuildInfo() (doc p.SimpleBSON, err error) {
	return simpleBSONConvert(p.NewBuildInfo())
}

func (m *MongoProxyHandler) ping() (doc p.SimpleBSON, err error) {
	msg := p.Ping{1}
	return simpleBSONConvert(msg)
}

func (m *MongoProxyHandler) defaultRetrun() (doc p.SimpleBSON, err error) {
	return simpleBSONConvert(p.NewDefaultServerReturn())
}

func (m *MongoProxyHandler) listDatabase() (doc p.SimpleBSON, err error) {
	if len(m.UserDbList.DBList) < 1 {
		m.GetUserDbList(m.User)
	}
	dbList := make([]p.Database, 0)
	log.Info("\n\n---listDatabase :  m.UserDbList.DBList length : ", len(m.UserDbList.DBList))
	for db, _ := range m.UserDbList.DBList {
		dbList = append(dbList, p.Database{
			Name:       db,
			Empty:      false,
			SizeOnDisk: 1000,
		})
	}

	msg := p.DbListResult{
		Databases: dbList,
		Ok:        1,
		TotalSize: 1000,
	}

	return simpleBSONConvert(msg)
}

func (m *MongoProxyHandler) getUserTablesByDb(db string) (tableMap map[string]bool) {
	tableMap = make(map[string]bool)
	if db, ok := m.UserDbList.DBList[db]; ok {
		if db.HaveAll == util.DB_PERMISSION_STATUS_ALLOW {
			tableMap["HasAll"] = true
		}
		for t, _ := range db.Tables {
			tableMap[t] = true
		}
	}
	return
}

func (m *MongoProxyHandler) listCollections(cmd interface{}, db string, query bson.D) (doc p.SimpleBSON, err error) {
	var res p.ListCollections
	queryCmd := p.BindListCollectionOp(query)
	log.Debug("queryCmd : ", queryCmd)
	err = m.runCmd(queryCmd, &res, db, "$cmd")
	if err != nil {
		res.Ok = 0
		res.Errmsg = fmt.Sprintf("error : %s", err)
		return simpleBSONConvert(res)
	}

	tableMap := m.getUserTablesByDb(db)
	if _, ok := tableMap["HasAll"]; ok || m.IsSupper == true {
		return simpleBSONConvert(res)
	}

	var result []p.CollectionResult
	for _, v := range res.Cursor.FirstBatch {
		if _, ok := tableMap[v.Name]; ok {
			result = append(result, v)
		}
	}
	log.Debug("CollectionResult : ", result)
	res.Cursor.FirstBatch = result
	return simpleBSONConvert(res)
}

func (m *MongoProxyHandler) getLog() (doc p.SimpleBSON, err error) {
	msg := p.GetLogResult{
		Names: []string{"global", "startupWarnings"},
		Ok:    1,
	}
	return simpleBSONConvert(msg)
}

func (m *MongoProxyHandler) notSupported(cmd string) (doc p.SimpleBSON, err error) {
	errBSON := bson.D{
		{"ok", 0},
		{"errmsg", fmt.Sprintf("not supported. (cmd: %s)", cmd)},
		{"code", 0},
	}
	return p.SimpleBSONConvert(errBSON)
}

func (m *MongoProxyHandler) permissionDenied(db, table string) (doc p.SimpleBSON, err error) {
	errBSON := bson.D{
		{"ok", 0},
		{"errmsg", fmt.Sprintf("Permission denied. (db: %s, table: %s)", db, table)},
		{"code", 0},
	}
	log.Error(errBSON)
	return p.SimpleBSONConvert(errBSON)
}

func (m *MongoProxyHandler) listIndexes(cmd string, db string, query bson.D) (doc p.SimpleBSON, err error) {
	c := p.BindListIndexOp(query)

	if ok := m.checkAuth(db, c.ListIndexes); !ok {
		return m.permissionDenied(db, c.ListIndexes)
	}

	var res interface{}
	err = m.runCmd(c, &res, db, c.ListIndexes)
	return simpleBSONConvert(res)
}

func (m *MongoProxyHandler) mapreduce(cmd string, db string, query bson.D) (doc p.SimpleBSON, err error) {
	c := p.BindMapReduceOp(query)

	if ok := m.checkAuth(db, c.MapReduce); !ok {
		return m.permissionDenied(db, c.MapReduce)
	}

	var res interface{}
	err = m.runCmd(c, &res, db, c.MapReduce)
	return simpleBSONConvert(res)
}

func (m *MongoProxyHandler) aggregate(cmd string, db string, query bson.D) (doc p.SimpleBSON, err error) {
	c := p.BindAggregateOp(query)

	if ok := m.checkAuth(db, c.Aggregate); !ok {
		return m.permissionDenied(db, c.Aggregate)
	}

	var res interface{}
	err = m.runCmd(c, &res, db, c.Aggregate)
	return simpleBSONConvert(res)
}

func (m *MongoProxyHandler) saslStart(cmd p.SaslCmd) (doc p.SimpleBSON, err error) {
	payloadParms := strings.Split(cmd.Payload, ",")
	r := p.SaslResult{
		Ok:             true,
		NotOk:          false,
		Done:           false,
		ConversationId: 1, // todo
	}
	if len(payloadParms) >= 4 {
		switch cmd.Mechanism {
		case "SCRAM-SHA-1":
			randCode := fmt.Sprintf("%s%s", payloadParms[3][2:], util.RandString(8))
			payload := fmt.Sprintf("r=%s,s=%s,i=%d", randCode, p.SaltS, p.SaltI)
			r.Payload = []byte(payload)

			loginRequest := p.UserLoginReuqest{
				User:       payloadParms[2][2:],
				Method:     cmd.Mechanism,
				ClientCode: payloadParms[3][2:],
				RandCode:   randCode,
				I:          p.SaltI,
				Salt:       p.SaltS,
				CreateTime: time.Now().Unix(),
				UpdateTime: time.Now().Unix(),
			}
			m.User = payloadParms[2][2:]

			p.UserLoginData[randCode] = loginRequest
			log.Info(fmt.Sprintf("user: %s start login.", payloadParms[2][2:]))
		default:
			r.ErrMsg = fmt.Sprintf("not supported cmd.Mechanism %s ", cmd.Mechanism)
		}
	} else {
		r.ErrMsg = "params error"
	}
	if len(r.ErrMsg) > 0 {
		r.Ok = false
		r.NotOk = true
		log.Error("saslStart error : ", r.ErrMsg)
	}
	return simpleBSONConvert(r)
}

func (m *MongoProxyHandler) saslContinue(cmd p.SaslCmd) (doc p.SimpleBSON, err error) {
	payloadParms := strings.Split(cmd.Payload, ",")
	r := p.SaslResult{
		Ok:             true,
		NotOk:          false,
		Done:           true,
		ConversationId: 1, // todo
	}
	if len(payloadParms) >= 3 {
		if u, ok := p.UserLoginData[payloadParms[1][2:]]; ok {
			u.UpdateTime = time.Now().Unix()
			u.C = payloadParms[0][2:]
			u.ValidateCode = payloadParms[2][2:]
			authMsg := fmt.Sprintf("n=%s,r=%s,r=%s,s=%s,i=%d,c=%s,r=%s",
				u.User, u.ClientCode, u.RandCode, u.Salt, u.I, u.C, u.RandCode)
			log.Info("authMsg", authMsg)
			// step1 : hash生成密码
			//hashPass, err := m.getUserPassHash(u.User)
			user, err := service.GetUser(u.User)
			//hashPass := p.PassHash("test", "test")

			if err != nil {
				log.Error("get user hashPassword error : ", err)
				r.ErrMsg = "Authentication failed."
				r.Ok = false
				r.NotOk = true
				return simpleBSONConvert(r)
			} else {
				// step2 : 盐值密码计算
				log.Info("user.MongodbPassword : ", user.MongodbPassword)
				saltedPass := p.SaltPassword(user.MongodbPassword, p.SaltB, u.I)
				log.Info("getUserPassHash", saltedPass)
				// step3 : server端签名字符串
				sign := p.ServerSignature(saltedPass, []byte(authMsg))

				payload := fmt.Sprintf("v=%s", sign)
				log.Info("payload:", payload)
				r.Payload = []byte(payload)

				// 加载用户权限数据库列表
				m.UserId = user.Id
				err = m.GetUserDbList(u.User)
				if err != nil {
					r.ErrMsg = "Authentication failed. user does not have [dataquery] system permission."
				}

			}
		} else {
			r.ErrMsg = "Authentication failed. user info not exists "
		}
	} else {
		r.Payload = []byte("")
	}
	if len(r.ErrMsg) > 0 {
		r.Ok = false
		log.Error("saslContinue error : ", r.ErrMsg)
	}
	return simpleBSONConvert(r)
}

func (m *MongoProxyHandler) getMore(ns string, query bson.D) (doc p.SimpleBSON, err error) {
	c := p.BindGetMoreOp(query)
	var res interface{}
	err = m.runCmd(c, &res, ns, "$cmd")
	return simpleBSONConvert(res)
}

func (m *MongoProxyHandler) checkLogin() bool {
	if len(m.User) < 1 {
		return false
	}
	return true
}

func (m *MongoProxyHandler) checkAuth(db string, table string) bool {
	if ok := m.checkLogin(); !ok {
		return false
	}
	if m.IsSupper == true {
		return true
	}
	if db, ok := m.UserDbList.DBList[db]; !ok {
		return false
	} else {
		if _, ok := db.Tables[table]; !ok && table != "$cmd" {
			return false
		}
	}
	return true
}

func (m *MongoProxyHandler) getBdId() (int64, error) {
	if db, ok := m.UserDbList.DBList[m.CurrentDb]; ok {
		return db.DbList.Id, nil
	} else {
		log.Error("get current db id error ")
		return 0, errors.New("get current db id error")
	}
}

func (m *MongoProxyHandler) NewEventLog(sql string) (e *models.EventLog) {
	dbId, _ := m.getBdId()
	e = &models.EventLog{
		Sql:            sql,
		Status:         0,
		CreateTime:     time.Now().Format("2006-01-02 15:04:05"),
		DbId:           dbId,
		UserId:         m.UserId,
		UpdateDatatime: time.Now().Format("2006-01-02 15:04:05"),
	}
	return e
}

func (m *MongoProxyHandler) runCmd(cmd interface{}, res interface{}, db string, table string) error {
	db = strings.Split(db, ".")[0]
	cmdJson, _ := json.Marshal(cmd)
	cmdStr := string(cmdJson)
	eventLog := m.NewEventLog(cmdStr)

	conn, err := m.dbConnSession(db)
	if err != nil {
		err = errors.New(fmt.Sprintf("init db error : %v ", err))
		log.Error(err)
		return err
	}
	dbCon := conn.Db
	err = dbCon.Run(cmd, res)
	eventLog.Status = util.SQL_EXEC_FINISHED
	if err != nil {
		if err == p.RunCmdTimeoutError {
			eventLog.Status = util.SQL_EXEC_TIMEOUT
		}
		log.Error(fmt.Sprintf("runCmd: error [%s], cmd : [%s]", err, cmdStr))
		err = errors.New(fmt.Sprintf("error: %s", err))
		return err
	}
	eventLog.UpdateDatatime = time.Now().Format("2006-01-02 15:04:05")
	service.RecordEvent(eventLog)
	return nil
}

func (m *MongoProxyHandler) find(cmd string, ns string, query bson.D) (doc p.SimpleBSON, err error) {
	findCmd := p.BindFindOp(query)

	if ok := m.checkAuth(ns, findCmd.Find); !ok {
		return m.permissionDenied(ns, findCmd.Find)
	}

	var res interface{}
	err = m.runCmd(findCmd, &res, ns, findCmd.Find)

	return simpleBSONConvert(res)
}

func (m *MongoProxyHandler) collStats(cmd string, ns string, query bson.D) (doc p.SimpleBSON, err error) {
	c := p.BindCollStatsOp(query)

	var res interface{}
	err = m.runCmd(c, &res, ns, c.CollStats)

	return simpleBSONConvert(res)
}

func (m *MongoProxyHandler) dbStats(cmd string, db string, query bson.D) (doc p.SimpleBSON, err error) {
	if len(m.User) < 1 || len(db) < 1 {
		return simpleBSONConvert(p.NewDbStatsResult())
	}

	c := p.BindDbStatsOp(query)

	if ok := m.checkAuth(db, ""); !ok {
		return m.permissionDenied(db, "")
	}

	var res interface{}
	err = m.runCmd(c, &res, db, "")

	return simpleBSONConvert(res)
}

// doc https://docs.mongodb.com/manual/reference/method/db.collection.explain/index.html#db.collection.explain
func (m *MongoProxyHandler) explain(cmd string, ns string, query bson.D, fields []byte) (doc p.SimpleBSON, err error) {
	var (
		table = "$cmd"
		res   interface{}
	)
	var e p.Explain
	err = bson.Unmarshal(fields, &e)
	if err != nil {
		log.Error("Unmarshal fields error ", err)
	}

	if t, ok := e.Explain[0].Value.(string); ok {
		table = t
	}
	err = m.runCmd(e, &res, ns, table)
	if err != nil {
		log.Error("run explain cmd error ", err)
	}

	return simpleBSONConvert(res)
}

func (m *MongoProxyHandler) distinct(cmd, db string, query bson.D) (doc p.SimpleBSON, err error) {
	c := p.BindDistinctOp(query)

	if ok := m.checkAuth(db, c.Distinct); !ok {
		return m.permissionDenied(db, c.Distinct)
	}

	var res interface{}
	err = m.runCmd(c, &res, db, c.Distinct)

	return simpleBSONConvert(res)
}

func (m *MongoProxyHandler) count(db string, query bson.D) (doc p.SimpleBSON, err error) {
	c := p.BindCountOp(query)

	if ok := m.checkAuth(db, c.Count); !ok {
		return m.permissionDenied(db, c.Count)
	}

	var res interface{}
	err = m.runCmd(c, &res, db, c.Count)

	return simpleBSONConvert(res)
}

func errorDoc(err error) (p.SimpleBSON) {
	var errorInfo string
	if err != nil && len(err.Error()) < 1 {
		errorInfo = err.Error()
	} else {
		errorInfo = "error"
	}
	errBSON := bson.D{{"ok", 0}, {"errmsg", errorInfo}}
	doc, _ := p.SimpleBSONConvert(errBSON)
	return doc
}

func (m *MongoProxyHandler) HandleQuery(msg p.Message) (r *p.ReplyMessage, err error) {
	var doc p.SimpleBSON
	r = &p.ReplyMessage{
		Flags:          8,
		CursorId:       0,
		StartingFrom:   0,
		NumberReturned: 1,
	}

	switch mm := msg.(type) {
	case *p.QueryMessage:
		query, err := mm.Query.ToBSOND()
		if err != nil {
			doc = errorDoc(errors.New(fmt.Sprintf("QueryMessage to bson error : ", err)))
			log.Error(doc)
		} else {
			doc, err = m.dispatch(mm.Namespace, query[0].Name, query, mm.Query.BSON)
		}
	default:
		doc = errorDoc(errors.New(fmt.Sprintf("QueryMessage: not supported. default %v", reflect.TypeOf(msg))))
		log.Error(err)
	}
	r.Docs = []p.SimpleBSON{doc}
	return r, err
}

// doc : https://docs.mongodb.com/master/reference/command/
// Reject all requests except query.
func (m *MongoProxyHandler) dispatch(db string, cmd string, query bson.D, fields []byte) (doc p.SimpleBSON, err error) {
	queryCmd := strings.ToLower(cmd)
	log.Info("dispatch cmd : ", queryCmd, " cmd:", cmd, " query : ", query.Map())
	switch queryCmd {
	// part 1 : 伪handler mongo server ，响应客户端连接请求
	case "getnonce":
		doc, err = m.getNonce()
	case "getlog":
		doc, err = m.getLog()
	case "ping", "logout":
		doc, err = m.ping()
	case "dbstats":
		doc, err = m.dbStats(cmd, db, query)
	case "ismaster":
		doc, err = m.isMaster(db)
	case "getlasterror":
		doc, err = m.getLastError()
	case "whatsmyuri", "replsetgetstatus", "getcmdlineopts", "serverstatus": // 用户登录前需要获取服务端信息，通用伪造返回
		doc, err = m.defaultRetrun()
	case "buildinfo":
		doc, err = m.getBuildInfo()
	case "collstats":
		doc, err = m.collStats(cmd, db, query)
		//doc, err = m.execCmd(query[0], db)
	case "getmore":
		doc, err = m.getMore(db, query)

	// part 2 : 登录后请求代理转发到mongo连接请求中处理再返回
	case "saslstart", "saslcontinue":
		saslCmd := p.SaslCmd{}
		err := bson.Unmarshal(fields, &saslCmd)
		if err != nil {
			log.Error("CommandMessage,bson.Unmarshal error : ", err, " queryBson : ", fields)
		}
		log.Info("saslCmd : ", saslCmd)
		if queryCmd == "saslstart" {
			doc, err = m.saslStart(saslCmd)
		} else {
			doc, err = m.saslContinue(saslCmd)
		}

	// part 3 : 需登录、执行权限; 及其他未实现功能
	case "update", "delete", "insert", "dropdatabase", "findandmodify", "createindexes", "deleteindexes", "reindex", "drop":
		log.Info("notSupported : ", cmd)
		doc, err = m.notSupported(cmd)
	default:
		doc, err = m.DispatchQuery(cmd, db, query, fields)
	}

	return doc, err
}

func (m *MongoProxyHandler) HandleOpCommand(msg p.Message) (r *p.CommandReplyMessage, err error) {
	var doc p.SimpleBSON
	r = &p.CommandReplyMessage{
		Metadata:   p.SimpleBSONEmpty(),
		OutputDocs: []p.SimpleBSON{},
	}

	switch mm := msg.(type) {
	case *p.CommandMessage:
		query, err := mm.CommandArgs.ToBSOND()
		if err != nil {
			doc = errorDoc(errors.New(fmt.Sprintf("CommandMessage to bson error : ", err)))
			log.Error(doc)
		} else {
			doc, err = m.dispatch(mm.DB, mm.CmdName, query, mm.CommandArgs.BSON)
		}
	}

	r.CommandReply = doc
	return r, err
}

func (m *MongoProxyHandler) execCmd(cmd interface{}, db string) (doc p.SimpleBSON, err error) {
	var res interface{}
	err = m.runCmd(cmd, &res, db, "$cmd")
	return simpleBSONConvert(res)
}

func (m *MongoProxyHandler) dbConnSession(db string) (mgoConnect, error) {
	return m.initDbConn(db)
}

func (m *MongoProxyHandler) HandleOpMsg(msg p.Message) (r *p.MessageMessage, err error) {
	switch mm := msg.(type) {
	case *p.MessageMessage:
		j, e := json.Marshal(mm)
		log.Info("HandleOpMsg : ", string(j), e)

		errBSON := bson.D{{"ok", 1}}
		doc, _ := p.SimpleBSONConvert(errBSON)
		r := &p.MessageMessage{
			FlagBits: 0,
			Sections: []p.MessageMessageSection{
				&p.BodySection{
					doc,
				},
			},
		}
		return r, nil
	default:
		log.Error(fmt.Sprintf("HandleOpMsg :not supported %v", reflect.TypeOf(msg)))
		return r, errors.New(fmt.Sprintf("HandleOpMsg :not supported %v", reflect.TypeOf(msg)))

	}
	return r, errors.New(fmt.Sprintf("HandleOpMsg :not supported %v", reflect.TypeOf(msg)))
}

func (m *MongoProxyHandler) GetUserDbList(user interface{}) error {
	if userDbList, err := service.GetUserDbList(int64(1), util.DB_TYPE_MONGO); err != nil {
		log.Error("MongoProxyHandler GetUserDbList error ", err)
		return err
	} else {
		m.UserDbList = userDbList
		m.IsSupper = userDbList.Super
	}
	return err
}

func (m *MongoProxyHandler) initDbConn(dbName string) (mgoConnect, error) {
	m.CurrentDb = dbName
	dbName = strings.Split(dbName, ".")[0]
	if db, ok := m.UserDbList.DBList[dbName]; ok {
		if conn, ok := GetMongoConnect(db.DbName); !ok {
			return InitMongoConnect(db.DbList)
		} else {
			return conn, nil
		}
	} else {
		log.Error(fmt.Sprintf("No database[%s] permission", db))
	}

	return mgoConnect{}, errors.New("database not exists.")
}
