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

func newHomepage(db *gorm.DB, opts ...gen.DOOption) homepage {
	_homepage := homepage{}

	_homepage.homepageDo.UseDB(db, opts...)
	_homepage.homepageDo.UseModel(&po.Homepage{})

	tableName := _homepage.homepageDo.TableName()
	_homepage.ALL = field.NewAsterisk(tableName)
	_homepage.ID = field.NewUint(tableName, "id")
	_homepage.CreatedAt = field.NewTime(tableName, "created_at")
	_homepage.UpdatedAt = field.NewTime(tableName, "updated_at")
	_homepage.DeletedAt = field.NewField(tableName, "deleted_at")
	_homepage.Name = field.NewString(tableName, "name")
	_homepage.HomepageURL = field.NewString(tableName, "homepage_url")
	_homepage.Actions = field.NewString(tableName, "version_actions")
	_homepage.NoUpdate = field.NewBool(tableName, "no_update")
	_homepage.Current = homepageHasOneCurrent{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("Current", "po.CurrentVersion"),
		Version: struct {
			field.RelationField
			CV struct {
				field.RelationField
			}
		}{
			RelationField: field.NewRelation("Current.Version", "po.Version"),
			CV: struct {
				field.RelationField
			}{
				RelationField: field.NewRelation("Current.Version.CV", "po.CurrentVersion"),
			},
		},
	}

	_homepage.Versions = homepageHasManyVersions{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("Versions", "po.Version"),
	}

	_homepage.fillFieldMap()

	return _homepage
}

type homepage struct {
	homepageDo homepageDo

	ALL         field.Asterisk
	ID          field.Uint
	CreatedAt   field.Time
	UpdatedAt   field.Time
	DeletedAt   field.Field
	Name        field.String
	HomepageURL field.String
	Actions     field.String
	NoUpdate    field.Bool
	Current     homepageHasOneCurrent

	Versions homepageHasManyVersions

	fieldMap map[string]field.Expr
}

func (h homepage) Table(newTableName string) *homepage {
	h.homepageDo.UseTable(newTableName)
	return h.updateTableName(newTableName)
}

func (h homepage) As(alias string) *homepage {
	h.homepageDo.DO = *(h.homepageDo.As(alias).(*gen.DO))
	return h.updateTableName(alias)
}

func (h *homepage) updateTableName(table string) *homepage {
	h.ALL = field.NewAsterisk(table)
	h.ID = field.NewUint(table, "id")
	h.CreatedAt = field.NewTime(table, "created_at")
	h.UpdatedAt = field.NewTime(table, "updated_at")
	h.DeletedAt = field.NewField(table, "deleted_at")
	h.Name = field.NewString(table, "name")
	h.HomepageURL = field.NewString(table, "homepage_url")
	h.Actions = field.NewString(table, "version_actions")
	h.NoUpdate = field.NewBool(table, "no_update")

	h.fillFieldMap()

	return h
}

func (h *homepage) WithContext(ctx context.Context) IHomepageDo { return h.homepageDo.WithContext(ctx) }

func (h homepage) TableName() string { return h.homepageDo.TableName() }

func (h homepage) Alias() string { return h.homepageDo.Alias() }

func (h homepage) Columns(cols ...field.Expr) gen.Columns { return h.homepageDo.Columns(cols...) }

func (h *homepage) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := h.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (h *homepage) fillFieldMap() {
	h.fieldMap = make(map[string]field.Expr, 10)
	h.fieldMap["id"] = h.ID
	h.fieldMap["created_at"] = h.CreatedAt
	h.fieldMap["updated_at"] = h.UpdatedAt
	h.fieldMap["deleted_at"] = h.DeletedAt
	h.fieldMap["name"] = h.Name
	h.fieldMap["homepage_url"] = h.HomepageURL
	h.fieldMap["version_actions"] = h.Actions
	h.fieldMap["no_update"] = h.NoUpdate

}

func (h homepage) clone(db *gorm.DB) homepage {
	h.homepageDo.ReplaceConnPool(db.Statement.ConnPool)
	return h
}

func (h homepage) replaceDB(db *gorm.DB) homepage {
	h.homepageDo.ReplaceDB(db)
	return h
}

type homepageHasOneCurrent struct {
	db *gorm.DB

	field.RelationField

	Version struct {
		field.RelationField
		CV struct {
			field.RelationField
		}
	}
}

