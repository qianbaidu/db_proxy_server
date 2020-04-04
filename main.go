package main

import (
	"github.com/prometheus/common/log"
	"github.com/qianbaidu/db-proxy/config"
	"github.com/qianbaidu/db-proxy/models"
	"github.com/qianbaidu/db-proxy/proxy"
	"net"
	"sync"
)

func MongoDbServer(ctx *sync.WaitGroup) {
	l, err := net.Listen(config.App.MongdbServer.Proto, config.App.MongdbServer.Addr)
	if err != nil {
		log.Fatal("start server error : ", err)
	}
	log.Info("start mongodb proxy server ", config.App.MongdbServer.Addr)

	for {
		c, err := l.Accept()
		if err != nil {
			log.Error(err)
			break
		}
		go func() {
			defer func() {
				if err := recover(); err != nil {
					c.Close()
				}
			}()
			proxy.RunMongo(c)
		}()
	}
	log.Info("")
	ctx.Done()
}

func MysqlServer(ctx *sync.WaitGroup) {
	l, err := net.Listen(config.App.MysqlServer.Proto, config.App.MysqlServer.Addr)
	if err != nil {
		log.Fatal("start server error : ", err)
	}
	log.Info("start mysql proxy server ", config.App.MysqlServer.Addr)

	for {
		c, err := l.Accept()
		if err != nil {
			log.Error(err)
			break
		}
		go func() {
			defer func() {
				if err := recover(); err != nil {
					c.Close()
				}
			}()
			proxy.RunMysql(c)
		}()
	}
	ctx.Done()
}

func init() {
	models.RegisterDb()
}

func main() {
	ctx := &sync.WaitGroup{}
	ctx.Add(1)
	go MongoDbServer(ctx)

	ctx.Add(1)
	go MysqlServer(ctx)
	ctx.Wait()

}
