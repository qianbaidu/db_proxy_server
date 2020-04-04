package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/prometheus/common/log"
)

type TablePermission struct {
	Id         int64  `orm:"column(id);size(20)"`
	Table      string `orm:"column(table_name);size(100)"`
	Status     int    `orm:"column(status);size(1)"`
	DbPermId   int64  `orm:"column(db_perm_id);size(20)"`
	HaveSecret int    `orm:"column(have_secret);size(1)"`
}

const tableName_TablePermission = "table_permission"

func (u *TablePermission) TableName() string {
	return tableName_TablePermission
}

func (u *TablePermission) GetTablesByDbPermId(permIds []int64) (tables []TablePermission, err error) {
	log.Info("start GetTablesByDbPermId")
	tables = make([]TablePermission, 0)
	q := orm.NewOrm().QueryTable(tableName_TablePermission)
	_, err = q.Filter("db_perm_id__in", permIds).Filter("status", 0).All(&tables)
	log.Info("tables:", tables, " err :", err)
	if err != nil {
		log.Error("GetTablesByDbId query error ", err)
		return tables, err
	}
	return tables, nil
}
