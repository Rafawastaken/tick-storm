-- name: InsertPrice :exec
INSERT INTO crypto_prices (exchange, coin_symbol, coin_price, trade_id, traded_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT DO NOTHING;

-- name: GetLatestPrice :one
SELECT * FROM crypto_prices
WHERE exchange = $1 AND coin_symbol = $2
ORDER BY traded_at DESC, id DESC
LIMIT 1;

-- name: ListPricesForCoin :many
SELECT * FROM crypto_prices
WHERE exchange = $1
  AND coin_symbol = $2
  AND (traded_at, id) < (sqlc.arg(before_time)::timestamptz, sqlc.arg(before_id)::bigint)
ORDER BY traded_at DESC, id DESC
LIMIT sqlc.arg(max_results);
