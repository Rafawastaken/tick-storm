package pgxkit

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TryAdvisoryLock attempts to acquire a session-level postgres advisory lock.
//
//   - On success: returns (release, true, nil). Caller MUST defer release().
//   - When another session holds the lock: returns (nil, false, nil).
//   - On infra error: returns (nil, false, err).
//
// The lock is bound to the underlying connection. release() unlocks AND
// returns the conn to the pool — the two are inseparable. Skipping release()
// strands the conn (and the lock with it) until pg eventually closes the
// session.
func TryAdvisoryLock(ctx context.Context, pool *pgxpool.Pool, key int64) (release func(), ok bool, err error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, false, err
	}

	var locked bool
	if err := conn.QueryRow(ctx, "SELECT pg_try_advisory_lock($1)", key).Scan(&locked); err != nil {
		conn.Release()
		return nil, false, err
	}
	if !locked {
		conn.Release()
		return nil, false, nil
	}

	return func() {
		// Best-effort unlock — if the conn is dead, pg releases the lock
		// automatically when the session ends. Use Background ctx so a
		// cancelled request ctx doesn't prevent cleanup.
		_, _ = conn.Exec(context.Background(), "SELECT pg_advisory_unlock($1)", key)
		conn.Release()
	}, true, nil
}
