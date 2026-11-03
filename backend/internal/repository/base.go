package repository

import "gorm.io/gorm"

type BaseRepository struct {
	db *gorm.DB
}

func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{
		db: db,
	}
}

func (r *BaseRepository) WithTx(tx *gorm.DB) *BaseRepository {
	return &BaseRepository{
		db: tx,
	}
}

func (r *BaseRepository) DB() *gorm.DB {
	return r.db
}
