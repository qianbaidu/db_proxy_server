package provider

import (
	"crypto/tls"
	"sync"
)

var defaultServer = NewDefaultServer()

type Server struct {
	serverVersion     string // e.g. "8.0.12"
	protocolVersion   int    // minimal 10
	capability        uint32 // server capability flag
	collationId       uint8
	defaultAuthMethod string // default authentication method, 'mysql_native_password'
	pubKey            []byte
	tlsConfig         *tls.Config
	cacheShaPassword  *sync.Map // 'user@host' -> SHA256(SHA256(PASSWORD))
}

func NewDefaultServer() *Server {
	return &Server{}
}
