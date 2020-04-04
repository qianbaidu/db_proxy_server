package proxy

import (
	"github.com/prometheus/common/log"
	"github.com/qianbaidu/db-proxy/proxy/mysql/provider"
	mysqlUtil "github.com/qianbaidu/db-proxy/util/mysql"
	"github.com/siddontang/go-mysql/packet"
	"net"
	"sync/atomic"
)

func NewMysqlProxyConn(conn net.Conn, h provider.Handler) (*provider.Conn, error) {
	p := provider.NewInMemoryProvider()
	var (
		packetConn    *packet.Conn
		defaultServer = provider.NewDefaultServer()
	)

	packetConn = packet.NewConn(conn)

	c := &provider.Conn{
		Conn:               packetConn,
		ServerConf:         defaultServer,
		CredentialProvider: p,
		H:                  h,
		ConnectionID:       atomic.AddUint32(&provider.BaseConnID, 1),
		Stmts:              make(map[uint32]*provider.Stmt),
		Salt:               mysqlUtil.RandomBuf(20),
	}

	c.Closed.Set(false)

	if err := c.Handshake(); err != nil {
		log.Error("Handshake error , ", err)
		c.Close()
		return nil, err
	}

	return c, nil
}
