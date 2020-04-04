package mongo

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/prometheus/common/log"
	"github.com/qianbaidu/db-proxy/models"
	"github.com/qianbaidu/db-proxy/util"
	"strings"
	"sync"
	"time"
)

var (
	mongoConnPool map[string]mgoConnect
	m             sync.RWMutex
)

type mgoConnect struct {
	Conn *mgo.Session
	Db   *mgo.Database
}

func init() {
	mongoConnPool = make(map[string]mgoConnect, 1024)
}

func GetMongoConnect(dbname string) (conn mgoConnect, ok bool) {
	m.Lock()
	defer m.Unlock()
	conn, ok = mongoConnPool[dbname]
	return
}

func InitMongoConnect(db models.DbList) (c mgoConnect, err error) {
	ip := strings.Split(db.Ip, ":")[0]
	dbUrl := fmt.Sprintf("%s:%s@%s:%s/%s", db.Username, db.Password, ip, db.Port, db.DbName)
	log.Info("connect mongo dns : ", dbUrl)

	conn, err := mgo.Dial(dbUrl)
	conn.SetMode(mgo.Eventual, true)
	conn.SetPoolTimeout(util.DEFAULT_MONGO_CONN_TIME_OUT * time.Second)
	conn.SetPoolLimit(util.DEFAULT_MONGO_CONN_POOL_SIZE)

	if err != nil {
		log.Error("mongo db connect error : ", err, " . connect url : ", dbUrl)
		return
	} else {
		c = mgoConnect{
			Conn: conn,
			Db:   conn.DB(db.DbName),
		}
		log.Info("-----------------------init db conn dbName : ", db.DbName)
		m.Lock()
		defer m.Unlock()
		mongoConnPool[db.Name] = c
		return c, nil
	}
}

