package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/Richtermnd/TgLogin/internal/config"
	"github.com/Richtermnd/TgLogin/internal/domain"
	"github.com/Richtermnd/TgLogin/internal/storage"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sqlx.DB
}

func New() *Storage {
	path := config.Config().Storage.URL
	db := connectToDB(path)
	s := &Storage{db: db}
	s.init()
	return s
}

func (s *Storage) SaveUser(ctx context.Context, user domain.User) error {
	const op = "storage.sqlite.SaveUser"
	stmt, args, err := sq.
		Insert(user.Table()).
		Columns(user.Columns()...).
		Values(user.Values()...).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = s.db.ExecContext(ctx, stmt, args...)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s: %w", op, storage.ErrAlreadyExist)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) UserByTGID(ctx context.Context, TGID int64) (domain.User, error) {
	const op = "storage.sqlite.UserByTGID"
	var user domain.User

	stmt, args, err := sq.
		Select(user.Columns()...).
		From(user.Table()).
		Where(sq.Eq{"tg_id": TGID}).
		ToSql()
	if err != nil {
		return user, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.db.GetContext(ctx, &user, stmt, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) UserByUsername(ctx context.Context, username string) (domain.User, error) {
	const op = "storage.sqlite.UserByUsername"
	var user domain.User

	stmt, args, err := sq.
		Select(user.Columns()...).
		From(user.Table()).
		Where(sq.Eq{"username": username}).
		ToSql()
	if err != nil {
		return user, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.db.GetContext(ctx, &user, stmt, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) UpdateLastLogin(ctx context.Context, TGID, lastLogin int64) error {
	const op = "storage.sqlite.UpdateLastLogin"
	var user domain.User
	stmt, args, err := sq.
		Update(user.Table()).
		Set("last_login", lastLogin).
		Where("tg_id", TGID).
		ToSql()
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
	}
	return err
}

func (s *Storage) init() {
	s.db.MustExec(`
	CREATE TABLE IF NOT EXISTS users (
    	tg_id INTEGER PRIMARY KEY,
    	first_name TEXT NOT NULL,
		last_name TEXT DEFAULT NULL,
		username TEXT DEFAULT NULL,
		photo_url TEXT DEFAULT NULL,
		last_login INTEGER,
		registered INTERGER
	);
	`)
}

func connectToDB(path string) *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
