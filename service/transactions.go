package service

import (
	"context"
	"errors"
	"log/slog"
	"net/url"
	"strconv"
	"time"

	"github.com/MaxIvanyshen/budget-tracker/database"
	"github.com/MaxIvanyshen/budget-tracker/database/sqlc"
)

func (s *Service) createExpense(ctx context.Context, userID int64, form url.Values) (*sqlc.Transactions, error) {
	amountStr := form.Get("amount")
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	category := form.Get("category")
	description := form.Get("description")
	dateStr := form.Get("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New("Invalid date format")
	}

	s.logger.LogAttrs(ctx, slog.LevelInfo, "Creating expense", slog.Any("amount", amount), slog.Any("category", category), slog.Any("description", description), slog.Any("date", dateStr))

	if amountStr == "" || description == "" || dateStr == "" {
		return nil, errors.New("All fields are required")
	}

	expense, err := s.queries.CreateTransaction(ctx, &sqlc.CreateTransactionParams{
		UserID:          userID,
		Amount:          amount,
		Description:     description,
		TransactionType: int64(database.TransactionTypeExpense),
		Category:        &category,
		Date:            date,
	})
	if err != nil {
		s.logger.LogAttrs(ctx, slog.LevelError, "Failed to create transaction", slog.Any("error", err))
		return nil, errors.New("Failed to create transaction")
	}
	return expense, nil
}

func (s *Service) deleteTransaction(ctx context.Context, userID, transactionID int64) error {
	err := s.queries.DeleteTransactionByIDAndUserID(ctx, &sqlc.DeleteTransactionByIDAndUserIDParams{
		ID:     transactionID,
		UserID: userID,
	})
	if err != nil {
		s.logger.LogAttrs(ctx, slog.LevelError, "Failed to delete transaction", slog.Any("error", err))
		return errors.New("Failed to delete transaction")
	}
	return nil
}
