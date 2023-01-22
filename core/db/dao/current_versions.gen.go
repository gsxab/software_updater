// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"software_updater/core/db/po"
)

func newCurrentVersion(db *gorm.DB, opts ...gen.DOOption) currentVersion {
	_currentVersion := currentVersion{}

	_currentVersion.currentVersionDo.UseDB(db, opts...)
	_currentVersion.currentVersionDo.UseModel(&po.CurrentVersion{})

	tableName := _currentVersion.currentVersionDo.TableName()
	_currentVersion.ALL = field.NewAsterisk(tableName)
	_currentVersion.ID = field.NewUint(tableName, "id")
	_currentVersion.CreatedAt = field.NewTime(tableName, "created_at")
	_currentVersion.UpdatedAt = field.NewTime(tableName, "updated_at")
	_currentVersion.DeletedAt = field.NewField(tableName, "deleted_at")
	_currentVersion.Name = field.NewString(tableName, "name")
	_currentVersion.ScheduledAt = field.NewTime(tableName, "scheduled_at")
	_currentVersion.CurrentVersionID = field.NewUint(tableName, "current_version_id")
	_currentVersion.Info = field.NewString(tableName, "info")
	_currentVersion.CurrentVersion = currentVersionHasOneCurrentVersion{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("CurrentVersion", "po.Version"),
	}

	_currentVersion.fillFieldMap()

	return _currentVersion
}

type currentVersion struct {
	currentVersionDo currentVersionDo

	ALL              field.Asterisk
	ID               field.Uint
	CreatedAt        field.Time
	UpdatedAt        field.Time
	DeletedAt        field.Field
	Name             field.String
	ScheduledAt      field.Time
	CurrentVersionID field.Uint
	Info             field.String
	CurrentVersion   currentVersionHasOneCurrentVersion

	fieldMap map[string]field.Expr
}

func (c currentVersion) Table(newTableName string) *currentVersion {
	c.currentVersionDo.UseTable(newTableName)
	return c.updateTableName(newTableName)
}

func (c currentVersion) As(alias string) *currentVersion {
	c.currentVersionDo.DO = *(c.currentVersionDo.As(alias).(*gen.DO))
	return c.updateTableName(alias)
}

func (c *currentVersion) updateTableName(table string) *currentVersion {
	c.ALL = field.NewAsterisk(table)
	c.ID = field.NewUint(table, "id")
	c.CreatedAt = field.NewTime(table, "created_at")
	c.UpdatedAt = field.NewTime(table, "updated_at")
	c.DeletedAt = field.NewField(table, "deleted_at")
	c.Name = field.NewString(table, "name")
	c.ScheduledAt = field.NewTime(table, "scheduled_at")
	c.CurrentVersionID = field.NewUint(table, "current_version_id")
	c.Info = field.NewString(table, "info")

	c.fillFieldMap()

	return c
}

func (c *currentVersion) WithContext(ctx context.Context) ICurrentVersionDo {
	return c.currentVersionDo.WithContext(ctx)
}

func (c currentVersion) TableName() string { return c.currentVersionDo.TableName() }

func (c currentVersion) Alias() string { return c.currentVersionDo.Alias() }

func (c *currentVersion) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := c.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (c *currentVersion) fillFieldMap() {
	c.fieldMap = make(map[string]field.Expr, 9)
	c.fieldMap["id"] = c.ID
	c.fieldMap["created_at"] = c.CreatedAt
	c.fieldMap["updated_at"] = c.UpdatedAt
	c.fieldMap["deleted_at"] = c.DeletedAt
	c.fieldMap["name"] = c.Name
	c.fieldMap["scheduled_at"] = c.ScheduledAt
	c.fieldMap["current_version_id"] = c.CurrentVersionID
	c.fieldMap["info"] = c.Info

}

func (c currentVersion) clone(db *gorm.DB) currentVersion {
	c.currentVersionDo.ReplaceConnPool(db.Statement.ConnPool)
	return c
}

func (c currentVersion) replaceDB(db *gorm.DB) currentVersion {
	c.currentVersionDo.ReplaceDB(db)
	return c
}

type currentVersionHasOneCurrentVersion struct {
	db *gorm.DB

	field.RelationField
}

func (a currentVersionHasOneCurrentVersion) Where(conds ...field.Expr) *currentVersionHasOneCurrentVersion {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a currentVersionHasOneCurrentVersion) WithContext(ctx context.Context) *currentVersionHasOneCurrentVersion {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a currentVersionHasOneCurrentVersion) Model(m *po.CurrentVersion) *currentVersionHasOneCurrentVersionTx {
	return &currentVersionHasOneCurrentVersionTx{a.db.Model(m).Association(a.Name())}
}

