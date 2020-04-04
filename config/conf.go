package config

import (
	"github.com/go-ini/ini"
	"github.com/prometheus/common/log"
)

type Mysql struct {
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
	UserName string `ini:"user"`
	Password string `ini:"password"`
	DbName   string `ini:"db_name"`
}

type Server struct {
	Proto string `ini:"proto"`
	Addr  string `ini:"addr"`
	Debug bool   `ini:"debug"`
}

type Env struct {
	Debug bool `ini:"debug"`
}

type Config struct {
	Mysql           *Mysql
	MysqlServer     *Server
	MongdbServer    *Server
	Env             *Env
	SqlParserServer *Server
}

var App *Config

func parser(section string, value interface{}) {
	conf, err := ini.Load("config/conf.ini")
	if err != nil {
		log.Fatalf("load config error %s", err)
	}
	err = conf.Section(section).MapTo(value)
	if err != nil {
		log.Fatalf("parser config error %s", err)
	}
}

func NewConfig() {
	App = &Config{
		Mysql:           &Mysql{},
		MysqlServer:     &Server{},
		MongdbServer:    &Server{},
		Env:             &Env{},
		SqlParserServer: &Server{},
	}
	parser("mysql", App.Mysql)
	parser("mongodb-server", App.MongdbServer)
	parser("mysql-server", App.MysqlServer)
	parser("env", App.Env)
	parser("sql-parser-server", App.SqlParserServer)

	log.Info("init load app config success. ")
}

func init() {
	if App == nil {
		log.Info("int app config ")
		NewConfig()
	}
}
