package postgres

import (
	"context"

	"gorm.io/gorm"
	"github.com/chuminh2001100/goPractice/domain"
	"github.com/chuminh2001100/goPractice/service"
)

type postgresStorgae struct {
	db *gorm.DB
}

func NewPostgresStorage(db *gorm.DB) *postgresStorgae {
	return &postgresStorgae{db: db}
}

func (s *postgresStorgae) CreateUser(ctx context.Context, u user.User) error {
	if err := s.db.Create(u).Error; err != nil {
		return err
	}

	return nil
}
