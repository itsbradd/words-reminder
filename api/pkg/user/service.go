package user

import (
	"context"
	"database/sql"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/sonngocme/words-reminder-be/db"
	"github.com/sonngocme/words-reminder-be/pkg"
	"github.com/sonngocme/words-reminder-be/pkg/jwt"
	"time"
)

var (
	ErrGenRefreshToken = errors.New("errored when generating refresh token")
	ErrGenAccessToken  = errors.New("errored when generating access token")
)

type Storage interface {
	SignUpUser(ctx context.Context, arg db.SignUpUserParams) (int64, error)
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

func (s *service) SignUpUser(ctx context.Context, arg db.SignUpUserParams) (int64, error) {
	return s.storage.SignUpUser(ctx, arg)
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

func (s *service) Login(ctx context.Context, info LoginInfo) (*Credentials, error) {
	err := info.Validate()
	if err != nil {
		return nil, &pkg.FailRes{
			StatusCode: fiber.StatusBadRequest,
			ErrorCode:  fiber.StatusBadRequest,
			Message:    "Invalid parameters",
			Errors:     err.(validation.Errors),
		}
	}

	user, err := s.GetUserByUsername(ctx, info.Username)
	if err != nil {
		return nil, &pkg.FailRes{
			StatusCode: fiber.StatusBadRequest,
			ErrorCode:  fiber.StatusBadRequest,
			Message:    "User not found",
			Errors: struct {
				Username string `json:"username"`
			}{
				Username: "Username is not valid",
			},
		}
	}

	err = s.VerifyPassword(user.Password, info.Password)
	if err != nil {
		return nil, &pkg.FailRes{
			StatusCode: fiber.StatusBadRequest,
			ErrorCode:  fiber.StatusBadRequest,
			Message:    "Password is not match",
			Errors: struct {
				Password string `json:"password"`
			}{
				Password: "Password is not match",
			},
		}
	}

	refresh, access, err := s.GenRefreshAndAccessToken(ctx, int64(user.ID))
	if err != nil {
		return nil, &pkg.FailRes{
			StatusCode: fiber.StatusInternalServerError,
			ErrorCode:  fiber.StatusInternalServerError,
			Message:    "Something went wrong",
			Errors:     nil,
		}
	}

	return &Credentials{
		RefreshToken: refresh,
		AccessToken:  access,
	}, nil
}