func (a homepageHasOneCurrent) Where(conds ...field.Expr) *homepageHasOneCurrent {
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

func (a homepageHasOneCurrent) WithContext(ctx context.Context) *homepageHasOneCurrent {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a homepageHasOneCurrent) Session(session *gorm.Session) *homepageHasOneCurrent {
	a.db = a.db.Session(session)
	return &a
}

func (a homepageHasOneCurrent) Model(m *po.Homepage) *homepageHasOneCurrentTx {
	return &homepageHasOneCurrentTx{a.db.Model(m).Association(a.Name())}
}

type homepageHasOneCurrentTx struct{ tx *gorm.Association }

func (a homepageHasOneCurrentTx) Find() (result *po.CurrentVersion, err error) {
	return result, a.tx.Find(&result)
}

func (a homepageHasOneCurrentTx) Append(values ...*po.CurrentVersion) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a homepageHasOneCurrentTx) Replace(values ...*po.CurrentVersion) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a homepageHasOneCurrentTx) Delete(values ...*po.CurrentVersion) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a homepageHasOneCurrentTx) Clear() error {
	return a.tx.Clear()
}

func (a homepageHasOneCurrentTx) Count() int64 {
	return a.tx.Count()
}

type homepageHasManyVersions struct {
	db *gorm.DB

	field.RelationField
}

func (a homepageHasManyVersions) Where(conds ...field.Expr) *homepageHasManyVersions {
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

func (a homepageHasManyVersions) WithContext(ctx context.Context) *homepageHasManyVersions {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a homepageHasManyVersions) Session(session *gorm.Session) *homepageHasManyVersions {
	a.db = a.db.Session(session)
	return &a
}

func (a homepageHasManyVersions) Model(m *po.Homepage) *homepageHasManyVersionsTx {
	return &homepageHasManyVersionsTx{a.db.Model(m).Association(a.Name())}
}

type homepageHasManyVersionsTx struct{ tx *gorm.Association }

func (a homepageHasManyVersionsTx) Find() (result []*po.Version, err error) {
	return result, a.tx.Find(&result)
}

func (a homepageHasManyVersionsTx) Append(values ...*po.Version) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a homepageHasManyVersionsTx) Replace(values ...*po.Version) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a homepageHasManyVersionsTx) Delete(values ...*po.Version) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a homepageHasManyVersionsTx) Clear() error {
	return a.tx.Clear()
}

func (a homepageHasManyVersionsTx) Count() int64 {
	return a.tx.Count()
}

type homepageDo struct{ gen.DO }

type IHomepageDo interface {
	gen.SubQuery
	Debug() IHomepageDo
	WithContext(ctx context.Context) IHomepageDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IHomepageDo
	WriteDB() IHomepageDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IHomepageDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IHomepageDo
	Not(conds ...gen.Condition) IHomepageDo
	Or(conds ...gen.Condition) IHomepageDo
	Select(conds ...field.Expr) IHomepageDo
	Where(conds ...gen.Condition) IHomepageDo
	Order(conds ...field.Expr) IHomepageDo
	Distinct(cols ...field.Expr) IHomepageDo
	Omit(cols ...field.Expr) IHomepageDo
	Join(table schema.Tabler, on ...field.Expr) IHomepageDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IHomepageDo
	RightJoin(table schema.Tabler, on ...field.Expr) IHomepageDo
	Group(cols ...field.Expr) IHomepageDo
	Having(conds ...gen.Condition) IHomepageDo
	Limit(limit int) IHomepageDo
	Offset(offset int) IHomepageDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IHomepageDo
	Unscoped() IHomepageDo
	Create(values ...*po.Homepage) error
	CreateInBatches(values []*po.Homepage, batchSize int) error
	Save(values ...*po.Homepage) error
	First() (*po.Homepage, error)
	Take() (*po.Homepage, error)
	Last() (*po.Homepage, error)
	Find() ([]*po.Homepage, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*po.Homepage, err error)
	FindInBatches(result *[]*po.Homepage, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*po.Homepage) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IHomepageDo
	Assign(attrs ...field.AssignExpr) IHomepageDo
	Joins(fields ...field.RelationField) IHomepageDo
	Preload(fields ...field.RelationField) IHomepageDo
	FirstOrInit() (*po.Homepage, error)
	FirstOrCreate() (*po.Homepage, error)
	FindByPage(offset int, limit int) (result []*po.Homepage, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IHomepageDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (h homepageDo) Debug() IHomepageDo {
	return h.withDO(h.DO.Debug())
}

func (h homepageDo) WithContext(ctx context.Context) IHomepageDo {
	return h.withDO(h.DO.WithContext(ctx))
}

func (h homepageDo) ReadDB() IHomepageDo {
	return h.Clauses(dbresolver.Read)
}

func (h homepageDo) WriteDB() IHomepageDo {
	return h.Clauses(dbresolver.Write)
}

func (h homepageDo) Session(config *gorm.Session) IHomepageDo {
	return h.withDO(h.DO.Session(config))
}

func (h homepageDo) Clauses(conds ...clause.Expression) IHomepageDo {
	return h.withDO(h.DO.Clauses(conds...))
}

func (h homepageDo) Returning(value interface{}, columns ...string) IHomepageDo {
	return h.withDO(h.DO.Returning(value, columns...))
}

func (h homepageDo) Not(conds ...gen.Condition) IHomepageDo {
	return h.withDO(h.DO.Not(conds...))
}

func (h homepageDo) Or(conds ...gen.Condition) IHomepageDo {
	return h.withDO(h.DO.Or(conds...))
}

func (h homepageDo) Select(conds ...field.Expr) IHomepageDo {
	return h.withDO(h.DO.Select(conds...))
}

func (h homepageDo) Where(conds ...gen.Condition) IHomepageDo {
	return h.withDO(h.DO.Where(conds...))
}

func (h homepageDo) Order(conds ...field.Expr) IHomepageDo {
	return h.withDO(h.DO.Order(conds...))
}

func (h homepageDo) Distinct(cols ...field.Expr) IHomepageDo {
	return h.withDO(h.DO.Distinct(cols...))
}

func (h homepageDo) Omit(cols ...field.Expr) IHomepageDo {
	return h.withDO(h.DO.Omit(cols...))
}

func (h homepageDo) Join(table schema.Tabler, on ...field.Expr) IHomepageDo {
	return h.withDO(h.DO.Join(table, on...))
}

func (h homepageDo) LeftJoin(table schema.Tabler, on ...field.Expr) IHomepageDo {
	return h.withDO(h.DO.LeftJoin(table, on...))
}

func (h homepageDo) RightJoin(table schema.Tabler, on ...field.Expr) IHomepageDo {
	return h.withDO(h.DO.RightJoin(table, on...))
}

func (h homepageDo) Group(cols ...field.Expr) IHomepageDo {
	return h.withDO(h.DO.Group(cols...))
}

func (h homepageDo) Having(conds ...gen.Condition) IHomepageDo {
	return h.withDO(h.DO.Having(conds...))
}

func (h homepageDo) Limit(limit int) IHomepageDo {
	return h.withDO(h.DO.Limit(limit))
}

func (h homepageDo) Offset(offset int) IHomepageDo {
	return h.withDO(h.DO.Offset(offset))
}

func (h homepageDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IHomepageDo {
	return h.withDO(h.DO.Scopes(funcs...))
}

func (h homepageDo) Unscoped() IHomepageDo {
	return h.withDO(h.DO.Unscoped())
}

func (h homepageDo) Create(values ...*po.Homepage) error {
	if len(values) == 0 {
		return nil
	}
	return h.DO.Create(values)
}

func (h homepageDo) CreateInBatches(values []*po.Homepage, batchSize int) error {
	return h.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (h homepageDo) Save(values ...*po.Homepage) error {
	if len(values) == 0 {
		return nil
	}
	return h.DO.Save(values)
}

func (h homepageDo) First() (*po.Homepage, error) {
	if result, err := h.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*po.Homepage), nil
	}
}

func (h homepageDo) Take() (*po.Homepage, error) {
	if result, err := h.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*po.Homepage), nil
	}
}

