package genfunc

const (
	genTnf = `
// TableName get sql table name.获取数据库表名
func (m *{{.StructName}}) TableName() string {
	return "{{.TableName}}"
}
`
	genBase = `
package {{.PackageName}}
import (
	"context"

	"github.com/jinzhu/gorm"
)

var globalIsRelated bool // 全局预加载

// prepare for other
type _BaseMgr struct {
	*gorm.DB
	ctx       *context.Context
	isRelated bool
}

// SetCtx set context
func (obj *_BaseMgr) SetCtx(c *context.Context) {
	obj.ctx = c
}

// GetDB get gorm.DB info
func (obj *_BaseMgr) GetDB() *gorm.DB {
	return obj.DB
}

// UpdateDB update gorm.DB info
func (obj *_BaseMgr) UpdateDB(db *gorm.DB) {
	obj.DB = db
}

// GetIsRelated Query foreign key Association.获取是否查询外键关联(gorm.Related)
func (obj *_BaseMgr) GetIsRelated() bool {
	return obj.isRelated
}

// SetIsRelated Query foreign key Association.设置是否查询外键关联(gorm.Related)
func (obj *_BaseMgr) SetIsRelated(b bool) {
	obj.isRelated = b
}

type SqlContext struct {
	Query   map[string]interface{}
	Not     map[string]interface{}
	Gt      map[string]interface{}
	Lt      map[string]interface{}
	In      map[string]interface{}
	Changes map[string]interface{}
	Order   []string
}

func NewSqlContext() *SqlContext {
	return &SqlContext{
		Query:   make(map[string]interface{}),
		Changes: make(map[string]interface{}),
		Not:     make(map[string]interface{}),
		Gt:      make(map[string]interface{}),
		Lt:      make(map[string]interface{}),
		In:      make(map[string]interface{}),
		Order:   []string{},
	}
}

// Option overrides behavior of Connect.
type Option interface {
	Apply(*SqlContext) *SqlContext
}

type OptionFunc func(*SqlContext) *SqlContext

func (f OptionFunc) Apply(o *SqlContext) *SqlContext {
	f(o)
	return o
}

// OpenRelated 打开全局预加载
func OpenRelated() {
	globalIsRelated = true
}

// CloseRelated 关闭全局预加载
func CloseRelated() {
	globalIsRelated = true
}

	`

	genlogic = `{{$obj := .}}{{$list := $obj.Em}}
type _{{$obj.StructName}}Mgr struct {
	*_BaseMgr
}

// {{$obj.StructName}}Mgr open func
func {{$obj.StructName}}Mgr(db *gorm.DB) *_{{$obj.StructName}}Mgr {
	if db == nil {
		panic(fmt.Errorf("{{$obj.StructName}}Mgr need init by db"))
	}
	return &_{{$obj.StructName}}Mgr{_BaseMgr: &_BaseMgr{DB: db, isRelated: globalIsRelated}}
}

// GetTableName get sql table name.获取数据库名字
func (obj *_{{$obj.StructName}}Mgr) GetTableName() string {
	return "{{$obj.TableName}}"
}

// Get 获取
func (obj *_{{$obj.StructName}}Mgr) Get() (result {{$obj.StructName}}, err error) {
	err = obj.DB.Table(obj.GetTableName()).Find(&result).Error
	{{GenPreloadList $obj.PreloadList false}}
	return
}

// Gets 获取批量结果
func (obj *_{{$obj.StructName}}Mgr) Gets() (results []*{{$obj.StructName}}, err error) {
	err = obj.DB.Table(obj.GetTableName()).Find(&results).Error
	{{GenPreloadList $obj.PreloadList true}}
	return
}

//////////////////////////option case ////////////////////////////////////////////
{{range $oem := $obj.Em}}

// With{{$oem.ColStructName}} {{$oem.ColName}}获取 {{$oem.Notes}}
func (obj *_{{$obj.StructName}}Mgr) With{{$oem.ColStructName}}({{CapLowercase $oem.ColStructName}} {{$oem.Type}}) Option {
	return OptionFunc(func(o *SqlContext) *SqlContext { o.Query["{{$oem.ColName}}"] = {{CapLowercase $oem.ColStructName}} ; return o;})
}

func (obj *_{{$obj.StructName}}Mgr) Gt{{$oem.ColStructName}}({{CapLowercase $oem.ColStructName}} {{$oem.Type}}) Option {
	return OptionFunc(func(o *SqlContext) *SqlContext { o.Gt["{{$oem.ColName}}"] = {{CapLowercase $oem.ColStructName}} ; return o;})
}

func (obj *_{{$obj.StructName}}Mgr) Lt{{$oem.ColStructName}}({{CapLowercase $oem.ColStructName}} {{$oem.Type}}) Option {
	return OptionFunc(func(o *SqlContext) *SqlContext { o.Lt["{{$oem.ColName}}"] = {{CapLowercase $oem.ColStructName}} ; return o;})
}

func (obj *_{{$obj.StructName}}Mgr) Not{{$oem.ColStructName}}({{CapLowercase $oem.ColStructName}} {{$oem.Type}}) Option {
	return OptionFunc(func(o *SqlContext) *SqlContext { o.Not["{{$oem.ColName}}"] = {{CapLowercase $oem.ColStructName}} ; return o;})
}

func (obj *_{{$obj.StructName}}Mgr) In{{$oem.ColStructName}}({{CapLowercase $oem.ColStructName}}s []*{{$oem.Type}}) Option {
	return OptionFunc(func(o *SqlContext) *SqlContext { o.In["{{$oem.ColName}}"] = {{CapLowercase $oem.ColStructName}}s ; return o;})
}
{{end}}

//////////////////////////update case ////////////////////////////////////////////
{{range $oem := $obj.Em}}
	// Change{{$oem.ColStructName}} {{$oem.ColName}}更新 {{$oem.Notes}}
	func (obj *_{{$obj.StructName}}Mgr) Change{{$oem.ColStructName}}({{CapLowercase $oem.ColStructName}} {{$oem.Type}}) Option {
		return OptionFunc(func(o *SqlContext) *SqlContext { o.Changes["{{$oem.ColName}}"] = {{CapLowercase $oem.ColStructName}} ; return o;})
	}
{{end}}

//////////////////////////order case ////////////////////////////////////////////
{{range $oem := $obj.Em}}
// Order{{$oem.ColStructName}} {{$oem.ColName}}排序 {{$oem.Notes}}
func (obj *_{{$obj.StructName}}Mgr) Order{{$oem.ColStructName}}(sortType bool) Option {
	orderSql := ""
	if sortType {
		orderSql = "{{$oem.ColName}} asc"
	}else {
		orderSql = "{{$oem.ColName}} desc"
	}
	return OptionFunc(func(o *SqlContext) *SqlContext { o.Order = append(o.Order, orderSql); return o })
} 
{{end}}

// GetByOption 功能选项模式获取
func (obj *_{{$obj.StructName}}Mgr) GetByOption(opts ...Option) (result {{$obj.StructName}}, err error) {
	sqlContext := NewSqlContext()
	for _, o := range opts {
		o.Apply(sqlContext)
	}

	//Query
	db := obj.DB.Table(obj.GetTableName()).Where(sqlContext.Query)

	if len(sqlContext.In) >0 {
		for col, val := range sqlContext.In {
			db = db.Where(col+" in (?)", val)
		}
	}
	if len(sqlContext.Not) >0 {
		for col, val := range sqlContext.Not {
			db = db.Where(col+" <> ?", val)
		}
	}
	if len(sqlContext.Lt) >0 {
		for col, val := range sqlContext.Lt {
			db = db.Where(col+" < ? ", val)
		}
	}
	if len(sqlContext.Gt) >0 {
		for col, val := range sqlContext.Gt {
			db = db.Where(col+" > ?", val)
		}
	}

	//Order
	if len(sqlContext.Order) > 0{
		for _, order := range sqlContext.Order {
			db = db.Order(order)
		}
	}

	//find
	err = db.Find(&result).Error

	{{GenPreloadList $obj.PreloadList false}}
	return
}

// GetByOptions 批量功能选项模式获取
func (obj *_{{$obj.StructName}}Mgr) GetByOptions(opts ...Option) (results []*{{$obj.StructName}}, err error) {
	sqlContext := NewSqlContext()
	for _, o := range opts {
		o.Apply(sqlContext)
	}

	//Query
	db := obj.DB.Table(obj.GetTableName()).Where(sqlContext.Query)

	if len(sqlContext.In) >0 {
		for col, val := range sqlContext.In {
			db = db.Where(col+" in (?)", val)
		}
	}
	if len(sqlContext.Not) >0 {
		for col, val := range sqlContext.Not {
			db = db.Where(col+" <> ?", val)
		}
	}
	if len(sqlContext.Lt) >0 {
		for col, val := range sqlContext.Lt {
			db = db.Where(col+" < ? ", val)
		}
	}
	if len(sqlContext.Gt) >0 {
		for col, val := range sqlContext.Gt {
			db = db.Where(col+" > ?", val)
		}
	}

	//Order
	if len(sqlContext.Order) > 0{
		for _, order := range sqlContext.Order {
			db = db.Order(order)
		}
	}

	//find
	err = db.Find(&results).Error

	{{GenPreloadList $obj.PreloadList true}}
	return
}

func (obj *_{{$obj.StructName}}Mgr) GetPageByOptions(pageIndex, pageSize int, opts ...Option) (results []*{{$obj.StructName}}, err error) {
	sqlContext := NewSqlContext()
	for _, o := range opts {
		o.Apply(sqlContext)
	}
	//Query
	db := obj.DB.Table(obj.GetTableName()).Where(sqlContext.Query)

	if len(sqlContext.In) >0 {
		for col, val := range sqlContext.In {
			db = db.Where(col+" in (?)", val)
		}
	}
	if len(sqlContext.Not) >0 {
		for col, val := range sqlContext.Not {
			db = db.Where(col+" <> ?", val)
		}
	}
	if len(sqlContext.Lt) >0 {
		for col, val := range sqlContext.Lt {
			db = db.Where(col+" < ? ", val)
		}
	}
	if len(sqlContext.Gt) >0 {
		for col, val := range sqlContext.Gt {
			db = db.Where(col+" > ?", val)
		}
	}

	//Order
	if len(sqlContext.Order) > 0{
		for _, order := range sqlContext.Order {
			db = db.Order(order)
		}
	}

	//offset
	db = db.Offset(pageIndex * pageSize)

	//limit
	db = db.Limit(pageSize)

	//find
	err = db.Find(&results).Error
	return
}

func (obj *_{{$obj.StructName}}Mgr) UpdateByOption(opts ...Option) (err error) {
    sqlContext := NewSqlContext()
	for _, o := range opts {
		o.Apply(sqlContext)
	}
	err = obj.DB.Table(obj.GetTableName()).Where(sqlContext.Query).Updates(sqlContext.Changes).Error
	return
}

// Count 计数
func (obj *_{{$obj.StructName}}Mgr) Count(opts ...Option) (num int, err error) {
	sqlContext :=  NewSqlContext()
	for _, o := range opts {
		o.Apply(sqlContext)
	}
	//Query
	db := obj.DB.Table(obj.GetTableName()).Where(sqlContext.Query)

	if len(sqlContext.In) >0 {
		for col, val := range sqlContext.In {
			db = db.Where(col+" in (?)", val)
		}
	}
	if len(sqlContext.Not) >0 {
		for col, val := range sqlContext.Not {
			db = db.Where(col+" <> ?", val)
		}
	}
	if len(sqlContext.Lt) >0 {
		for col, val := range sqlContext.Lt {
			db = db.Where(col+" < ? ", val)
		}
	}
	if len(sqlContext.Gt) >0 {
		for col, val := range sqlContext.Gt {
			db = db.Where(col+" > ?", val)
		}
	}
	err = obj.DB.Table(obj.GetTableName()).Where(sqlContext.Query).Count(&num).Error
	return
}

// Delete 计数
func (obj *_{{$obj.StructName}}Mgr) Delete(opts ...Option) (err error) {
	sqlContext := NewSqlContext()

	for _, o := range opts {
		o.Apply(sqlContext)
	}
	//Query
	db := obj.DB.Table(obj.GetTableName()).Where(sqlContext.Query)

	if len(sqlContext.In) >0 {
		for col, val := range sqlContext.In {
			db = db.Where(col+" in (?)", val)
		}
	}
	if len(sqlContext.Not) >0 {
		for col, val := range sqlContext.Not {
			db = db.Where(col+" <> ?", val)
		}
	}
	if len(sqlContext.Lt) >0 {
		for col, val := range sqlContext.Lt {
			db = db.Where(col+" < ? ", val)
		}
	}
	if len(sqlContext.Gt) >0 {
		for col, val := range sqlContext.Gt {
			db = db.Where(col+" > ?", val)
		}
	}
	err = obj.DB.Table(obj.GetTableName()).Where(sqlContext.Query).Delete(Sector{}).Error
	return
}
//////////////////////////batch case ////////////////////////////////////////////

{{range $oem := $obj.Em}}
// GetFrom{{$oem.ColStructName}} 通过{{$oem.ColName}}获取内容 {{$oem.Notes}} {{if $oem.IsMulti}}
func (obj *_{{$obj.StructName}}Mgr) GetFrom{{$oem.ColStructName}}({{CapLowercase $oem.ColStructName}} {{$oem.Type}}) (batchResults []*{{$obj.StructName}}, err error) {
	err = obj.DB.Table(obj.GetTableName()).Where("{{$oem.ColName}} = ?", {{CapLowercase $oem.ColStructName}}).Find(&batchResults).Error
	{{GenPreloadList $obj.PreloadList true}}
	return
}
{{else}}
func (obj *_{{$obj.StructName}}Mgr)  GetFrom{{$oem.ColStructName}}({{CapLowercase $oem.ColStructName}} {{$oem.Type}}) (result {{$obj.StructName}}, err error) {
	err = obj.DB.Table(obj.GetTableName()).Where("{{$oem.ColName}} = ?", {{CapLowercase $oem.ColStructName}}).Find(&result).Error
	{{GenPreloadList $obj.PreloadList false}}
	return
}
{{end}}
// GetBatchFrom{{$oem.ColStructName}} 批量唯一主键查找 {{$oem.Notes}}
func (obj *_{{$obj.StructName}}Mgr) GetBatchFrom{{$oem.ColStructName}}({{CapLowercase $oem.ColStructName}}s []{{$oem.Type}}) (batchResults []*{{$obj.StructName}}, err error) {
	err = obj.DB.Table(obj.GetTableName()).Where("{{$oem.ColName}} IN (?)", {{CapLowercase $oem.ColStructName}}s).Find(&batchResults).Error
	{{GenPreloadList $obj.PreloadList true}}
	return
}
 {{end}}
 //////////////////////////primary index case ////////////////////////////////////////////
 {{range $ofm := $obj.Primay}}
 // {{GenFListIndex $ofm 1}} primay or index 获取唯一内容
 func (obj *_{{$obj.StructName}}Mgr) {{GenFListIndex $ofm 1}}({{GenFListIndex $ofm 2}}) (result {{$obj.StructName}}, err error) {
	err = obj.DB.Table(obj.GetTableName()).Where("{{GenFListIndex $ofm 3}}", {{GenFListIndex $ofm 4}}).Find(&result).Error
	{{GenPreloadList $obj.PreloadList false}}
	return
}
 {{end}}

 {{range $ofm := $obj.Index}}
 // {{GenFListIndex $ofm 1}}  获取多个内容
 func (obj *_{{$obj.StructName}}Mgr) {{GenFListIndex $ofm 1}}({{GenFListIndex $ofm 2}}) (results []*{{$obj.StructName}}, err error) {
	err = obj.DB.Table(obj.GetTableName()).Where("{{GenFListIndex $ofm 3}}", {{GenFListIndex $ofm 4}}).Find(&results).Error
	{{GenPreloadList $obj.PreloadList true}}
	return
}
 {{end}}

`
	genPreload = `if err == nil && obj.isRelated { {{range $obj := .}}{{if $obj.IsMulti}}
		{
			var info []{{$obj.ForeignkeyStructName}}  // {{$obj.Notes}} 
			err = obj.DB.New().Table("{{$obj.ForeignkeyTableName}}").Where("{{$obj.ForeignkeyCol}} = ?", result.{{$obj.ColStructName}}).Find(&info).Error
			if err != nil {
				return
			}
			result.{{$obj.ForeignkeyStructName}}List = info
		}  {{else}} 
		{
			var info {{$obj.ForeignkeyStructName}}  // {{$obj.Notes}} 
			err = obj.DB.New().Table("{{$obj.ForeignkeyTableName}}").Where("{{$obj.ForeignkeyCol}} = ?", result.{{$obj.ColStructName}}).Find(&info).Error
			if err != nil {
				return
			}
			result.{{$obj.ForeignkeyStructName}} = info
		} {{end}} {{end}}
	}
`
	genPreloadMulti = `if err == nil && obj.isRelated {
		for i := 0; i < len(results); i++ { {{range $obj := .}}{{if $obj.IsMulti}}
		{
			var info []{{$obj.ForeignkeyStructName}}  // {{$obj.Notes}} 
			err = obj.DB.New().Table("{{$obj.ForeignkeyTableName}}").Where("{{$obj.ForeignkeyCol}} = ?", results[i].{{$obj.ColStructName}}).Find(&info).Error
			if err != nil {
				return
			}
			results[i].{{$obj.ForeignkeyStructName}}List = info
		}  {{else}} 
		{
			var info {{$obj.ForeignkeyStructName}}  // {{$obj.Notes}} 
			err = obj.DB.New().Table("{{$obj.ForeignkeyTableName}}").Where("{{$obj.ForeignkeyCol}} = ?", results[i].{{$obj.ColStructName}}).Find(&info).Error
			if err != nil {
				return
			}
			results[i].{{$obj.ForeignkeyStructName}} = info
		} {{end}} {{end}}
	}
}`
)
