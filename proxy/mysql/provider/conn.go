package provider

import (
	"net"
	"sync/atomic"

	. "github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/packet"
	"github.com/siddontang/go/sync2"

	//"github.com/qianbaidu/db-proxy/proxy/common"
)

/*
   Conn acts like a MySQL server connection, you can use MySQL client to communicate with it.
*/
type Conn struct {
	*packet.Conn

	ServerConf     *Server
	capability     uint32
	authPluginName string
	ConnectionID   uint32
	status         uint16
	Salt           []byte // should be 8 + 12 for auth-plugin-data-part-1 and auth-plugin-data-part-2

	CredentialProvider  CredentialProvider
	user                string
	password            string
	cachingSha2FullAuth bool

	H Handler

	Stmts  map[uint32]*Stmt
	StmtID uint32

	Closed sync2.AtomicBool
	//DbList *map[string]common.Database
}

type DatabasesInfo struct {
	Addr     string
	User     string
	Password string
	DbName   string
	Charset  string
}

var BaseConnID uint32 = 10000

// NewConn: create connection with default server settings
func NewConn(conn net.Conn, user string, password string, h Handler) (*Conn, error) {
	p := NewInMemoryProvider()
	p.AddUser(user, password)
	salt, _ := RandomBuf(20)

	var packetConn *packet.Conn
	if defaultServer.tlsConfig != nil {
		packetConn = packet.NewTLSConn(conn)
	} else {
		packetConn = packet.NewConn(conn)
	}

	c := &Conn{
		Conn:               packetConn,
		ServerConf:         defaultServer,
		CredentialProvider: p,
		H:                  h,
		ConnectionID:       atomic.AddUint32(&BaseConnID, 1),
		Stmts:              make(map[uint32]*Stmt),
		Salt:               salt,
	}
	c.SetClose(false)

	if err := c.Handshake(); err != nil {
		c.SetClose(true)
		return nil, err
	}

	return c, nil
}

// NewCustomizedConn: create connection with customized server settings
func NewCustomizedConn(conn net.Conn, serverConf *Server, p CredentialProvider, h Handler) (*Conn, error) {
	var packetConn *packet.Conn
	if serverConf.tlsConfig != nil {
		packetConn = packet.NewTLSConn(conn)
	} else {
		packetConn = packet.NewConn(conn)
	}

	salt, _ := RandomBuf(20)
	c := &Conn{
		Conn:               packetConn,
		ServerConf:         serverConf,
		CredentialProvider: p,
		H:                  h,
		ConnectionID:       atomic.AddUint32(&BaseConnID, 1),
		Stmts:              make(map[uint32]*Stmt),
		Salt:               salt,
	}
	c.SetClose(false)

	if err := c.Handshake(); err != nil {
		c.SetClose(true)
		return nil, err
	}

	return c, nil
}

func (c *Conn) Handshake() error {
	if err := c.writeInitialHandshake(); err != nil {
		return err
	}

	if err := c.readHandshakeResponse(); err != nil {
		if err == ErrAccessDenied {
			err = NewDefaultError(ER_ACCESS_DENIED_ERROR, c.user, c.LocalAddr().String(), "Yes")
		}
		c.writeError(err)
		return err
	}

	if err := c.writeOK(nil); err != nil {
		return err
	}

	c.ResetSequence()

	return nil
}

func (c *Conn) SetClose(b bool) {
	c.Closed.Set(b)
	c.Conn.Close()
}

func (c *Conn) GetClosed() bool {
	return c.Closed.Get()
}

func (c *Conn) GetUser() string {
	return c.user
}

func (c *Conn) GetConnectionID() uint32 {
	return c.ConnectionID
}

func (c *Conn) IsAutoCommit() bool {
	return c.status&SERVER_STATUS_AUTOCOMMIT > 0
}

func (c *Conn) IsInTransaction() bool {
	return c.status&SERVER_STATUS_IN_TRANS > 0
}

func (c *Conn) SetInTransaction() {
	c.status |= SERVER_STATUS_IN_TRANS
}

func (c *Conn) ClearInTransaction() {
	c.status &= ^SERVER_STATUS_IN_TRANS
}
