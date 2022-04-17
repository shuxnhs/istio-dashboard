# DataModel数据模型

为了进一步减少数据库操作的开发量，避免开发在model层重复编写代码实现基本的数据库操作，故结合gorm在model层提供了 \models\dataModel.go数据库操作模型。

下面将简单介绍下数据模型的设计与使用。



## 编写你的model类

假设你有一张ip详情表配置表，编写你的model层

```go
package models

// @ip详情表
type IpDetail struct {
	Id       int64 `gorm:"primary_key"`
	AppId    int64
	Type     int8
	Name     string
	IpDetail string
	Status   int8
	Ctime    int64
	Utime    int64
	ObjType  int8
}

// gorm的dbTabler的实现
func (i *IpDetail) TableName() string {
	return IPDETAIL
}

// @业务模型
type IpDetailModel struct {
	IpDetail
}

func NewIpDetailModel() *IpDetailModel {
	return &IpDetailModel{IpDetail{}}
}

func (i *IpDetailModel) GetTableStruct(isSlice bool) interface{} {
	if isSlice {
		return &[]IpDetail{}
	}
	return &IpDetail{}
}
```

在dataModel中定义了获取表struct的interface，每个model的业务模型都应该去实现这个方法GetTableStruct()，才可以使用dataModel的一系列操作。

```go
// 业务模型接口
type Model interface {
	GetTableStruct(isSlice bool) interface{}
}
```



## 数据库常用操作

使用dataModel，我们不用在models层再去重复写CURD操作的代码，便可以直接执行包括查询，增加，删除，更新，获取列表，聚合等操作。

#### 简单：操作实例

以下为对ipDetail的操作

  ```go
  // 插入新数据
  model := models.NewDataModel()
  res, err := model.Insert(&models.IpDetail{
      AppId:    2,
      Type:     5,
      Name:     "service-insert",
      IpDetail: "",
      Status:   0,
      Ctime:    0,
      Utime:    0,
      ObjType:  0,
  })
  if res >= 1 && err == nil {
    return "插入成功", err
  }
    return "插入失败", err


  // 查询：根据type获取数据
  model := models.NewDataModel()
  ipdetail, err := model.GetDataBy(models.NewIpDetailModel(), "type", 5, nil)
  return ipdetail.(*models.IpDetail), err


  // 更新：根据type更新数据
  model := models.NewDataModel()
  return model.UpdateAll(models.NewIpDetailModel(), 
        func(db *gorm.DB) *gorm.DB {
          return db.Where("type = ?", 5)}, 
        map[string]interface{}{"Name": "update-data"})


  // 删除数据
  model := models.NewDataModel()
  return model.DeleteAll(models.NewIpDetailModel(), 
        func(db *gorm.DB) *gorm.DB {
          return db.Where("type = ?", 5)})
  ```



#### 聚合操作

+ 获取最小值记录

   ```go
    func (d *DataModel) Min(bussinessModel Model, field string, 
         									 whereScopes Scopes) (interface{}, error)
    
    // 使用实例，获取type最小的一列
    d := models.NewDataModel()
    got, err := d.Min(models.NewIpDetailModel(), "type"
          func(db *gorm.DB) *gorm.DB {return db})
   ```

+ 获取最小值记录

   ```go
    func (d *DataModel) Max(bussinessModel Model, field string, 
         									 whereScopes Scopes) (interface{}, error)
    
    // 使用实例，获取type最大的一列
    d := models.NewDataModel()
    got, err := d.Max(models.NewIpDetailModel(), "type"
          func(db *gorm.DB) *gorm.DB {return db})
   ```



+ 获取总和

   ```go
   func (d *DataModel) Sum(bussinessModel Model, field string, whereScopes Scopes) int64
   
   // 使用实例，获取type的总和
    d := models.NewDataModel()
    got, err := d.Sum(models.NewIpDetailModel(), "type"
          func(db *gorm.DB) *gorm.DB {return db})
   ```



