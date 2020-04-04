package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/prometheus/common/log"
	mysqlUtil "github.com/qianbaidu/db-proxy/util/mysql"
)

type User struct {
	Id                         int64  `orm:"column(id);size(20)"`
	User                       string `orm:"column(user);size(16)"`
	MysqlPassword              string `orm:"column(mysql_password);size(41)"`
	MysqlReadPermission        int    `orm:"column(mysql_read_permission);size(1)"`
	MysqlReadWritePermission   int    `orm:"column(mysql_read_write_permission);size(1)"`
	MongodbPassword            string `orm:"column(mongodb_password);size(41)"`
	MongodbReadPermission      int    `orm:"column(mongodb_read_permission);size(1)"`
	MongodbReadWritePermission int    `orm:"column(mongodb_read_write_permission);size(1)"`
}

const tablename_User = "user"

func (u *User) TableName() string {
	return tablename_User
}

func (u *User) Read(fields ...string) error {
	o := orm.NewOrm()
	if err := o.Read(u, fields...); err != nil {
		log.Error("read rule by id error : ", err)
		return err
	}
	return nil
}

func (u *User) GetByUserName() error {
	q := orm.NewOrm().QueryTable(tablename_User)
	err := q.Filter("user", u.User).One(u)
	if err != nil {
		log.Error("GetByUser query error : ", err)
		return mysqlUtil.ErrUserNotExists
	}
	return nil
}

func (u *User) Insert() (int64, error) {
	if id, err := orm.NewOrm().Insert(u); err != nil {
		log.Error("insert user error :", err)
		return id, err
	} else {
		return id, nil
	}
}