func (h homepageDo) Last() (*po.Homepage, error) {
	if result, err := h.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*po.Homepage), nil
	}
}

func (h homepageDo) Find() ([]*po.Homepage, error) {
	result, err := h.DO.Find()
	return result.([]*po.Homepage), err
}

func (h homepageDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*po.Homepage, err error) {
	buf := make([]*po.Homepage, 0, batchSize)
	err = h.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (h homepageDo) FindInBatches(result *[]*po.Homepage, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return h.DO.FindInBatches(result, batchSize, fc)
}

func (h homepageDo) Attrs(attrs ...field.AssignExpr) IHomepageDo {
	return h.withDO(h.DO.Attrs(attrs...))
}

func (h homepageDo) Assign(attrs ...field.AssignExpr) IHomepageDo {
	return h.withDO(h.DO.Assign(attrs...))
}

func (h homepageDo) Joins(fields ...field.RelationField) IHomepageDo {
	for _, _f := range fields {
		h = *h.withDO(h.DO.Joins(_f))
	}
	return &h
}

func (h homepageDo) Preload(fields ...field.RelationField) IHomepageDo {
	for _, _f := range fields {
		h = *h.withDO(h.DO.Preload(_f))
	}
	return &h
}

func (h homepageDo) FirstOrInit() (*po.Homepage, error) {
	if result, err := h.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*po.Homepage), nil
	}
}

func (h homepageDo) FirstOrCreate() (*po.Homepage, error) {
	if result, err := h.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*po.Homepage), nil
	}
}

func (h homepageDo) FindByPage(offset int, limit int) (result []*po.Homepage, count int64, err error) {
	result, err = h.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = h.Offset(-1).Limit(-1).Count()
	return
}

func (h homepageDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = h.Count()
	if err != nil {
		return
	}

	err = h.Offset(offset).Limit(limit).Scan(result)
	return
}

func (h homepageDo) Scan(result interface{}) (err error) {
	return h.DO.Scan(result)
}

func (h homepageDo) Delete(models ...*po.Homepage) (result gen.ResultInfo, err error) {
	return h.DO.Delete(models)
}

func (h *homepageDo) withDO(do gen.Dao) *homepageDo {
	h.DO = *do.(*gen.DO)
	return h
}
