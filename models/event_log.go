package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/prometheus/common/log"
)

type EventLog struct {
	Id             int64  `orm:"column(id);size(20)"`
	Sql            string `orm:"column(sql);size(16)"`
	Status         int    `orm:"column(status);size(2)"`
	CreateTime     string `orm:"column(create_time);size(6)"`
	DbId           int64  `orm:"column(db_id);size(20)"`
	UserId         int64  `orm:"column(user_id);size(20)"`
	UpdateDatatime string `orm:"column(update_datatime);size(6)"`
}

const tablename_EventLog = "event_log"

func (e *EventLog) TableName() string {
	return tablename_EventLog
}

func (e *EventLog) InsertData() (int64, error) {
	if id, err := orm.NewOrm().Insert(e); err != nil {
		log.Error("insert query log error :", err)
		return id, err
	} else {
		return id, nil
	}
}
