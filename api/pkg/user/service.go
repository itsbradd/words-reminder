package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sonngocme/words-reminder-be/db"
	"github.com/sonngocme/words-reminder-be/pkg"
	"github.com/sonngocme/words-reminder-be/pkg/jwt"
	"time"
)

var (
	ErrGenRefreshToken = errors.New("errored when generating refresh token")
	ErrGenAccessToken  = errors.New("errored when generating access token")
	ErrUsernameExists  = errors.New("username already exists")
)

type Storage interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (int64, error)
	SetUserRefreshToken(ctx context.Context, arg db.SetUserRefreshTokenParams) error
	GetUserByUsername(ctx context.Context, username string) (db.User, error)
}

type PassHasher interface {
	HashPassword(pass string) (string, error)
	VerifyPassword(hashedPassword, pass string) error
}

type service struct {
	storage    Storage
	passHasher PassHasher
	jwt        JWTService
}

var _ Service = (*service)(nil)

func NewService(storage Storage, passHasher PassHasher, jwt JWTService) Service {
	return &service{
		storage:    storage,
		passHasher: passHasher,
		jwt:        jwt,
	}
}

func (s *service) CreateUser(ctx context.Context, arg db.CreateUserParams) (int64, error) {
	return s.storage.CreateUser(ctx, arg)
}

func (s *service) HashPassword(pass string) (string, error) {
	res, err := s.passHasher.HashPassword(pass)
	return res, err
}

func (s *service) VerifyPassword(hashedPass, pass string) error {
	return s.passHasher.VerifyPassword(hashedPass, pass)
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

func (s *service) GetUserByUsername(ctx context.Context, username string) (db.User, error) {
	return s.storage.GetUserByUsername(ctx, username)
}

func (s *service) GenRefreshAndAccessToken(ctx context.Context, id int64) (string, string, error) {
	token, err := s.jwt.NewWithClaims(jwt.MapClaims{},
		s.jwt.GenIssuerClaim("Words Reminder"),
		s.jwt.GenSubjectClaim(id),
		s.jwt.GenAudienceClaim("Words Reminder"),
		s.jwt.GenIssueAtClaim(time.Now()),
	)
	if err != nil {
		return "", "", ErrGenRefreshToken
	}

	err = s.SetUserRefreshToken(ctx, id, token)
	if err != nil {
		return "", "", ErrGenRefreshToken
	}

	accessToken, err := s.jwt.NewWithClaims(jwt.MapClaims{},
		s.jwt.GenIssuerClaim("Words Reminder"),
		s.jwt.GenSubjectClaim(id),
		s.jwt.GenAudienceClaim("Words Reminder"),
		s.jwt.GenIssueAtClaim(time.Now()),
		s.jwt.GenExpireTimeClaim(time.Now().Add(1*time.Minute)),
	)
	if err != nil {
		return "", "", ErrGenAccessToken
	}
	return token, accessToken, nil
}

func (s *service) IsUsernameExists(ctx context.Context, username string) error {
	_, err := s.GetUserByUsername(ctx, username)

	// Err is nil when return data has value (user exists)
	if err == nil {
		return ErrUsernameExists
	}

	// Errored when querying to the database
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fiber.ErrInternalServerError
	}
	// Username has not been taken
	return nil
}

func (s *service) Login(ctx context.Context, info LoginInfo) (*Credentials, error) {
	user, err := s.GetUserByUsername(ctx, info.Username)
	if err != nil {
		return nil, &pkg.FailRes{
			StatusCode: fiber.StatusBadRequest,
			ErrorCode:  fiber.StatusBadRequest,
			Message:    "user not found",
			Errors: struct {
				Username string `json:"username"`
			}{
				Username: "username is not valid",
			},
		}
	}

	err = s.VerifyPassword(user.Password, info.Password)
	if err != nil {
		return nil, &pkg.FailRes{
			StatusCode: fiber.StatusBadRequest,
			ErrorCode:  fiber.StatusBadRequest,
			Message:    "password is not match",
			Errors: struct {
				Password string `json:"password"`
			}{
				Password: "password is not match",
			},
		}
	}

	refresh, access, err := s.GenRefreshAndAccessToken(ctx, int64(user.ID))
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return &Credentials{
		RefreshToken: refresh,
		AccessToken:  access,
	}, nil
}
