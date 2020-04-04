package mysql

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func CalcPassword(scramble, password []byte) []byte {
	if len(password) == 0 {
		return nil
	}

	// stage1Hash = SHA1(password)
	crypt := sha1.New()
	crypt.Write(password)
	stage1 := crypt.Sum(nil)

	// scrambleHash = SHA1(scramble + SHA1(stage1Hash))
	// inner Hash
	crypt.Reset()
	crypt.Write(stage1)
	hash := crypt.Sum(nil)

	// outer Hash
	crypt.Reset()
	crypt.Write(scramble)
	crypt.Write(hash)
	scramble = crypt.Sum(nil)

	// token = scrambleHash XOR stage1Hash
	for i := range scramble {
		scramble[i] ^= stage1[i]
	}
	return scramble
}

func MysqlNativePasswordAuth(reply, salt []byte, mysqlNativePassword string) bool {
	/*
		SERVER:  recv(reply)
				 hash_stage1=xor(reply, sha1(salt,hash))
				 candidate_hash2=sha1(hash_stage1)
				 check(candidate_hash2==hash)
	*/
	if len(reply) == 0 {
		return false
	}

	if mysqlNativePassword == "" {
		return false
	}

	if strings.Contains(mysqlNativePassword, "*") {
		mysqlNativePassword = mysqlNativePassword[1:]
	}

	hash, err := hex.DecodeString(mysqlNativePassword)
	if err != nil {
		return false
	}

	// scramble = SHA1(salt+hash)
	crypt := sha1.New()
	crypt.Write(salt)
	crypt.Write(hash)
	scramble := crypt.Sum(nil)

	// token = scramble XOR stage1Hash
	for i := range scramble {
		scramble[i] ^= reply[i]
	}
	hashStage1 := scramble

	crypt.Reset()
	crypt.Write(hashStage1)
	candidateHash2 := crypt.Sum(nil)

	return bytes.Equal(candidateHash2, hash)
}

// NativePassword generates a mysql_native_password string.
func NativePassword(password string) string {
	if len(password) == 0 {
		return ""
	}

	// native = sha1(sha1(password))

	hash := sha1.New()
	hash.Write([]byte(password))
	s1 := hash.Sum(nil)

	hash.Reset()
	hash.Write(s1)
	s2 := hash.Sum(nil)

	s := strings.ToUpper(hex.EncodeToString(s2))

	return fmt.Sprintf("*%s", s)
}

func RandomBuf(size int) []byte {
	buf := make([]byte, size)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < size; i++ {
		buf[i] = byte(rand.Intn(127))
		if buf[i] == 0 || buf[i] == byte('$') {
			buf[i]++
		}
	}
	return buf
}