type currentVersionHasOneCurrentVersionTx struct{ tx *gorm.Association }

func (a currentVersionHasOneCurrentVersionTx) Find() (result *po.Version, err error) {
	return result, a.tx.Find(&result)
}

func (a currentVersionHasOneCurrentVersionTx) Append(values ...*po.Version) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a currentVersionHasOneCurrentVersionTx) Replace(values ...*po.Version) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a currentVersionHasOneCurrentVersionTx) Delete(values ...*po.Version) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a currentVersionHasOneCurrentVersionTx) Clear() error {
	return a.tx.Clear()
}

func (a currentVersionHasOneCurrentVersionTx) Count() int64 {
	return a.tx.Count()
}

type currentVersionDo struct{ gen.DO }

type ICurrentVersionDo interface {
	gen.SubQuery
	Debug() ICurrentVersionDo
	WithContext(ctx context.Context) ICurrentVersionDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ICurrentVersionDo
	WriteDB() ICurrentVersionDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ICurrentVersionDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ICurrentVersionDo
	Not(conds ...gen.Condition) ICurrentVersionDo
	Or(conds ...gen.Condition) ICurrentVersionDo
	Select(conds ...field.Expr) ICurrentVersionDo
	Where(conds ...gen.Condition) ICurrentVersionDo
	Order(conds ...field.Expr) ICurrentVersionDo
	Distinct(cols ...field.Expr) ICurrentVersionDo
	Omit(cols ...field.Expr) ICurrentVersionDo
	Join(table schema.Tabler, on ...field.Expr) ICurrentVersionDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ICurrentVersionDo
	RightJoin(table schema.Tabler, on ...field.Expr) ICurrentVersionDo
	Group(cols ...field.Expr) ICurrentVersionDo
	Having(conds ...gen.Condition) ICurrentVersionDo
	Limit(limit int) ICurrentVersionDo
	Offset(offset int) ICurrentVersionDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ICurrentVersionDo
	Unscoped() ICurrentVersionDo
	Create(values ...*po.CurrentVersion) error
	CreateInBatches(values []*po.CurrentVersion, batchSize int) error
	Save(values ...*po.CurrentVersion) error
	First() (*po.CurrentVersion, error)
	Take() (*po.CurrentVersion, error)
	Last() (*po.CurrentVersion, error)
	Find() ([]*po.CurrentVersion, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*po.CurrentVersion, err error)
	FindInBatches(result *[]*po.CurrentVersion, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*po.CurrentVersion) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ICurrentVersionDo
	Assign(attrs ...field.AssignExpr) ICurrentVersionDo
	Joins(fields ...field.RelationField) ICurrentVersionDo
	Preload(fields ...field.RelationField) ICurrentVersionDo
	FirstOrInit() (*po.CurrentVersion, error)
	FirstOrCreate() (*po.CurrentVersion, error)
	FindByPage(offset int, limit int) (result []*po.CurrentVersion, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ICurrentVersionDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (c currentVersionDo) Debug() ICurrentVersionDo {
	return c.withDO(c.DO.Debug())
}

func (c currentVersionDo) WithContext(ctx context.Context) ICurrentVersionDo {
	return c.withDO(c.DO.WithContext(ctx))
}

func (c currentVersionDo) ReadDB() ICurrentVersionDo {
	return c.Clauses(dbresolver.Read)
}

func (c currentVersionDo) WriteDB() ICurrentVersionDo {
	return c.Clauses(dbresolver.Write)
}

func (c currentVersionDo) Session(config *gorm.Session) ICurrentVersionDo {
	return c.withDO(c.DO.Session(config))
}

func (c currentVersionDo) Clauses(conds ...clause.Expression) ICurrentVersionDo {
	return c.withDO(c.DO.Clauses(conds...))
}

func (c currentVersionDo) Returning(value interface{}, columns ...string) ICurrentVersionDo {
	return c.withDO(c.DO.Returning(value, columns...))
}

func (c currentVersionDo) Not(conds ...gen.Condition) ICurrentVersionDo {
	return c.withDO(c.DO.Not(conds...))
}

func (c currentVersionDo) Or(conds ...gen.Condition) ICurrentVersionDo {
	return c.withDO(c.DO.Or(conds...))
}

func (c currentVersionDo) Select(conds ...field.Expr) ICurrentVersionDo {
	return c.withDO(c.DO.Select(conds...))
}

func (c currentVersionDo) Where(conds ...gen.Condition) ICurrentVersionDo {
	return c.withDO(c.DO.Where(conds...))
}

func (c currentVersionDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) ICurrentVersionDo {
	return c.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (c currentVersionDo) Order(conds ...field.Expr) ICurrentVersionDo {
	return c.withDO(c.DO.Order(conds...))
}

func (c currentVersionDo) Distinct(cols ...field.Expr) ICurrentVersionDo {
	return c.withDO(c.DO.Distinct(cols...))
}

func (c currentVersionDo) Omit(cols ...field.Expr) ICurrentVersionDo {
	return c.withDO(c.DO.Omit(cols...))
}

func (c currentVersionDo) Join(table schema.Tabler, on ...field.Expr) ICurrentVersionDo {
	return c.withDO(c.DO.Join(table, on...))
}

func (c currentVersionDo) LeftJoin(table schema.Tabler, on ...field.Expr) ICurrentVersionDo {
	return c.withDO(c.DO.LeftJoin(table, on...))
}

func (c currentVersionDo) RightJoin(table schema.Tabler, on ...field.Expr) ICurrentVersionDo {
	return c.withDO(c.DO.RightJoin(table, on...))
}

func (c currentVersionDo) Group(cols ...field.Expr) ICurrentVersionDo {
	return c.withDO(c.DO.Group(cols...))
}

func (c currentVersionDo) Having(conds ...gen.Condition) ICurrentVersionDo {
	return c.withDO(c.DO.Having(conds...))
}

func (c currentVersionDo) Limit(limit int) ICurrentVersionDo {
	return c.withDO(c.DO.Limit(limit))
}

func (c currentVersionDo) Offset(offset int) ICurrentVersionDo {
	return c.withDO(c.DO.Offset(offset))
}

func (c currentVersionDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ICurrentVersionDo {
	return c.withDO(c.DO.Scopes(funcs...))
}

func (c currentVersionDo) Unscoped() ICurrentVersionDo {
	return c.withDO(c.DO.Unscoped())
}

func (c currentVersionDo) Create(values ...*po.CurrentVersion) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Create(values)
}

func (c currentVersionDo) CreateInBatches(values []*po.CurrentVersion, batchSize int) error {
	return c.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (c currentVersionDo) Save(values ...*po.CurrentVersion) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Save(values)
}

func (c currentVersionDo) First() (*po.CurrentVersion, error) {
	if result, err := c.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*po.CurrentVersion), nil
	}
}

func (c currentVersionDo) Take() (*po.CurrentVersion, error) {
	if result, err := c.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*po.CurrentVersion), nil
	}
}

