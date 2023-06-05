package data

import (
	"context"
	"fmt"
	"time"

	"github.com/MrSwartz/event/internal/config"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ChClient struct {
	db        driver.Conn
	tableName string
	dbName    string
}

func NewClickHouseDB(ctx context.Context, c config.Clickhouse) (*ChClient, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", c.Host, c.Port)},
		Auth: clickhouse.Auth{
			Database: c.DBName,
			Username: c.Username,
			Password: c.Password,
		},
		Debug:           c.Debug,
		DialTimeout:     time.Duration(c.DialTimeout) * time.Second,
		MaxOpenConns:    c.MaxOpenConns,
		MaxIdleConns:    c.MaxIdleConns,
		ConnMaxLifetime: time.Duration(c.ConnMaxLifetime) * time.Second,
	})

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	// не лучшая идея, но так проще сделать первый запуск сервиса когда ещё не создана таблица
	if err := conn.Exec(ctx, fmt.Sprintf(createTemplate, c.DBName, c.TableName)); err != nil {
		return nil, err
	}

	return &ChClient{
		db:        conn,
		tableName: c.TableName,
		dbName:    c.DBName,
	}, nil
}
