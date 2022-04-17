package model

import (
	"errors"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

/**
 * @Title		  基于数据库的数据基类
 * @Description	  提供更常用的数据库操作接口，避免Model子类的重复开发工作
 * @Author		  hexiaohong <hexiaohong@kingdee.com> 2020-08-17
 */
type DataModel struct {
}

func NewDataModel() *DataModel {
	return &DataModel{}
}

// 业务模型接口
type Model interface {
	GetTableStruct(isSlice bool) interface{}
	TableName() string
}

type Scopes func(db *gorm.DB) *gorm.DB

// V2版本的可以不用去实现GetTableStruct方法
// V2版本新增的方法

// 增加相关操作

/**
 * @Description  			      		批量插入数据
 * @param 	interface	data 			对应的结构体切片
 * @return  int64		rowsAffected	影响的条数
 */
func (d *DataModel) InsertMore(data interface{}) (int64, error) {
	appDB := DB.Create(data)
	return appDB.RowsAffected, appDB.Error
}

/**
 * @Description  			      		批量插入数据
 * @param 	interface	data 			对应的结构体切片
 * @return  int64		rowsAffected	影响的条数
 */
func (d *DataModel) TxInsertMore(tx *gorm.DB, data interface{}) (int64, error) {
	appDB := tx.Create(data)
	return appDB.RowsAffected, appDB.Error
}

// 查询相关操作

/**
 * @Description  			      获取一条纪录
 * @param  string     field       需要查询的字段名，通常为主键或带有唯一索引的字段
 * @param  interface  value       查询的值
 * @param  []string   selectField 需要select的字段
 * @return interface  dataStruct  返回对应的map
 */
func (d *DataModel) GetDataByToMap(bussinessModel Model, field string, value interface{}, selectField []string) (map[string]interface{}, error) {
	results, err := d.GetDataMoreByToMap(bussinessModel, field, value, 1, selectField, "")
	if len(results) == 0 {
		return nil, err
	}
	return results[0], err
}

/**
 * @Description  			      获取多条纪录
 * @param  string     field       需要查询的字段名，通常为主键或带有唯一索引的字段
 * @param  interface  value       查询的值
 * @param  int        limit       需要获取的条数,0表示无限制获取全部
 * @param  []string   selectField 需要select的字段
 * @param  string     order		  需要排序的字段，如："age desc, name",  "age desc"
 * @return interface  dataStruct  返回对应的map
 */
func (d *DataModel) GetDataMoreByToMap(bussinessModel Model, field string, value interface{}, limit int, selectField []string, order string) ([]map[string]interface{}, error) {
	tableName := bussinessModel.TableName()
	var results []map[string]interface{}
	appDB := DB
	if order != "" {
		appDB = appDB.Order(order)
	}
	if limit != 0 {
		appDB = appDB.Limit(limit)
	}
	appDB = appDB.Where(field+" = ?", value)
	if !reflect.DeepEqual(selectField, []string{"*"}) && selectField != nil {
		appDB = appDB.Select(selectField)
	}
	err := appDB.Table(tableName).Find(&results).Error
	return results, err
}

/**
 * @Description  			      获取一条纪录
 * @param  string     field       需要查询的字段名，通常为主键或带有唯一索引的字段
 * @param  interface  value       查询的值
 * @param  []string   selectField 需要select的字段
 * @return interface  dataStruct  返回对应的结构体
 */
func (d *DataModel) GetDataBy(bussinessModel Model, field string, value interface{}, selectField []string) (interface{}, error) {
	db := DB
	tableStruct := bussinessModel.GetTableStruct(false)
	db = db.Limit(1).Where(field+" = ?", value)
	if !reflect.DeepEqual(selectField, []string{"*"}) && selectField != nil {
		db = db.Select(selectField)
	}
	err := db.Find(tableStruct).Error
	return tableStruct, err
}

/**
 * @Description  			      获取多条纪录
 * @param  string     field       需要查询的字段名，通常为主键或带有唯一索引的字段
 * @param  interface  value       查询的值
 * @param  int        limit       需要获取的条数,0表示无限制获取全部
 * @param  []string   selectField 需要select的字段
 * @param  string     order		  需要排序的字段，如："age desc, name",  "age desc"
 * @return interface  dataStruct  返回对应的结构体
 */
func (d *DataModel) GetDataMoreBy(bussinessModel Model, field string, value interface{}, limit int, selectField []string, order string) (interface{}, error) {
	tableStruct := bussinessModel.GetTableStruct(true)
	appDB := DB
	if order != "" {
		appDB = appDB.Order(order)
	}
	if limit != 0 {
		appDB = appDB.Limit(limit)
	}
	appDB = appDB.Where(field+" = ?", value)
	if !reflect.DeepEqual(selectField, []string{"*"}) && selectField != nil {
		appDB = appDB.Select(selectField)
	}
	err := appDB.Find(tableStruct).Error
	return tableStruct, err
}

/**
 * @Description  			      						动态条件获取一条纪录
 * @param  func (db *gorm.DB) *gorm.DB	whereScopes		动态设置查询条件
 * @param  []string   					selectField 	需要select的字段
 * @return interface  					dataStruct  	返回对应的结构体
 */
func (d *DataModel) GetData(bussinessModel Model, whereScopes Scopes, selectField []string) (interface{}, error) {
	tableStruct := bussinessModel.GetTableStruct(false)
	limitScopes := func(db *gorm.DB) *gorm.DB {
		return db.Limit(1)
	}
	selectScopes := func(db *gorm.DB) *gorm.DB {
		return db
	}
	if !reflect.DeepEqual(selectField, []string{"*"}) && selectField != nil {
		selectScopes = func(db *gorm.DB) *gorm.DB {
			return db.Select(selectField)
		}
	}
	err := DB.Scopes(limitScopes, whereScopes, selectScopes).Find(tableStruct).Error
	return tableStruct, err
}

/**
 * @Description  			      						动态条件获取多条纪录
 * @param  func (db *gorm.DB) *gorm.DB 	whereScopes		动态设置查询条件
 * @param  []string   					selectField 	需要select的字段
 * @return interface  					dataStruct  	返回对应的结构体
 */
func (d *DataModel) GetList(bussinessModel Model, whereScopes Scopes, selectField []string) (interface{}, error) {
	tableStruct := bussinessModel.GetTableStruct(true)
	selectScopes := func(db *gorm.DB) *gorm.DB {
		return db
	}
	if !reflect.DeepEqual(selectField, []string{"*"}) && selectField != nil {
		selectScopes = func(db *gorm.DB) *gorm.DB {
			return db.Select(selectField)
		}
	}
	err := DB.Scopes(whereScopes, selectScopes).Find(tableStruct).Error
	return tableStruct, err
}

// 聚合操作

/**
 * @Description  			      						动态条件获取纪录条数
 * @param  func (db *gorm.DB) *gorm.DB 	whereScopes		动态设置查询条件
 * @return int64  						count  			返回对应的条数
 */
func (d *DataModel) Count(bussinessModel Model, whereScopes Scopes) (int64, error) {
	var count int64 = 0
	tableStruct := bussinessModel.GetTableStruct(true)
	err := DB.Scopes(whereScopes).Find(tableStruct).Count(&count).Error
	return count, err
}

/**
 * @Description  			      						取某字段最大的那条记录
 * @param  string     					field       	需要查询的字段名
 * @param  func (db *gorm.DB) *gorm.DB 	whereScopes		动态设置查询条件
 * @return interface  					dataStruct  	返回对应的结构体
 */
func (d *DataModel) Max(bussinessModel Model, field string, whereScopes Scopes) (interface{}, error) {
	tableStruct := bussinessModel.GetTableStruct(false)
	maxScopes := func(db *gorm.DB) *gorm.DB {
		return db.Order(field + " desc")
	}
	err := DB.Scopes(whereScopes, maxScopes).First(tableStruct).Error
	return tableStruct, err
}

/**
 * @Description  			      						取某字段最小的那条记录
 * @param  string     					field       	需要查询的字段名
 * @param  func (db *gorm.DB) *gorm.DB 	whereScopes		动态设置查询条件
 * @return interface  					dataStruct  	返回对应的结构体
 */
func (d *DataModel) Min(bussinessModel Model, field string, whereScopes Scopes) (interface{}, error) {
	tableStruct := bussinessModel.GetTableStruct(false)
	minScopes := func(db *gorm.DB) *gorm.DB {
		return db.Order(field)
	}
	err := DB.Scopes(whereScopes, minScopes).First(tableStruct).Error
	return tableStruct, err
}

/**
 * @Description  			      						取某字段的总和
 * @param  string     					field       	需要查询的字段名
 * @param  func (db *gorm.DB) *gorm.DB 	whereScopes		动态设置查询条件
 * @return int64  						sum  			返回对应字段的总和
 */
func (d *DataModel) Sum(bussinessModel Model, field string, whereScopes Scopes) int64 {
	tableStruct := bussinessModel.GetTableStruct(false)
	var total int64 = 0
	selectScopes := func(db *gorm.DB) *gorm.DB {
		return db.Select("sum(" + field + ") as total")
	}
	row := DB.Model(tableStruct).Scopes(selectScopes, whereScopes).Row()
	_ = row.Scan(&total)
	return total
}

// 删除相关操作

/**
 * @Description  			      						动态条件删除纪录
 * @param  func (db *gorm.DB) *gorm.DB 	whereScopes		动态设置查询条件
 * @return  int64		rowsAffected					影响的条数
 */
func (d *DataModel) DeleteAll(bussinessModel Model, whereScopes Scopes) (int64, error) {
	// important: 禁止全表删除,防止误操作
	if whereScopes == nil || reflect.DeepEqual(whereScopes, func(db *gorm.DB) *gorm.DB { return db }) {
		return 0, errors.New("Delete All Is Forbidden ")
	}
	tableStruct := bussinessModel.GetTableStruct(false)
	appDB := DB.Scopes(whereScopes).Delete(tableStruct)
	return appDB.RowsAffected, appDB.Error
}

// 增加相关操作

/**
 * @Description  			      		插入数据
 * @param 	interface	data 			对应的结构体
 * @return  int64		rowsAffected	影响的条数
 */
func (d *DataModel) Insert(data interface{}) (int64, error) {
	appDB := DB.Create(data)
	return appDB.RowsAffected, appDB.Error
}

/**
 * @Description  			      		插入数据
 * @param 	interface	data 			对应的结构体
 * @return  int64		rowsAffected	影响的条数
 */
func (d *DataModel) TxInsert(tx *gorm.DB, data interface{}) (int64, error) {
	appDB := tx.Create(data)
	return appDB.RowsAffected, appDB.Error
}

// 更新相关操作

/**
 * @Description  			      						动态条件更新纪录
 * @param  func (db *gorm.DB) *gorm.DB 	whereScopes		动态设置查询条件
 * @return int64  						sum  			返回对应字段的总和
 */
func (d *DataModel) UpdateAll(bussinessModel Model, whereScopes Scopes, updateMap map[string]interface{}) (int64, error) {
	tableStruct := bussinessModel.GetTableStruct(false)
	appDB := DB.Model(tableStruct).Scopes(whereScopes).Updates(updateMap)
	return appDB.RowsAffected, appDB.Error
}

/**
 * @Description  			      						动态条件事务更新纪录
 * @param  func (db *gorm.DB) *gorm.DB 	whereScopes		动态设置查询条件
 * @return int64  						sum  			返回对应字段的总和
 */
func (d *DataModel) TxUpdateAll(tx *gorm.DB, bussinessModel Model, whereScopes Scopes, updateMap map[string]interface{}) (int64, error) {
	tableStruct := bussinessModel.GetTableStruct(false)
	appDB := tx.Model(tableStruct).Scopes(whereScopes).Updates(updateMap)
	return appDB.RowsAffected, appDB.Error
}

// 直接执行sql语句

var ExpectMethod = map[string]bool{"INSERT": true, "UPDATE": true, "DELETE": true}

/**
 * @Description
 * @param	string			sql		直接执行sql语句
 * @return  interface{}		返回对应的结果(ToArray->ToMap)
 */
func (d *DataModel) ExecuteSql(sql string) (interface{}, error) {
	a := make([]interface{}, 0)
	cols := make([]string, 0)

	if ExpectMethod[strings.ToUpper(strings.Split(sql, " ")[0])] {
		if err := DB.Exec(sql).Error; err != nil {
			return nil, err
		}
	} else {
		raw := DB.Raw(sql)
		if err := raw.Error; err != nil {
			return nil, err
		}

		rows, err := raw.Rows()

		if err != nil {
			return nil, err
		}

		cols, _ = rows.Columns()
		for rows.Next() {
			columns := make([]interface{}, len(cols))
			columnPointers := make([]interface{}, len(cols))
			for i, _ := range columns {
				columnPointers[i] = &columns[i]
			}

			if err := rows.Scan(columnPointers...); err != nil {
				return nil, err
			}
			m := make(map[string]interface{})

			//for i, colName := range cols {
			for i := 0; i < len(cols); i++ {
				val := *columnPointers[i].(*interface{})
				formatValue := ""
				switch val.(type) {
				case []uint8:
					formatValue = strings.Join(strings.Fields(string([]byte(val.([]uint8)[:]))), "")
				}
				m[cols[i]] = formatValue

			}

			a = append(a, m)
		}
	}

	return a, nil
}

// 事务相关
/**
 * @Description  			   开启事务
 * @return  *gorm.DB		appDB	*gorm.DB
 */
func (d *DataModel) Begin() (appDB *gorm.DB, err error) {
	appDB = DB.Begin()
	err = appDB.Error
	return
}

/**
 * @Description  			   回滚事务
 * @return  *gorm.DB		appDB	*gorm.DB
 */
func (d *DataModel) Rollback() (appDB *gorm.DB, err error) {
	appDB = DB.Rollback()
	err = appDB.Error
	return
}

/**
 * @Description  			   提交事务
 * @return  *gorm.DB		appDB	*gorm.DB
 */
func (d *DataModel) Commit() (appDB *gorm.DB, err error) {
	appDB = DB.Commit()
	err = appDB.Error
	return
}
