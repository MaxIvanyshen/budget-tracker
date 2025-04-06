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
		"difference":     differenceInPercent,
		"expenses":       expenses,
	}, nil
}

func (s *Service) getIncomePageInfo(ctx context.Context, userID int64) (map[string]any, error) {
	totalThisMonth, err := s.queries.GetTotalTransactionsByUserIDAndTransactionTypeForThisMonth(ctx, &sqlc.GetTotalTransactionsByUserIDAndTransactionTypeForThisMonthParams{
		UserID:          userID,
		TransactionType: int64(database.TransactionTypeIncome),
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
		TransactionType: int64(database.TransactionTypeIncome),
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
		TransactionType: int64(database.TransactionTypeIncome),
	})
	if err != nil {
		s.logger.LogAttrs(ctx, slog.LevelError, "Failed to get transactions count", slog.Any("error", err))
		return nil, err
	}
	year, err := s.queries.GetTotalTransactionsThisYearByUserIDAndTransactionType(ctx, &sqlc.GetTotalTransactionsThisYearByUserIDAndTransactionTypeParams{
		UserID:          userID,
		TransactionType: int64(database.TransactionTypeIncome),
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

	incomes, err := s.queries.GetTransactionsByUserIDAndTransactionType(ctx, &sqlc.GetTransactionsByUserIDAndTransactionTypeParams{
		UserID:          userID,
		TransactionType: int64(database.TransactionTypeIncome),
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
		"difference":     differenceInPercent,
		"incomes":        incomes,
	}, nil
}

func (s *Service) getDashboardPageInfo(ctx context.Context, userID int64) (map[string]any, error) {
	transactions, err := s.queries.GetLatestTransactionsByUserID(ctx, userID)
	if err != nil {
		s.logger.LogAttrs(ctx, slog.LevelError, "Failed to get transactions", slog.Any("error", err))
		return nil, err
	}

	totalBalance, err := s.queries.GetTotalBalanceByUserID(ctx, userID)
	if err != nil {
		s.logger.LogAttrs(ctx, slog.LevelError, "Failed to get total balance", slog.Any("error", err))
		return nil, err
	}
	totalBalanceLastMonth, err := s.queries.GetTotalBalanceForLastMonthByUserID(ctx, userID)
	if err != nil {
		s.logger.LogAttrs(ctx, slog.LevelError, "Failed to get total balance for last month", slog.Any("error", err))
		return nil, err
	}

	if totalBalance == nil {
		totalBalance = new(float64)
	}
	if totalBalanceLastMonth == nil {
		totalBalanceLastMonth = new(float64)
	}
	balanceDiffInPercent := *totalBalance / *totalBalanceLastMonth * 100 - 100
	if totalBalanceLastMonth == nil || *totalBalanceLastMonth == 0 {
		balanceDiffInPercent = 100
	} else if totalBalance == nil || *totalBalance == 0 {
		balanceDiffInPercent = 0
	}

	totalIncome, err := s.queries.GetTotalTransactionsByUserIDAndTransactionTypeForThisMonth(ctx, &sqlc.GetTotalTransactionsByUserIDAndTransactionTypeForThisMonthParams{
		UserID:          userID,
		TransactionType: int64(database.TransactionTypeIncome),
	})
	if err != nil {
		s.logger.LogAttrs(ctx, slog.LevelError, "Failed to get total transactions", slog.Any("error", err))
		return nil, err
	}
	if totalIncome == nil {
		totalIncome = new(float64)
	}
	incomeLastMonth, err := s.queries.GetTotalTransactionsByUserIDAndTransactionTypeForLastMonth(ctx, &sqlc.GetTotalTransactionsByUserIDAndTransactionTypeForLastMonthParams{
		UserID:          userID,
		TransactionType: int64(database.TransactionTypeIncome),
	})
	if err != nil {
		s.logger.LogAttrs(ctx, slog.LevelError, "Failed to get total transactions", slog.Any("error", err))
		return nil, err
	}
	if incomeLastMonth == nil {
		incomeLastMonth = new(float64)
	}

	incomeDiff := *totalIncome / *incomeLastMonth * 100 - 100
	if incomeLastMonth == nil || *incomeLastMonth == 0 {
		incomeDiff = 100
	} else if totalIncome == nil || *totalIncome == 0 {
		incomeDiff = 0
	}

	totalExpense, err := s.queries.GetTotalTransactionsByUserIDAndTransactionTypeForThisMonth(ctx, &sqlc.GetTotalTransactionsByUserIDAndTransactionTypeForThisMonthParams{
		UserID:          userID,
		TransactionType: int64(database.TransactionTypeExpense),
	})
	if err != nil {
		s.logger.LogAttrs(ctx, slog.LevelError, "Failed to get total transactions", slog.Any("error", err))
		return nil, err
	}
	if totalExpense == nil {
		totalExpense = new(float64)
	}

	expenseLastMonth, err := s.queries.GetTotalTransactionsByUserIDAndTransactionTypeForLastMonth(ctx, &sqlc.GetTotalTransactionsByUserIDAndTransactionTypeForLastMonthParams{
		UserID:          userID,
		TransactionType: int64(database.TransactionTypeExpense),
	})
	if err != nil {
		s.logger.LogAttrs(ctx, slog.LevelError, "Failed to get total transactions", slog.Any("error", err))
		return nil, err
	}
	if expenseLastMonth == nil {
		expenseLastMonth = new(float64)
	}
	expenseDiff := *totalExpense / *expenseLastMonth * 100 - 100
	if expenseLastMonth == nil || *expenseLastMonth == 0 {
		expenseDiff = 100
	} else if totalExpense == nil || *totalExpense == 0 {
		expenseDiff = 0
	}

	categorySummary, err := s.queries.GetCategorySummaryByUserID(ctx, userID)
	if err != nil {
		s.logger.LogAttrs(ctx, slog.LevelError, "Failed to get category summary", slog.Any("error", err))
		return nil, err
	}

	categories := make([]string, len(categorySummary))
	categoryData := make([]float64, len(categorySummary))
	for i, category := range categorySummary {
		categories[i] = *category.Category
		categoryData[i] = *category.TotalAmount
	}

	monthlyOverview, err := s.queries.GetMonthlyOverviewByUserID(ctx, userID)
	if err != nil {
		s.logger.LogAttrs(ctx, slog.LevelError, "Failed to get monthly overview", slog.Any("error", err))
		return nil, err
	}

	months := make([]string, len(monthlyOverview))
	monthlyDataIncomes := make([]int64, len(monthlyOverview))
	monthlyDataExpenses := make([]int64, len(monthlyOverview))
	for i, month := range monthlyOverview {
		months[i] = month.MonthName.(string)
		monthlyDataIncomes[i] = month.Income.(int64)
		monthlyDataExpenses[i] = month.Expenses.(int64)
	}

	return map[string]any{
		"transactions":        transactions,
		"totalBalance":        *totalBalance,
		"balanceDiff":         balanceDiffInPercent,
		"totalIncome":         *totalIncome,
		"incomeDiff":          incomeDiff,
		"totalExpense":        *totalExpense,
		"expenseDiff":         expenseDiff,
		"categories":          categories,
		"categoryData":        categoryData,
		"monthlyDataIncomes":  monthlyDataIncomes,
		"monthlyDataExpenses": monthlyDataExpenses,
		"monthlyDataMonths":   months,
	}, nil
}