func (c currentVersionDo) Last() (*po.CurrentVersion, error) {
	if result, err := c.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*po.CurrentVersion), nil
	}
}

func (c currentVersionDo) Find() ([]*po.CurrentVersion, error) {
	result, err := c.DO.Find()
	return result.([]*po.CurrentVersion), err
}

func (c currentVersionDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*po.CurrentVersion, err error) {
	buf := make([]*po.CurrentVersion, 0, batchSize)
	err = c.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (c currentVersionDo) FindInBatches(result *[]*po.CurrentVersion, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return c.DO.FindInBatches(result, batchSize, fc)
}

func (c currentVersionDo) Attrs(attrs ...field.AssignExpr) ICurrentVersionDo {
	return c.withDO(c.DO.Attrs(attrs...))
}

func (c currentVersionDo) Assign(attrs ...field.AssignExpr) ICurrentVersionDo {
	return c.withDO(c.DO.Assign(attrs...))
}

func (c currentVersionDo) Joins(fields ...field.RelationField) ICurrentVersionDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Joins(_f))
	}
	return &c
}

func (c currentVersionDo) Preload(fields ...field.RelationField) ICurrentVersionDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Preload(_f))
	}
	return &c
}

func (c currentVersionDo) FirstOrInit() (*po.CurrentVersion, error) {
	if result, err := c.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*po.CurrentVersion), nil
	}
}

func (c currentVersionDo) FirstOrCreate() (*po.CurrentVersion, error) {
	if result, err := c.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*po.CurrentVersion), nil
	}
}

func (c currentVersionDo) FindByPage(offset int, limit int) (result []*po.CurrentVersion, count int64, err error) {
	result, err = c.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = c.Offset(-1).Limit(-1).Count()
	return
}

func (c currentVersionDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = c.Count()
	if err != nil {
		return
	}

	err = c.Offset(offset).Limit(limit).Scan(result)
	return
}

func (c currentVersionDo) Scan(result interface{}) (err error) {
	return c.DO.Scan(result)
}

func (c currentVersionDo) Delete(models ...*po.CurrentVersion) (result gen.ResultInfo, err error) {
	return c.DO.Delete(models)
}

func (c *currentVersionDo) withDO(do gen.Dao) *currentVersionDo {
	c.DO = *do.(*gen.DO)
	return c
}
