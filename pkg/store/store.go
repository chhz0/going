package store

import (
	"context"
	"fmt"

	"github.com/chhz0/going/pkg/store/query"
	"gorm.io/gorm"
)

// DBProvider 定义数据库连接获取接口
type DBProvider interface {
	DB(ctx context.Context) *gorm.DB
}

// Store 泛型数据存储
type Store[T any] struct {
	dbProvider DBProvider
	logger     Logger
}

// NewStore 创建新的数据存储实例
func NewStore[T any](provider DBProvider) *Store[T] {
	return &Store[T]{
		dbProvider: provider,
		logger:     emptyLogger{},
	}
}

// WithLogger 设置日志记录器
func (s *Store[T]) WithLogger(logger Logger) *Store[T] {
	s.logger = logger
	return s
}

// Create 创建新记录
func (s *Store[T]) Create(ctx context.Context, entity *T) error {
	db := s.dbProvider.DB(ctx)
	if err := db.Create(entity).Error; err != nil {
		s.logError(ctx, "create failed", err)
		return err
	}
	return nil
}

// Update 更新记录
func (s *Store[T]) Update(ctx context.Context, entity *T, opts ...query.Option) error {
	db := s.dbProvider.DB(ctx)
	options := query.BuildOptions(opts)
	db = query.ApplyOptions(db, options)

	if err := db.Updates(entity).Error; err != nil {
		s.logError(ctx, "update failed", err)
		return err
	}
	return nil
}

// Delete 删除记录
func (s *Store[T]) Delete(ctx context.Context, opts ...query.Option) error {
	db := s.dbProvider.DB(ctx)
	options := query.BuildOptions(opts)
	db = query.ApplyOptions(db, options)

	var model T
	if err := db.Delete(&model).Error; err != nil {
		s.logError(ctx, "delete failed", err)
		return err
	}
	return nil
}

// Get 查询单条记录
func (s *Store[T]) Get(ctx context.Context, opts ...query.Option) (*T, error) {
	db := s.dbProvider.DB(ctx)
	options := query.BuildOptions(opts)
	db = query.ApplyOptions(db, options)

	var entity T
	if err := db.First(&entity).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			s.logError(ctx, "get failed", err)
		}
		return nil, err
	}
	return &entity, nil
}

// List 查询多条记录
func (s *Store[T]) List(ctx context.Context, opts ...query.Option) ([]*T, error) {
	db := s.dbProvider.DB(ctx)
	options := query.BuildOptions(opts)
	db = query.ApplyOptions(db, options)

	var entities []*T
	if err := db.Find(&entities).Error; err != nil {
		s.logError(ctx, "list failed", err)
		return nil, err
	}
	return entities, nil
}

// Count 统计记录数
func (s *Store[T]) Count(ctx context.Context, opts ...query.Option) (int64, error) {
	db := s.dbProvider.DB(ctx)
	options := query.BuildOptions(opts)
	db = query.ApplyOptions(db, options)

	var count int64
	if err := db.Model(new(T)).Count(&count).Error; err != nil {
		s.logError(ctx, "count failed", err)
		return 0, err
	}
	return count, nil
}

// Transaction 执行数据库事务
func (s *Store[T]) Transaction(ctx context.Context, fn func(tx *TxStore[T]) error) error {
	db := s.dbProvider.DB(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		txStore := &TxStore[T]{
			db:     tx,
			logger: s.logger,
		}
		return fn(txStore)
	})
}

func (s *Store[T]) logError(ctx context.Context, msg string, err error) {
	var model T
	s.logger.Error(ctx, fmt.Sprintf("%T %s", model, msg), "error", err)
}
