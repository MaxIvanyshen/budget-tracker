package service

import (
	"context"
	"log/slog"

	"github.com/MaxIvanyshen/budget-tracker/database"
	"github.com/MaxIvanyshen/budget-tracker/database/sqlc"
)

func (s *Service) getExpensesPageInfo(ctx context.Context, userID int64) (map[string]any, error) {
	totalThisMonth, err := s.queries.GetTotalTransactionsByUserIDAndTransactionTypeForThisMonth(ctx, &sqlc.GetTotalTransactionsByUserIDAndTransactionTypeForThisMonthParams{
		UserID:          userID,
		TransactionType: int64(database.TransactionTypeExpense),
	})
	if err != nil {
		s.logger.LogAttrs(ctx, slog.LevelError, "Failed to get total transactions", slog.Any("error", err))
		return nil, err
	}
	if totalThisMonth == nil {
		totalThisMonth = new(float64)
	}
	totalLastMonth, err := s.queries.GetTotalTransactionsByUserIDAndTransactionTypeForLastMonth(ctx, &sqlc.GetTotalTransactionsByUserIDAndTransactionTypeForLastMonthParams{
		UserID:          userID,
		TransactionType: int64(database.TransactionTypeExpense),
	})
	if err != nil {
		s.logger.LogAttrs(ctx, slog.LevelError, "Failed to get total transactions", slog.Any("error", err))
		return nil, err
	}
	if totalLastMonth == nil {
		totalLastMonth = new(float64)
	}
	count, err := s.queries.GetTransactionsCountByUserIDAndTransactionType(ctx, &sqlc.GetTransactionsCountByUserIDAndTransactionTypeParams{
		UserID:          userID,
		TransactionType: int64(database.TransactionTypeExpense),
	})
	if err != nil {
		s.logger.LogAttrs(ctx, slog.LevelError, "Failed to get transactions count", slog.Any("error", err))
		return nil, err
	}
	year, err := s.queries.GetTotalTransactionsThisYearByUserIDAndTransactionType(ctx, &sqlc.GetTotalTransactionsThisYearByUserIDAndTransactionTypeParams{
		UserID:          userID,
		TransactionType: int64(database.TransactionTypeExpense),
	})
	if err != nil {
		s.logger.LogAttrs(ctx, slog.LevelError, "Failed to get average transactions", slog.Any("error", err))
		return nil, err
	}
	if year == nil {
		year = new(float64)
	}

	differenceInPercent := *totalThisMonth / *totalLastMonth * 100 - 100
	if totalLastMonth == nil || *totalLastMonth == 0 {
		differenceInPercent = 100
		if totalThisMonth == nil || *totalThisMonth == 0 {
			differenceInPercent = 0
		}
	}

	expenses, err := s.queries.GetTransactionsByUserIDAndTransactionType(ctx, &sqlc.GetTransactionsByUserIDAndTransactionTypeParams{
		UserID:          userID,
		TransactionType: int64(database.TransactionTypeExpense),
	})
	if err != nil {
		s.logger.LogAttrs(ctx, slog.LevelError, "Failed to get transactions", slog.Any("error", err))
		return nil, err
	}

	return map[string]any{
		"totalThisMonth": *totalThisMonth,
		"totalLastMonth": *totalLastMonth,
		"count":          count,
		"year":           *year,
		"difference":     int64(differenceInPercent),
		"expenses":       expenses,
	}, nil
}
