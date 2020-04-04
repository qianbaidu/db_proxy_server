package provider

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"github.com/prometheus/common/log"
	"time"
)

type Nonce struct {
	Nonce string `bson:"nonce"`
}

type SaslCmd struct {
	Start          int    `bson:"saslStart,omitempty"`
	Continue       int    `bson:"saslContinue,omitempty"`
	ConversationId int    `bson:"conversationId,omitempty"`
	Mechanism      string `bson:"mechanism,omitempty"`
	Payload        string `bson:"payload"`
}

type SaslResult struct {
	Ok    bool `bson:"ok"`
	NotOk bool `bson:"code"` // Server <= 2.3.2 returns ok=1 & code>0 on errors (WTF?)
	Done  bool

	ConversationId int `bson:"conversationId"`
	Payload        []byte
	ErrMsg         string `bson:"errmsg"`
}

type UserLoginReuqest struct {
	Id           int32
	CreateTime   int64
	UpdateTime   int64
	User         string
	Method       string
	ClientCode   string
	RandCode     string
	ValidateCode string
	I            int
	C            string
	Salt         string
}

// todo 给成每次自动生成
var SaltB = []byte{195, 48, 163, 156, 57, 205, 41, 124, 4, 6, 112, 213, 191, 113, 145, 139}

const (
	SaltS = "wzCjnDnNKXwEBnDVv3GRiw=="
	SaltI = 10000
)

var UserLoginData map[string]UserLoginReuqest

func init() {
	if UserLoginData == nil {
		UserLoginData = make(map[string]UserLoginReuqest, 0)
		//go CheckUserLoginData()
	}
}

func CheckUserLoginData() {
	log.Info("start CheckUserLoginData")
	for {
		for k, v := range UserLoginData {
			current := time.Now().Unix()
			if current-v.UpdateTime > 3*60 {
				delete(UserLoginData, k)
			}
		}
		time.Sleep(time.Second * 60)
	}
}

var b64 = base64.StdEncoding

func PassHash(user string, pass string) string {
	credsum := md5.New()
	credsum.Write([]byte(user + ":mongo:" + pass))

	return hex.EncodeToString(credsum.Sum(nil))
}

func SaltPassword(pass string, salt []byte, iterCount int) []byte {
	h := sha1.New
	mac := hmac.New(h, []byte(pass))
	mac.Write(salt)
	mac.Write([]byte{0, 0, 0, 1})
	ui := mac.Sum(nil)
	hi := make([]byte, len(ui))
	copy(hi, ui)
	for i := 1; i < iterCount; i++ {
		mac.Reset()
		mac.Write(ui)
		mac.Sum(ui[:0])
		for j, b := range ui {
			hi[j] ^= b
		}
	}

	return hi
}

func ServerSignature(saltedPass []byte, authMsg []byte) []byte {
	newHash := sha1.New
	mac := hmac.New(newHash, saltedPass)
	mac.Write([]byte("Server Key"))
	serverKey := mac.Sum(nil)

	mac = hmac.New(newHash, serverKey)
	mac.Write(authMsg)
	serverSignature := mac.Sum(nil)

	encoded := make([]byte, b64.EncodedLen(len(serverSignature)))
	b64.Encode(encoded, serverSignature)
	return encoded
}
