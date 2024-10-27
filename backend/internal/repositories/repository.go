package repositories

import (
	"gorm.io/gorm"
)

type Repository[T any] struct {
    db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) *Repository[T] {
    return &Repository[T]{db: db}
}

func (r *Repository[T]) FindAll() ([]T, error) {
    var entities []T
    if err := r.db.Find(&entities).Error; err != nil {
        return nil, err
    }
    return entities, nil
}

func (r *Repository[T]) FindOneByID(id uint) (*T, error) {
    var entity T
    if err := r.db.First(&entity, id).Error; err != nil {
        return nil, err
    }
    return &entity, nil
}

func (r *Repository[T]) Create(entity *T) error {
    return r.db.Create(entity).Error
}

func (r *Repository[T]) Update(entity *T) error {
    return r.db.Save(entity).Error
}

func (r *Repository[T]) Delete(id uint) error {
    return r.db.Delete(new(T), id).Error
}
