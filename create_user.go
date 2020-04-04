package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/qianbaidu/db-proxy/models"
	"strings"
)

func mysqlPassword(password string) string {
	//step 1 : stage1Hash = SHA1(password)
	hash := sha1.New()
	hash.Write([]byte(password))
	s1 := hash.Sum(nil)

	hash.Reset()
	hash.Write(s1)
	s2 := hash.Sum(nil)

	s := strings.ToUpper(hex.EncodeToString(s2))

	return fmt.Sprintf("*%s", s)
}

func mongoDbPassword(user string, pass string) string {
	credsum := md5.New()
	credsum.Write([]byte(user + ":mongo:" + pass))

	return hex.EncodeToString(credsum.Sum(nil))
}

func main() {
	user := flag.String("u", "", "user name")
	password := flag.String("p", "", "user password")

	flag.Parse()
	if len(*user) < 1 || len(*password) < 1 {
		flag.PrintDefaults()
		panic("User name and password cannot be empty")
	}
	models.RegisterDb()
	u := models.User{
		User:                       *user,
		MysqlPassword:              mysqlPassword(*password),
		MysqlReadPermission:        1,
		MysqlReadWritePermission:   0,
		MongodbPassword:            mongoDbPassword(*user, *password),
		MongodbReadPermission:      1,
		MongodbReadWritePermission: 0,
	}
	if id, err := u.Insert(); err != nil {
		panic(fmt.Sprintf("insert user error : %s", err))
	} else {
		fmt.Println("create user success, id :", id)
	}
}
