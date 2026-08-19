-- name: InsertCryptoPrices :exec
INSERT INTO crypto_prices (coin_symbol, coin_price)
VALUES ($1, $2);

-- name: GetLatestPrice :one
SELECT * FROM crypto_prices
WHERE coin_symbol = $1
ORDER BY created_at DESC
LIMIT 1;