+ 获取数据总数

   ```go
   func (d *DataModel) GetList(bussinessModel Model, whereScopes Scopes, 
      													selectField []string) (interface{}, error)
   
   // 使用实例：获取ipdetail数据表条数
    d := models.NewDataModel()
    got, err := d.Sum(models.NewIpDetailModel(),
          func(db *gorm.DB) *gorm.DB {return db})
   ```



> gorm提供了Scopes：func(*DB) *DB
>
> 可以用于添加动态条件，dataModel使用了scopes可以用于传递我们services过程查询，删除，更新等操作的动态条件



#### 数据查询

查询结果需要转化为对应的struct或struct切片

+ 获取一条数据

  ```go
  func (d *DataModel) GetDataBy(bussinessModel Model, field string, 
        value interface{}, selectField []string) (interface{}, error) 
  
  // 操作实例，获取type=5的一条数据
  model := models.NewDataModel()
  ipdetail, err := model.GetDataBy(models.NewIpDetailModel(), "type", 5, nil)
  return ipdetail.(*models.IpDetail), err
  ```



+ 获取多条数据

  ```go
  func (d *DataModel) GetDataMoreBy(bussinessModel Model, field string, 
      value interface{}, limit int64, selectField []string, order string) 
  		(interface{}, error)
  
  // 操作实例，获取type=5的一条数据
  model := models.NewDataModel()
  ipdetail, err := model.GetDataMoreBy(models.NewIpDetailModel(), "type", 5, 
                                       2, nil, "type desc")
  return ipdetail.(*[]models.IpDetail), err
  ```



+ 根据条件获取一条数据

  ```go
  func (d *DataModel) GetData(bussinessModel Model, whereScopes Scopes, 
        											selectField []string) (interface{}, error)
  
  // 操作实例
  model := models.NewDataModel()
  got, err := d.GetData(models.NewIpDetailModel(),  
      func(db *gorm.DB) *gorm.DB {
  				return db.Where("app_id = ?", 202293).Order("type desc")
  		}, []string{"*"})
  // sql: select * from ip_detail where app_id = 202293 order by type desc;
  return got.(*models.IpDetail), err
  ```



+ 根据条件获取多条数据

  ```go
  func (d *DataModel) GetList(bussinessModel Model, whereScopes Scopes, 
                              selectField []string) (interface{}, error)
  
  // 操作实例
  model := models.NewDataModel()
  got, err := d.GetList(models.NewIpDetailModel(),  
      func(db *gorm.DB) *gorm.DB {
  				return db.Limit(2).Where("app_id = ?", 202293).Order("type desc").Offset(1)
  		}, []string{"*"})
  // sql: select * from ip_detail where app_id = 202293 order by type desc limit 1,2;
  return got.(*[]models.IpDetail), err
  ```



#### 数据更新

+ 根据条件更新

  ```go
  // 返回影响条数，错误
  func (d *DataModel) UpdateAll(bussinessModel Model, whereScopes Scopes, 
        												updateMap map[string]interface{}) (int64, error)
  
  // 操作实例
  d := models.NewDataModel()
  got, err := d.UpdateAll(models.NewIpDetailModel(), 
  	func(db *gorm.DB) *gorm.DB {
  			return db.Where("type = ?", 3)}, 
    map[string]interface{}{"type": 4,}
  )
  ```



#### 数据删除

+ 根据条件删除

  	```go
  // whereScope禁止传递nil或db，防止全表删除误操作
  func (d *DataModel) DeleteAll(bussinessModel Model, whereScopes Scopes) (int64, error)

  // 操作实例
  d := models.NewDataModel()
  affectRow, err := d.DeleteAll(models.NewIpDetailModel(),
  func(db *gorm.DB) *gorm.DB {return db.Where("type = ?", 3)})
  ```


+ 软删除：可以使用更新接口操作



#### 数据添加

+ 插入数据
  ```go
// 返回影响条数，错误
func (d *DataModel) Insert(data interface{}) (int64, error)

// 操作实例
d := models.NewDataModel()
data := &models.IpDetail{
        AppId:    2,
        Type:     5,
        Name:     "service-insert",
        IpDetail: "",
        Status:   0,
        Ctime:    0,
        Utime:    0,
        ObjType:  0,
}
affectRow, err := d.Insert(data)
```


+ todo：批量插入
