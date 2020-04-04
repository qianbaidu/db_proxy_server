package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/prometheus/common/log"
)

type DbPermission struct {
	Id                int64 `orm:"column(id);size(20)"`
	HaveAll           int   `orm:"column(have_all);size(1)"`
	HaveSecretColumns int   `orm:"column(have_secret_columns);size(1)"`
	Status            int   `orm:"column(status);size(1)"`
	DbId              int64 `orm:"column(db_id);size(20)"`
	UserId            int64 `orm:"column(user_id);size(20)"`
}

const tableName_DbPermission = "db_permission"

func (u *DbPermission) TableName() string {
	return tableName_DbPermission
}

type UserDbPerList struct {
	DbList
	DbPermission
	PermId int64 `orm:"column(perm_id);size(20)"`
}

func (u *DbPermission) GetUserDbList(userId int64, dbType int) (dbList []UserDbPerList, err error) {
	sql := `select A.*,B.*,A.id as perm_id from(select * from db_permission where user_id = ? ) as A 
			left join db_list  as B ON A.db_id = B.id where B.status = 0  and db_type = ?;`
	dbList = make([]UserDbPerList, 0)
	if _, err := orm.NewOrm().Raw(sql, userId, dbType).QueryRows(&dbList); err != nil {
		log.Error("GetUserDbList query error ", err)
		return dbList, err
	}
	return dbList, nil
}
