package proxy

import (
	"github.com/prometheus/common/log"
	"github.com/qianbaidu/db-proxy/proxy/mongo"
	"github.com/qianbaidu/db-proxy/proxy/mysql"
	"net"
)

func RunMysql(c net.Conn) {
	conn, _ := NewMysqlProxyConn(c, &mysql.MysqlProxyHandler{})
	for {
		err := mysql.HandleProxy(conn)
		if err != nil {
			log.Error("HandleMysqlProxy : ", err)
			conn.H.Destory()
			//c.Close()
			break
		}
	}
}

func RunMongo(c net.Conn) {
	conn, _ := NewMongoProxyConn(c, &mongo.MongoProxyHandler{})
	for {
		err := mongo.HandleProxy(conn)
		if err != nil {
			log.Error("HandleMongoProxy : ", err)
			c.Close()
			break
		}
	}
}
