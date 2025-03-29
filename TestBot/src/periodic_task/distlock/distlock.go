package distlock

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"

	"context"
	"time"
)

func AcquireLock(conn *pgxpool.Conn, lockName string, ownerID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	var owner string
	err = tx.QueryRow(ctx, "SELECT owner_id FROM distlocks WHERE lock_name = $1 FOR UPDATE", lockName).Scan(&owner)
	if err != nil {
		if err == pgx.ErrNoRows {
			_, err = tx.Exec(ctx, "INSERT INTO distlocks (lock_name, owner_id) VALUES ($1, $2)", lockName, ownerID)
			if err != nil {
				return false, err
			}

			err = tx.Commit(ctx)
			return err == nil, err
		}
		return false, err
	}

	if owner != ownerID {
		return false, nil
	}

	err = tx.Commit(ctx)
	return err == nil, err
}

func ReleaseLock(conn *pgxpool.Conn, lockName string, ownerID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer conn.Release()

	_, err := conn.Exec(ctx, "DELETE FROM distlocks WHERE lock_name = $1 AND owner_id = $2", lockName, ownerID)
	return err
}
