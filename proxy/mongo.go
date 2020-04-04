package proxy

import (
	"github.com/qianbaidu/db-proxy/proxy/mongo/provider"
	"net"
	"sync/atomic"
)

func NewMongoProxyConn(conn net.Conn, h provider.Handler) (*provider.Conn, error) {
	c := &provider.Conn{
		Conn: conn,

		H:            h,
		ConnectionID: atomic.AddUint32(&provider.BaseConnID, 1),
	}

	return c, nil
}
