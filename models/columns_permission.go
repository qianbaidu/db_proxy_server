package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/prometheus/common/log"
)

type ColumnsPermission struct {
	Id          int64  `orm:"column(id);size(20)"`
	TablePermId int64  `orm:"column(table_perm_id);size(20)"`
	ColumnName  string `orm:"column(column_name);size(100)"`
	Status      int    `orm:"column(status);size(1)"`
	DbId        int64  `orm:"column(db_perm_id);size(20)"`
}

const tableName_ColumnsPermission = "columns_permission"

func (u *ColumnsPermission) TableName() string {
	return tableName_ColumnsPermission
}

func (u *ColumnsPermission) GetByTableIds(tableIds []int64) (cols []ColumnsPermission, err error) {
	cols = make([]ColumnsPermission, 0)
	_, err = orm.NewOrm().QueryTable(tableName_ColumnsPermission).
		Filter("table_perm_id", tableIds).
		Filter("status", 0).
		All(&cols)
	if err != nil {
		log.Error("GetTablesByDbId query error ", err)
		return cols, err
	}
	return cols, nil
}
