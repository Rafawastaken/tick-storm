-- name: InsertPrice :exec
INSERT INTO crypto_prices (exchange, coin_symbol, coin_price, trade_id)
VALUES ($1, $2, $3, $4)
ON CONFLICT DO NOTHING;

-- name: GetLatestPrice :one
SELECT * FROM crypto_prices
WHERE exchange = $1 AND coin_symbol = $2
ORDER BY created_at DESC, id DESC
LIMIT 1;

-- name: ListPricesForCoin :many
SELECT * FROM crypto_prices
WHERE exchange = $1 AND coin_symbol = $2
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg(max_results);