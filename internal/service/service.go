package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Richtermnd/TgLogin/internal/config"
	"github.com/Richtermnd/TgLogin/internal/domain"
	"github.com/Richtermnd/TgLogin/internal/storage"
	"github.com/Richtermnd/TgLogin/pkg/tglogin"
)

var (
	ErrNotTelegramData = errors.New("not telegram data")
	ErrExpired         = errors.New("expired")
)

type Storage interface {
	SaveUser(ctx context.Context, user domain.User) error
	UserByTGID(ctx context.Context, TGID int64) (domain.User, error)
	UserByUsername(ctx context.Context, username string) (domain.User, error)
	UpdateLastLogin(ctx context.Context, TGID, lastLogin int64) error
}

type Service struct {
	log     *slog.Logger
	storage Storage
}

func New(log *slog.Logger, storage Storage) *Service {
	return &Service{
		log:     log,
		storage: storage,
	}
}

func (s *Service) User(ctx context.Context, userData tglogin.TelegramUserData) (domain.User, error) {
	const op = "service.User"
	log := s.log.With(slog.String("op", op), slog.Int64("TGID", userData.TGID))

	// Get user from storage
	log.Info("Get user from storage")
	user, err := s.storage.UserByTGID(ctx, userData.TGID)
	if err != nil {
		// Handle error
		if errors.Is(err, storage.ErrNotFound) {
			log.Info("User not found")
		} else {
			log.Error("Failed to get user from storage", slog.String("error", err.Error()))
		}
		return domain.User{}, err
	}
	log.Info("user found")
	return user, nil
}

func (s *Service) RegisterUser(ctx context.Context, userData tglogin.TelegramUserData) error {
	const op = "service.RegisterUser"
	log := s.log.With(slog.String("op", op), slog.Int64("TGID", userData.TGID))

	// Convert TelegramUserData to domaim.User
	user := telegramUserDataToUser(userData)

	// Registered at first login time
	user.Registered = user.LastLogin

	// Save user in storage
	err := s.storage.SaveUser(ctx, user)
	if err != nil {
		log.Error("Failed to save user", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("User registered")
	return nil
}

func (s *Service) Login(ctx context.Context, userData tglogin.TelegramUserData) error {
	const op = "service.Login"
	log := s.log.With(slog.String("op", op))
	user := telegramUserDataToUser(userData)
	err := s.storage.UpdateLastLogin(ctx, user.TGID, user.LastLogin)
	if err != nil {
		log.Warn("Failed to update last login", slog.String("error", err.Error()))
	}
	return err
}

func (s *Service) IsAuthentificated(ctx context.Context, userData tglogin.TelegramUserData) bool {
	const op = "service.IsAuthentificated"
	log := s.log.With(slog.String("op", op))
	log.Info("Check user")
	if !tglogin.IsTelegramAuthorization(userData, config.Config().Token) {
		log.Info("Not telegram data")
		return false
	}
	if !tglogin.IsExpiredData(userData.AuthDate, config.Config().LoginTTL) {
		log.Info("Expired telegram data")
		return false
	}
	return true
}

func telegramUserDataToUser(data tglogin.TelegramUserData) domain.User {
	var user domain.User
	user.TGID = data.TGID
	user.FirstName = data.FirstName
	user.LastName = data.LastName
	user.Username = data.Username
	user.PhotoURL = data.PhotoURL
	user.LastLogin = data.AuthDate
	return user
}
