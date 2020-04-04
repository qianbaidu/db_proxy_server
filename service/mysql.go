package service

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/common/log"
	"github.com/qianbaidu/db-proxy/config"
	"github.com/qianbaidu/db-proxy/models"
	"github.com/qianbaidu/db-proxy/util"
	mysqlUtil "github.com/qianbaidu/db-proxy/util/mysql"
)

func MysqlLoginAuth(userName string, salt, auth []byte) (int64, error) {
	u := &models.User{User: userName}
	if err := u.GetByUserName(); err != nil {
		return 0, err
	}
	if u.MysqlReadPermission == util.DB_PERMISSION_STATUS_NOT_ALLOW {
		return 0, mysqlUtil.ErrAccessDenied
	}
	if mysqlUtil.MysqlNativePasswordAuth(auth, salt, u.MysqlPassword) == false {
		return 0, mysqlUtil.ErrPasswordWrong
	}

	return u.Id, nil
}

func GetUser(userName string) (*models.User, error) {
	u := &models.User{User: userName}
	if err := u.GetByUserName(); err != nil {
		return u, err
	}
	if u.MongodbReadPermission == util.DB_PERMISSION_STATUS_NOT_ALLOW {
		return u, mysqlUtil.ErrAccessDenied
	}

	return u, nil
}

type UserDbList struct {
	DBList map[string]UserDbPerms
	Super  bool
}

type UserDbPerms struct {
	models.DbList
	models.DbPermission
	Tables map[string]UserDbTable
}

type UserDbTable struct {
	models.TablePermission
	Columns []models.ColumnsPermission
}

func GetUserDbList(userId int64, dbType int) (userDbList UserDbList, err error) {
	log.Info("GetUserDbList start")
	userDbList = UserDbList{}
	var (
		userPerDbList []models.UserDbPerList
		tableM        models.TablePermission
		dbPermM       models.DbPermission
		colM          models.ColumnsPermission
		isSuper       = true
	)
	//query user permission db list
	userPerDbList, err = dbPermM.GetUserDbList(userId, dbType)
	if err != nil {
		log.Error("GetUserDbList query GetUserDbList error ", err)
		return userDbList, err
	}
	log.Info("userPerDbList:", userPerDbList)

	//tablePermList, err := GetDbTablesPerm(userPerDbList)
	permIds := make([]int64, 0)
	for _, v := range userPerDbList {
		if v.HaveAll == util.DB_PERMISSION_STATUS_ALLOW && v.HaveSecretColumns == util.DB_PERMISSION_STATUS_ALLOW {
			continue
		}
		permIds = append(permIds, v.PermId)
	}

	tableIds := make([]int64, 0)
	permTables := make([]models.TablePermission, 0)
	if len(permIds) > 1 {
		// querey db tables by db_perm_id
		permTables, err = tableM.GetTablesByDbPermId(permIds)
		if err != nil {
			log.Error("query GetTablesByDbPermId error ", err)
			return userDbList, err
		}
		// reset table list to map
		//tableMapList := make(map[int64]models.TablePermission, 0)
		for _, v := range permTables {
			tableIds = append(tableIds, v.Id)
		}
	}

	log.Info("tableIds:", tableIds)
	log.Info("PermId", permIds)

	colMapList := make(map[int64][]models.ColumnsPermission, 0)

	if len(tableIds) > 0 {
		// querey table columns by table_perm_id
		cols, err := colM.GetByTableIds(tableIds)
		if err != nil {
			log.Error("query GetByTableIds error ", err)
			return userDbList, err
		}
		// reset columns permission array to map
		for _, v := range cols {
			if colsList, ok := colMapList[v.TablePermId]; ok {
				colsList = append(colsList, v)
				colMapList[v.TablePermId] = colsList
			} else {
				colsList := make([]models.ColumnsPermission, 0)
				colsList = append(colsList, v)
				colMapList[v.TablePermId] = colsList
			}
		}
		log.Info("colMapList", colMapList)
	}
	// reset table data
	userPermtables := make(map[int64]UserDbTable, 0)
	for _, v := range permTables {
		tableMap := UserDbTable{TablePermission: v}
		if cols, ok := colMapList[v.Id]; ok {
			tableMap.Columns = cols
		}
		userPermtables[v.DbPermId] = tableMap
	}
	log.Info("userPermtables:", userPermtables)

	db := make(map[string]UserDbPerms, 0)
	for _, v := range userPerDbList {
		if v.HaveAll == util.DB_PERMISSION_STATUS_NOT_ALLOW || v.HaveSecretColumns == util.DB_PERMISSION_STATUS_NOT_ALLOW {
			isSuper = false
		}
		usertable := UserDbPerms{DbList: v.DbList}
		usertable.HaveAll = v.DbPermission.HaveAll
		if table, ok := userPermtables[v.PermId]; ok {
			t := make(map[string]UserDbTable, 0)
			t[table.Table] = table
			usertable.Tables = t
		}

		db[v.DbName] = usertable
	}
	userDbList.DBList = db
	userDbList.Super = isSuper

	log.Info("userDbList:", userDbList)
	j, e := json.Marshal(userDbList)
	log.Infof("userDbList: %s \nerr : %v", string(j), e)

	return userDbList, nil
}

func ParserSql(sql string, user string) (res SqlParserResp, err error) {
	url := fmt.Sprintf(util.API_MYSQL_PARSER, config.App.SqlParserServer.Addr)
	res, err = parserSql(url, sql, user)
	if err != nil {
		log.Error("parser sql err ", err)
	}
	return res, err
}

type SqlParser struct {
	Table  []string `json:"table"`
	Column []struct {
		IsSelect   string `json:"is_select"`
		ColumnName string `json:"column_name"`
		TableName  string `json:"table_name"`
	} `json:"column"`
	Alias map[string]string `json:"alias"`
}

type SqlParserResp struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Result  SqlParser `json:"result"`
}

func parserSql(url string, sql string, user string) (res SqlParserResp, err error) {
	type Params struct{ Sql string `json:"sql"` }
	postData := Params{sql}
	resp, err := util.HttpPost(url, postData)
	if err != nil {
		log.Error("ParserSql error, response : ", resp, err)
		return res, err
	}

	if err = util.CheckRestFulResponseStatus(string(resp)); err != nil {
		return res, err
	}

	parser := SqlParserResp{}
	if err := json.Unmarshal(resp, &parser); err != nil {
		return res, err
	}

	return parser, nil
}
