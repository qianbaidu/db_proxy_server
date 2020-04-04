package mongo

import (
	"github.com/globalsign/mgo/bson"
	p "github.com/qianbaidu/db-proxy/proxy/mongo/provider"
	"github.com/qianbaidu/db-proxy/util"
	"github.com/prometheus/common/log"
	"strings"
)

func (m *MongoProxyHandler) DispatchQuery(cmd string, db string, query bson.D, fields []byte) (doc p.SimpleBSON, err error) {
	if ok := m.checkLogin(); !ok {
		return simpleBSONConvert(p.Unauthorized())
	}

	table := "$cmd"
	if r, ok := query[0].Value.(string); ok && len(r) > 0 {
		table = r
	}
	if ok := m.checkAuth(db, table); !ok {
		return m.permissionDenied(db, table)
	}

	switch strings.ToLower(cmd) {
	case "listdatabases":
		doc, err = m.listDatabase()
	case "listCollections", "listcollections":
		doc, err = m.listCollections(cmd, db, query)
	case "find":
		doc, err = m.find(cmd, db, query)
	case "explain":
		doc, err = m.explain(cmd, db, query, fields)
	case "distinct":
		doc, err = m.distinct(cmd, db, query)
	case "count":
		doc, err = m.count(db, query)
	case "listindexes":
		doc, err = m.listIndexes(cmd, db, query)
	case "aggregate":
		doc, err = m.aggregate(cmd, db, query)
	case "mapreduce":
		doc, err = m.mapreduce(cmd, db, query)

	//case
	//	"getpreverror", "getmore", "plancachelistqueryshapes", "plancachelistfilters", "plancachelistplans",
	//	"usersinfo", "getshardmap", "getshardversion", "listshards",
	//	"getparameter", "listcommands", "cursorinfo", "hostinfo",
	//	"availablequeryoptions", "connpoolstats", "connectionstatus", "datasize", "dbhash", "diaglogging",
	//	"driveroidtest", "features", "isself", "netstat", "profile", "shardconnpoolstats", "top",
	//	"validate":

	default:
		log.Info("cmd ", cmd, ", query : ", query.Map())
		util.JsonPrint(query)
		doc, err = m.execCmd(cmd, db)
	}
	return doc, err
}
