CREATE TABLE crypto_prices (
    id          BIGSERIAL   PRIMARY KEY,
    exchange    VARCHAR(20) NOT NULL,
    coin_symbol VARCHAR(20) NOT NULL,
    coin_price  NUMERIC(18, 8) NOT NULL,
    trade_id    BIGINT,
    -- traded_at is market time (from the exchange); created_at is ingestion
    -- time. Their difference is the ingestion lag.
    traded_at   TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Mirrors the ORDER BY of the read queries: equality columns first,
-- then the sort columns in the same direction.
CREATE INDEX idx_crypto_prices_lookup
    ON crypto_prices (exchange, coin_symbol, traded_at DESC, id DESC);

-- Idempotency: reconnecting to a stream replays trades. Partial, because
-- synthetic or backfilled rows have no upstream id and must not collide.
CREATE UNIQUE INDEX uq_crypto_prices_trade
    ON crypto_prices (exchange, coin_symbol, trade_id)
    WHERE trade_id IS NOT NULL;
