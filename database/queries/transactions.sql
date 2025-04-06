-- name: CreateTransaction :one
INSERT INTO transactions (
    user_id,
    amount,
    description,
    transaction_type,
    category,
    created_at,
    updated_at
) VALUES (
    sqlc.arg(user_id),
    sqlc.arg(amount),
    sqlc.arg(description),
    sqlc.arg(transaction_type),
    sqlc.arg(category),
    sqlc.arg(date),
    sqlc.arg(date)
) RETURNING *;
-- name: GetTransactionByID :one
SELECT * FROM transactions WHERE id = sqlc.arg(id);
-- name: GetTransactionsByUserID :many
SELECT * FROM transactions WHERE user_id = sqlc.arg(user_id);
-- name: GetTransactionsByUserIDAndTransactionType :many
SELECT * FROM transactions WHERE user_id = sqlc.arg(user_id) AND transaction_type = sqlc.arg(transaction_type);
-- name: GetLatestTransactionsByUserID :many
SELECT * FROM transactions WHERE user_id = sqlc.arg(user_id) ORDER BY updated_at DESC;
-- name: GetLatestTransactionsByUserIDAndTransactionType :many
SELECT * FROM transactions WHERE user_id = sqlc.arg(user_id) AND transaction_type = sqlc.arg(transaction_type) ORDER BY updated_at DESC;
-- name: GetTotalTransactionsByUserIDAndTransactionTypeForThisMonth :one
SELECT SUM(amount) FROM transactions WHERE user_id = sqlc.arg(user_id) AND transaction_type = sqlc.arg(transaction_type) AND strftime('%Y-%m', created_at) = strftime('%Y-%m', 'now') ORDER BY updated_at DESC;
-- name: GetTotalTransactionsByUserIDAndTransactionTypeForLastMonth :one
SELECT SUM(amount) FROM transactions WHERE user_id = sqlc.arg(user_id) AND transaction_type = sqlc.arg(transaction_type) AND strftime('%Y-%m', created_at) = strftime('%Y-%m', 'now', '-1 month') ORDER BY updated_at DESC;
-- name: GetTransactionsCountByUserIDAndTransactionType :one
SELECT COUNT(*) FROM transactions WHERE user_id = sqlc.arg(user_id) AND transaction_type = sqlc.arg(transaction_type);
-- name: GetTotalTransactionsThisYearByUserIDAndTransactionType :one
SELECT SUM(amount) FROM transactions WHERE user_id = sqlc.arg(user_id) AND transaction_type = sqlc.arg(transaction_type) AND strftime('%Y', created_at) = strftime('%Y', 'now') ORDER BY updated_at DESC;
-- name: DeleteTransactionByIDAndUserID :exec
DELETE FROM transactions WHERE id = sqlc.arg(id) AND user_id = sqlc.arg(user_id);
-- name: GetTotalBalanceByUserID :one
SELECT SUM(CASE WHEN transaction_type = 1 THEN amount ELSE -amount END) AS total_balance FROM transactions WHERE user_id = sqlc.arg(user_id);
-- name: GetTotalBalanceForLastMonthByUserID :one
SELECT SUM(CASE WHEN transaction_type = 1 THEN amount ELSE -amount END) AS total_balance FROM transactions WHERE user_id = sqlc.arg(user_id) AND strftime('%Y-%m', created_at) = strftime('%Y-%m', 'now', '-1 month');
-- name: GetCategorySummaryByUserID :many
SELECT 
    category,
    SUM(amount) as total_amount
FROM 
    transactions
WHERE 
    transaction_type = 2
    AND strftime('%m', updated_at) = strftime('%m', 'now')
    AND strftime('%Y', updated_at) = strftime('%Y', 'now')
    AND user_id = sqlc.arg(user_id)
GROUP BY 
    category
ORDER BY 
    total_amount DESC;
-- name: GetMonthlyOverviewByUserID :many
WITH months AS (
    SELECT 
        datetime('now', '-3 months', 'start of month') AS month_start
    UNION SELECT 
        datetime('now', '-2 months', 'start of month')
    UNION SELECT 
        datetime('now', '-1 month', 'start of month')
    UNION SELECT 
        datetime('now', 'start of month')
    UNION SELECT 
        datetime('now', '+1 month', 'start of month')
    UNION SELECT 
        datetime('now', '+2 months', 'start of month')
)

SELECT 
    strftime('%m', months.month_start) AS month_num,
    CASE 
        WHEN strftime('%m', months.month_start) = '01' THEN 'Jan'
        WHEN strftime('%m', months.month_start) = '02' THEN 'Feb'
        WHEN strftime('%m', months.month_start) = '03' THEN 'Mar'
        WHEN strftime('%m', months.month_start) = '04' THEN 'Apr'
        WHEN strftime('%m', months.month_start) = '05' THEN 'May'
        WHEN strftime('%m', months.month_start) = '06' THEN 'Jun'
        WHEN strftime('%m', months.month_start) = '07' THEN 'Jul'
        WHEN strftime('%m', months.month_start) = '08' THEN 'Aug'
        WHEN strftime('%m', months.month_start) = '09' THEN 'Sep'
        WHEN strftime('%m', months.month_start) = '10' THEN 'Oct'
        WHEN strftime('%m', months.month_start) = '11' THEN 'Nov'
        WHEN strftime('%m', months.month_start) = '12' THEN 'Dec'
    END AS month_name,
    COALESCE(SUM(CASE WHEN t.transaction_type = 1 THEN t.amount ELSE 0 END), 0) AS income,
    COALESCE(SUM(CASE WHEN t.transaction_type = 2 THEN t.amount ELSE 0 END), 0) AS expenses
FROM 
    months
LEFT JOIN 
    transactions t ON strftime('%Y-%m', t.updated_at) = strftime('%Y-%m', months.month_start)
WHERE 
    t.user_id = sqlc.arg(user_id)
GROUP BY 
    month_num, month_name
ORDER BY 
    strftime('%Y-%m', months.month_start);
