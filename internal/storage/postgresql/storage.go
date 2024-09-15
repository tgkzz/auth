package postgresql

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tgkzz/auth/internal/domain/models"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

type Storage struct {
	dbConn *pgxpool.Pool
}

func New(ctx context.Context, databaseUrl string) (*Storage, error) {
	conn, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(ctx); err != nil {
		return nil, err
	}

	return &Storage{dbConn: conn}, nil
}

func (s *Storage) CreateNewUser(ctx context.Context, user models.User) (userId int64, err error) {
	const op = "storage.postgresql.CreateNewUser"

	query, args, err := sq.Insert("users").
		Columns("username", "pass_hash", "role").
		Values(user.Username, user.PassHash, user.Role).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if err = s.dbConn.QueryRow(ctx, query, args...).Scan(&userId); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) { // username unique constraint error
			if pgErr.Code == "23505" {
				return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
			}
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return
}

func (s *Storage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	const op = "storage.postgresql.GetUserByUsername"

	query, args, err := sq.Select("*").
		From("users").
		Where(sq.Eq{"username": username}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var res models.User
	if err = s.dbConn.QueryRow(ctx, query, args...).Scan(&res.ID, &res.Username, &res.PassHash, &res.Role); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &res, nil
}
