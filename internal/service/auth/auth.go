package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/tgkzz/auth/internal/domain/models"
	"github.com/tgkzz/auth/internal/storage/postgresql"
	"github.com/tgkzz/auth/pkg/jwt"
	"github.com/tgkzz/auth/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserStorage interface {
	CreateNewUser(ctx context.Context, user models.User) (userId int64, err error)
	GetUserByUsername(ctx context.Context, username string) (user *models.User, err error)
}

type Auth struct {
	log        *slog.Logger
	usrStorage UserStorage
}

func New(logger *slog.Logger, userStorage UserStorage) *Auth {
	return &Auth{
		log:        logger,
		usrStorage: userStorage,
	}
}

func (a *Auth) Register(ctx context.Context, username string, password string) (userId int64, err error) {
	const op = "auth.Register"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", logger.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	uid, err := a.usrStorage.CreateNewUser(ctx, models.User{Username: username, PassHash: passHash, Role: 1})
	if err != nil {
		log.Error("failed to create new user", logger.Err(err))

		if errors.Is(err, postgresql.ErrUserExists) {
			return 0, ErrInvalidCredentials
		}

		return 0, err
	}

	return uid, nil
}

func (a *Auth) Login(ctx context.Context, username string, password string) (token string, err error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	user, err := a.usrStorage.GetUserByUsername(ctx, username)
	if err != nil {
		log.Error("failed to get user", logger.Err(err))

		if errors.Is(err, postgresql.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", logger.Err(err))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	token, err = jwt.NewToken(*user)
	if err != nil {
		a.log.Info("failed to generate token", logger.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return
}
