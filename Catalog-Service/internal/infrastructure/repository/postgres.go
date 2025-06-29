package repository

import (
	dto "catalog-service/internal/application/DTO"
	"catalog-service/internal/domain"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func transaction(ctx context.Context, pool *pgxpool.Pool, fn func(tx pgx.Tx) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := recover(); err != nil {
			rollbackCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			tx.Rollback(rollbackCtx)
			panic(err)
		}
	}()

	if err := fn(tx); err != nil {
		rollbackCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		logrus.Errorf("Transaction error: %s", err.Error())
		_ = tx.Rollback(rollbackCtx)
		return err
	}

	commitCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return tx.Commit(commitCtx)
}

func (pg *PostgresRepository) AddItem(item *domain.Item) *dto.Status {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := transaction(ctx, pg.pool, func(tx pgx.Tx) error {
		_, err := tx.Exec(context.Background(), `INSERT INTO items VALUES ($1, $2, $3, $4)`, item.ID, item.Name, item.Category, item.Price)
		return err
	})
	if err != nil {
		return &dto.Status{Code: 2, Status: dto.Error, Message: err.Error()}
	}
	return &dto.Status{Code: 0, Status: dto.Success, Message: "Item added to catalog"}
}

func (pg *PostgresRepository) GetItems(category ...string) ([]*domain.Item, *dto.Status) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var items []*domain.Item
	err := transaction(ctx, pg.pool, func(tx pgx.Tx) error {
		var rows pgx.Rows
		var err error

		if len(category) == 0 {
			rows, err = tx.Query(ctx, `SELECT id, name, category, price FROM items`)
		} else {
			rows, err = tx.Query(ctx, `SELECT id, name, category, price FROM items WHERE category = $1`, category[0])
		}

		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var item domain.Item
			if err := rows.Scan(&item.ID, &item.Name, &item.Category, &item.Price); err != nil {
				return err
			}
			items = append(items, &item)
		}
		return rows.Err()
	})

	if err != nil {
		return nil, &dto.Status{Code: 2, Status: dto.Error, Message: err.Error()}
	}
	if len(items) == 0 {
		return nil, &dto.Status{Code: 0, Status: dto.NotFound, Message: "items not found"}
	}
	return items, &dto.Status{Code: 0, Status: dto.Success, Message: "items found"}
}

func (pg *PostgresRepository) GetItem(id string) (*domain.Item, *dto.Status) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var item domain.Item
	err := transaction(ctx, pg.pool, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx,
			`SELECT id, name, category, price FROM items WHERE id = $1`, id,
		).Scan(&item.ID, &item.Name, &item.Category, &item.Price)
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &dto.Status{Code: 0, Status: dto.NotFound, Message: "item not found"}
		}
		return nil, &dto.Status{Code: 2, Status: dto.Error, Message: err.Error()}
	}
	return &item, &dto.Status{Code: 0, Status: dto.Success, Message: "item found"}
}

func (pg *PostgresRepository) DeleteItem(id string) *dto.Status {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := transaction(ctx, pg.pool, func(tx pgx.Tx) error {
		_, err := tx.Exec(context.Background(), `DELETE FROM items WHERE id=$1`, id)
		return err
	})
	if err != nil {
		return &dto.Status{Code: 2, Status: dto.Error, Message: err.Error()}
	} else {
		return &dto.Status{Code: 0, Status: dto.Success, Message: "item with id=" + id + " deleted"}
	}
}
