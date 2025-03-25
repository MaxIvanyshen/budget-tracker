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
