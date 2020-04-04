package mongo

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/prometheus/common/log"
	p "github.com/qianbaidu/db-proxy/proxy/mongo/provider"

	"gopkg.in/mgo.v2/bson"
	"io"
)

var (
	mconn *mgo.Session
	err   error
)

type isMasterResult struct {
	IsMaster       bool
	Secondary      bool
	Primary        string
	Hosts          []string
	Passives       []string
	Tags           bson.D
	Msg            string
	SetName        string `bson:"setName"`
	MaxWireVersion int    `bson:"maxWireVersion"`
}

func NewIsMaster() {
	return
}

type respBson struct {
	ConversationId int    `bson:"conversationId"`
	Done           bool   `bson:"done"`
	Payload        string `bson:"payload"`
	Ok             int    `bson:"ok"`
	//NotOk          bool   `bson:"code"`
	//ErrMsg         string
}

type authBson struct {
	SaslStart      int    `bson:"saslStart"`
	Mechanism      string `bson:"mechanism"`
	Payload        string `bson:"payload"`
	AutoAuthorize  int    `bson:"autoAuthorize"`
	Ok             bool   `bson:"ok"`
	Done           bool   `bson:"done"`
	ConversationId int    `bson:"conversationId"`
}

type find struct {
	Find        string      `bson:"find"`
	Filter      interface{} `bson:"filter"`
	Skip        int         `bson:"skip"`
	Limit       int         `bson:"limit"`
	BatchSize   int         `bson:"batchSize"`
	SingleBatch int         `bson:"singleBatch"`
}

func HandleProxy(c *p.Conn) (err error) {
	m, err := p.ReadMessage(c.Conn)
	if err != nil {
		if err == io.EOF {
			log.Info("io.EOF exit")
			return err
		}
		log.Info("ReadMessage error ", err)
		return p.NewStackErrorf("got error reading from client: %s", err)
	}

	var reply p.Message
	header := p.MessageHeader{
		0,
		17,
		m.Header().RequestID,
		p.OP_REPLY,
	}
	switch m.Header().OpCode {
	case p.OP_QUERY, p.OP_GET_MORE:
		log.Debug("------OP_QUERY, OP_GET_MORE : ", m.Header().OpCode)
		r, err := c.H.HandleQuery(m)
		if err != nil {
			log.Error("HandleQuery error ,", err)
		}
		r.SetHeader(header)
		reply = r
	case p.OP_COMMAND:
		log.Debug("------OP_COMMAND : ", m.Header().OpCode)
		r, err := c.H.HandleOpCommand(m)
		if err != nil {
			log.Error("HandleOpCommand error ,", err)
		}
		header.OpCode = p.OP_COMMAND_REPLY
		r.SetHeader(header)
		reply = r
	case p.OP_MSG:
		log.Debug("------OP_MSG : ", m.Header().OpCode)
		r, err := c.H.HandleOpMsg(m)
		if err != nil {
			log.Error("OP_MSG error ,", err)
		}
		header.OpCode = p.OP_MSG
		r.SetHeader(header)
		reply = r

	case p.OP_INSERT, p.OP_UPDATE, p.OP_DELETE:
		log.Debug("not supported OP_INSERT, OP_UPDATE, OP_DELETE")
		return p.NewStackErrorf("not supported now")
	//case p.OP_KILL_CURSORS:
	// todo
	default:
		log.Error(fmt.Sprintf("------not supported now, code: %v", m.Header().OpCode))
		return p.NewStackErrorf("not supported now")
	}
	err = p.SendMessage(reply, c.Conn)

	if err != nil {
		log.Error("send message error ", err)
	} else {
		log.Debug("send messaage success")
	}

	return err
}
