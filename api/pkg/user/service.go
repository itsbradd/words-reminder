package user

import (
	"context"
	"database/sql"
	"github.com/sonngocme/words-reminder-be/db"
)

type Storage interface {
	SignUpUser(ctx context.Context, arg db.SignUpUserParams) (int64, error)
	SetUserRefreshToken(ctx context.Context, arg db.SetUserRefreshTokenParams) error
}

type PassHasher interface {
	HashPassword(pass string) (string, error)
	VerifyPassword(hashedPassword, pass string) error
}

type service struct {
	storage    Storage
	passHasher PassHasher
}

var _ Service = (*service)(nil)

func NewService(storage Storage, passHasher PassHasher) Service {
	return &service{
		storage:    storage,
		passHasher: passHasher,
	}
}

func (s *service) SignUpUser(ctx context.Context, arg db.SignUpUserParams) (int64, error) {
	return s.storage.SignUpUser(ctx, arg)
}

func (s *service) HashPassword(pass string) (string, error) {
	res, err := s.passHasher.HashPassword(pass)
	return res, err
}

func (s *service) SetUserRefreshToken(ctx context.Context, id int64, token string) error {
	return s.storage.SetUserRefreshToken(ctx, db.SetUserRefreshTokenParams{
		RefreshToken: sql.NullString{
			Valid:  true,
			String: token,
		},
		ID: int32(id),
	})
}
