package query

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Option 查询选项函数类型
type Option func(*Options)

// Options 查询选项
type Options struct {
	filter     any
	limit      int
	offset     int
	clauses    []clause.Expression
	preloads   []string
	selects    []string
	omit       []string
	joins      []string
	group      string
	having     any
	distinct   bool
	scopes     []func(*gorm.DB) *gorm.DB
	customFunc func(*gorm.DB) *gorm.DB
}

// BuildOptions 构建查询选项
func BuildOptions(opts []Option) *Options {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// ApplyOptions 应用查询选项到数据库连接
func ApplyOptions(db *gorm.DB, options *Options) *gorm.DB {
	if options == nil {
		return db
	}

	if options.filter != nil {
		db = db.Where(options.filter)
	}

	if options.limit > 0 {
		db = db.Limit(options.limit)
	}

	if options.offset > 0 {
		db = db.Offset(options.offset)
	}

	for _, clause := range options.clauses {
		db = db.Clauses(clause)
	}

	for _, preload := range options.preloads {
		db = db.Preload(preload)
	}

	if len(options.selects) > 0 {
		db = db.Select(options.selects)
	}

	if len(options.omit) > 0 {
		db = db.Omit(options.omit...)
	}

	for _, join := range options.joins {
		db = db.Joins(join)
	}

	if options.group != "" {
		db = db.Group(options.group)
	}

	if options.having != nil {
		db = db.Having(options.having)
	}

	if options.distinct {
		db = db.Distinct()
	}

	for _, scope := range options.scopes {
		db = scope(db)
	}

	if options.customFunc != nil {
		db = options.customFunc(db)
	}

	return db
}

// 基础查询选项

// WithFilter 设置过滤条件
func WithFilter(filter interface{}) Option {
	return func(o *Options) {
		o.filter = filter
	}
}

// WithLimit 设置查询限制
func WithLimit(limit int) Option {
	return func(o *Options) {
		o.limit = limit
	}
}

// WithOffset 设置查询偏移
func WithOffset(offset int) Option {
	return func(o *Options) {
		o.offset = offset
	}
}

// WithClauses 添加GORM子句
func WithClauses(clauses ...clause.Expression) Option {
	return func(o *Options) {
		o.clauses = append(o.clauses, clauses...)
	}
}

// 关系加载选项

// WithPreload 预加载关联
func WithPreload(preloads ...string) Option {
	return func(o *Options) {
		o.preloads = append(o.preloads, preloads...)
	}
}

// WithSelect 选择特定字段
func WithSelect(selects ...string) Option {
	return func(o *Options) {
		o.selects = selects
	}
}

// WithOmit 忽略特定字段
func WithOmit(omits ...string) Option {
	return func(o *Options) {
		o.omit = omits
	}
}

// WithJoin 添加JOIN语句
func WithJoin(joins ...string) Option {
	return func(o *Options) {
		o.joins = append(o.joins, joins...)
	}
}

// 分组与聚合选项

// WithGroup 设置分组字段
func WithGroup(group string) Option {
	return func(o *Options) {
		o.group = group
	}
}

// WithHaving 设置HAVING条件
func WithHaving(having interface{}) Option {
	return func(o *Options) {
		o.having = having
	}
}

// WithDistinct 设置DISTINCT查询
func WithDistinct(distinct bool) Option {
	return func(o *Options) {
		o.distinct = distinct
	}
}

// 高级选项

// WithScope 添加GORM作用域
func WithScope(scope func(*gorm.DB) *gorm.DB) Option {
	return func(o *Options) {
		o.scopes = append(o.scopes, scope)
	}
}

// WithCustom 自定义数据库函数
func WithCustom(fn func(*gorm.DB) *gorm.DB) Option {
	return func(o *Options) {
		o.customFunc = fn
	}
}

// 简写方法

// F WithFilter的简写
func F(filter interface{}) Option {
	return WithFilter(filter)
}

// L WithLimit的简写
func L(limit int) Option {
	return WithLimit(limit)
}

// O WithOffset的简写
func O(offset int) Option {
	return WithOffset(offset)
}

// C WithClauses的简写
func C(clauses ...clause.Expression) Option {
	return WithClauses(clauses...)
}

// P WithPreload的简写
func P(preloads ...string) Option {
	return WithPreload(preloads...)
}

// J WithJoin的简写
func J(joins ...string) Option {
	return WithJoin(joins...)
}

// D WithDistinct的简写
func D(distinct bool) Option {
	return WithDistinct(distinct)
}

// G WithGroup的简写
func G(group string) Option {
	return WithGroup(group)
}

// H WithHaving的简写
func H(having interface{}) Option {
	return WithHaving(having)
}

// S WithScope的简写
func S(scope func(*gorm.DB) *gorm.DB) Option {
	return WithScope(scope)
}

// CF WithCustom的简写
func CF(fn func(*gorm.DB) *gorm.DB) Option {
	return WithCustom(fn)
}

// 链式构建器

// Q 链式构建查询条件
func (o *Options) Q(opts ...Option) *Options {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Where 添加WHERE条件
func (o *Options) Where(condition interface{}) *Options {
	return o.Q(F(condition))
}

// Preload 添加预加载
func (o *Options) Preload(preloads ...string) *Options {
	return o.Q(P(preloads...))
}

// Join 添加JOIN语句
func (o *Options) Join(joins ...string) *Options {
	return o.Q(J(joins...))
}

// Limit 设置限制
func (o *Options) Limit(limit int) *Options {
	return o.Q(L(limit))
}

// Offset 设置偏移
func (o *Options) Offset(offset int) *Options {
	return o.Q(O(offset))
}

// Distinct 设置DISTINCT
func (o *Options) Distinct(distinct bool) *Options {
	return o.Q(D(distinct))
}

// Group 设置分组
func (o *Options) Group(group string) *Options {
	return o.Q(G(group))
}

// Having 设置HAVING条件
func (o *Options) Having(having interface{}) *Options {
	return o.Q(H(having))
}

// Custom 添加自定义函数
func (o *Options) Custom(fn func(*gorm.DB) *gorm.DB) *Options {
	return o.Q(CF(fn))
}

// Page 设置分页参数
func (o *Options) Page(page, pageSize int) *Options {
	return o.Q(L(pageSize), O((page-1)*pageSize))
}
