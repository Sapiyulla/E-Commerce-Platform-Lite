package repository

import (
	"context"
	"errors"
	"order-service/internal/domain"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	yugapool "github.com/yugabyte/pgx/v5/pgxpool"
)

type YugabyteDBRepository struct {
	ctx  time.Duration
	pool *yugapool.Pool
}

var ErrOrderAlreadyExists = errors.New("order already exists")

func NewYugabyteDBRepository(ctx time.Duration, pool *yugapool.Pool) *YugabyteDBRepository {
	return &YugabyteDBRepository{
		ctx:  ctx,
		pool: pool,
	}
}

func (r *YugabyteDBRepository) Add(order *domain.Order) error {
	reqCtx, cancel := context.WithTimeout(context.Background(), r.ctx)
	defer cancel()

	_, err := r.pool.Exec(
		reqCtx,
		`INSERT INTO orders (id, item_id, payed) VALUES ($1, $2, $3)`,
		order.ID, order.Item_id, order.Payed,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return ErrOrderAlreadyExists
		}
	}
	return err
}

func (r *YugabyteDBRepository) Pay(id string) error {
	reqCtx, cancel := context.WithTimeout(context.Background(), r.ctx)
	defer cancel()

	_, err := r.pool.Exec(
		reqCtx,
		`UPDATE orders SET payed = true WHERE payed IS DISTINCT FROM true AND id = $1`,
		id,
	)
	return err
}

func InitMigrate(pool *yugapool.Pool) {
	_, err := pool.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS orders (
			id UUID PRIMARY KEY,
			item_id UUID UNIQUE NOT NULL,
			payed BOOLEAN DEFAULT false
		);`,
	)
	if err != nil {
		logrus.Errorf("init migrate: %s", err.Error())
		return
	}
	logrus.Info("init migrate: succes")
}
