package provider

import (
	"github.com/siddontang/go/sync2"
	"net"
)

var BaseConnID uint32 = 10000

type Conn struct {
	Conn net.Conn

	ServerConf     *Server
	capability     uint32
	authPluginName string
	ConnectionID   uint32
	status         uint16
	Salt           []byte

	user                string
	password            string
	cachingSha2FullAuth bool

	H Handler

	Closed sync2.AtomicBool
	//DbList *map[string]common.Database

	RequestInfo Request
}

type Request struct {
	Id           int32
	CreateTime   int
	UpdateTime   int
	User         string
	Method       string
	ClientCode   string
	ResponseCode string
	ValidateCode string
	I            int
	C            string
	Salt         string
}
