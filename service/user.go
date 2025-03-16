package service

import (
	"context"
	"errors"
	"log/slog"
	"net/url"

	"github.com/MaxIvanyshen/budget-tracker/database/sqlc"
	"github.com/MaxIvanyshen/budget-tracker/utils"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) registerUser(ctx context.Context, form url.Values) (*sqlc.Users, error) {
	name := form.Get("name")
	password := form.Get("password")
	confirm := form.Get("confirm-password")
	terms := form.Get("terms-and-privacy")
	email := form.Get("email")

	if name == "" || password == "" || confirm == "" || email == "" {
		return nil, errors.New("All fields are required")
	}

	if password != confirm {
		return nil, errors.New("Passwords do not match")
	}

	err := utils.ValidatePassword(password)
	if err != nil {
		return nil, err
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("Failed to hash password")
	}

	password = string(encryptedPassword)

	if terms != "on" {
		return nil, errors.New("You must agree to terms and privacy policy")
	}

	_, err = s.queries.GetUserByEmail(ctx, email)
	if err == nil {
		return nil, errors.New("User with this email already exists")
	}

	user, err := s.queries.CreateUser(ctx, &sqlc.CreateUserParams{
		Name:     name,
		Password: password,
		Email:    email,
	})
	if err != nil {
		s.logger.LogAttrs(ctx, slog.LevelError, "Failed to create user", slog.Any("error", err))
		return nil, errors.New("Failed to create user")
	}
	return user, nil
}

func (s *Service) loginUser(ctx context.Context, form url.Values) (*sqlc.Users, error) {
	email := form.Get("email")
	password := form.Get("password")

	if email == "" || password == "" {
		return nil, errors.New("All fields are required")
	}

	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("User not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("Invalid email or password")
	}

	return user, nil
}
