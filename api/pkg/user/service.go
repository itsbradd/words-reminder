package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/itsbradd/words-reminder-be/db"
	"github.com/itsbradd/words-reminder-be/pkg"
	"github.com/itsbradd/words-reminder-be/pkg/jwt"
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

func (s *service) SignUp(ctx context.Context, info SignUpInfo) (*Credentials, error) {
	if err := s.IsUsernameExists(ctx, info.Username); err != nil {
		if errors.Is(err, ErrUsernameExists) {
			return nil, &pkg.FailResponse[any]{
				StatusCode: fiber.StatusBadRequest,
				ErrorCode:  fiber.StatusBadRequest,
				Message:    "username already exists",
				Errors: struct {
					Username string `json:"username"`
				}{
					Username: "username already exists",
				},
			}
		}
		return nil, fiber.ErrInternalServerError
	}

	hashedPass, err := s.passHasher.HashPassword(info.Password)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	// Can return err duplicate username, above check doesn't in a transaction
	userId, err := s.CreateUser(ctx, db.CreateUserParams{
		Username: info.Username,
		Password: hashedPass,
	})
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	refresh, access, err := s.GenRefreshAndAccessToken(ctx, userId)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return &Credentials{
		RefreshToken: refresh,
		AccessToken:  access,
	}, nil
}

func (s *service) Login(ctx context.Context, info LoginInfo) (*Credentials, error) {
	user, err := s.GetUserByUsername(ctx, info.Username)
	if err != nil {
		return nil, &pkg.FailResponse[any]{
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

	err = s.passHasher.VerifyPassword(user.Password, info.Password)
	if err != nil {
		return nil, &pkg.FailResponse[any]{
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
