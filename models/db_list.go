package models

type DbList struct {
	Id       int64  `orm:"column(id);size(20)"`
	Name     string `orm:"column(name);size(20)"`
	Ip       string `orm:"column(ip);size(64)"`
	Port     int    `orm:"column(port);size(5)"`
	DbName   string `orm:"column(db_name);size(32)"`
	Username string `orm:"column(username);size(32)"`
	Password string `orm:"column(password);size(128)"`
	Encode   int    `orm:"column(encode);size(1)"`
	Status   int    `orm:"column(status);size(1)"`
	DbType   int    `orm:"column(db_type);size(1)"`
}

const tableName_DbList = "db_list"

func (u *DbList) TableName() string {
	return tableName_DbList
}


