package data

import (
	"context"
	"event/internal/config"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ChClient struct {
	db driver.Conn
}

func NewClickHouseDB(ctx context.Context, c config.ClickHouse) (*ChClient, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", c.Host, c.Port)},
		Auth: clickhouse.Auth{
			Database: c.DBName,
			Username: c.Username,
			Password: c.Password,
		},
		Debug:           false,
		DialTimeout:     time.Second,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	})

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	return &ChClient{db: conn}, err
}
