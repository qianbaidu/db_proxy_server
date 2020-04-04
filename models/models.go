package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/common/log"
	"github.com/qianbaidu/db-proxy/config"
)

func RegisterDb() {
	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		log.Error("orm.RegisterDriver eror : ", err)
	}
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		config.App.Mysql.UserName,
		config.App.Mysql.Password,
		config.App.Mysql.Host,
		config.App.Mysql.Port,
		config.App.Mysql.DbName,
	)

	log.Info("init db dns : ", dns)
	err = orm.RegisterDataBase("default", "mysql", dns)
	if err != nil {
		log.Error("orm.RegisterDataBase eror : ", err)
	}
	orm.RegisterModel(new(User), new(DbList), new(DbPermission), new(TablePermission), new(ColumnsPermission), new(EventLog))
}

