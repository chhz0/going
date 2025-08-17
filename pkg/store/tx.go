package store

import (
	"context"
	"fmt"

	"github.com/chhz0/going/pkg/store/query"
	"gorm.io/gorm"
)

// TxStore 事务存储
type TxStore[T any] struct {
	db     *gorm.DB
	logger Logger
}

func (s *TxStore[T]) Create(ctx context.Context, entity *T) error {
	if err := s.db.Create(entity).Error; err != nil {
		s.logError(ctx, "tx create failed", err)
		return err
	}
	return nil
}

func (s *TxStore[T]) Update(ctx context.Context, entity *T, opts ...query.Option) error {
	options := query.BuildOptions(opts)
	db := query.ApplyOptions(s.db, options)

	if err := db.Updates(entity).Error; err != nil {
		s.logError(ctx, "tx update failed", err)
		return err
	}
	return nil
}

func (s *TxStore[T]) Delete(ctx context.Context, opts ...query.Option) error {
	options := query.BuildOptions(opts)
	db := query.ApplyOptions(s.db, options)

	var model T
	if err := db.Delete(&model).Error; err != nil {
		s.logError(ctx, "tx delete failed", err)
		return err
	}
	return nil
}

func (s *TxStore[T]) Get(ctx context.Context, opts ...query.Option) (*T, error) {
	options := query.BuildOptions(opts)
	db := query.ApplyOptions(s.db, options)

	var entity T
	if err := db.First(&entity).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			s.logError(ctx, "tx get failed", err)
		}
		return nil, err
	}
	return &entity, nil
}

func (s *TxStore[T]) List(ctx context.Context, opts ...query.Option) ([]*T, error) {
	options := query.BuildOptions(opts)
	db := query.ApplyOptions(s.db, options)

	var entities []*T
	if err := db.Find(&entities).Error; err != nil {
		s.logError(ctx, "tx list failed", err)
		return nil, err
	}
	return entities, nil
}

func (s *TxStore[T]) Count(ctx context.Context, opts ...query.Option) (int64, error) {
	options := query.BuildOptions(opts)
	db := query.ApplyOptions(s.db, options)

	var count int64
	if err := db.Model(new(T)).Count(&count).Error; err != nil {
		s.logError(ctx, "tx count failed", err)
		return 0, err
	}
	return count, nil
}

func (s *TxStore[T]) logError(ctx context.Context, msg string, err error) {
	var model T
	s.logger.Error(ctx, fmt.Sprintf("%T %s", model, msg), "error", err)
}
