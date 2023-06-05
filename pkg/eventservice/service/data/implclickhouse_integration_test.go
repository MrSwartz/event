package data

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/MrSwartz/event/internal/config"

	"github.com/stretchr/testify/require"
)

func TestInsert(t *testing.T) {
	cnf := config.Clickhouse{
		Host:            os.Getenv("CLICKHOUSE_HOST"),     // "127.0.0.1"
		Port:            os.Getenv("CLICKHOUSE_PORT"),     // "9000"
		Password:        os.Getenv("CLICKHOUSE_PASSWORD"), // "qwerty123"
		Username:        os.Getenv("CLICKHOUSE_USER"),     // "default"
		DBName:          os.Getenv("CLICKHOUSE_NAME"),     // "test"
		TableName:       "test",
		DialTimeout:     2,
		MaxOpenConns:    2,
		MaxIdleConns:    2,
		ConnMaxLifetime: 10,
		Debug:           true,
	}

	ctx := context.Background()
	db, err := NewClickHouseDB(ctx, cnf)
	require.NoError(t, err)

	err = db.Ping(ctx)
	require.NoError(t, err)

	defer db.CloseData()

	err = db.db.Exec(
		ctx,
		fmt.Sprintf(createTemplate, cnf.DBName, cnf.TableName),
	)
	require.NoError(t, err)

	defer func() {
		err = db.db.Exec(
			ctx,
			fmt.Sprintf(dropTemplate, cnf.DBName, cnf.TableName),
		)
		require.NoError(t, err)
	}()

	event := DataEventModel{
		ClientTime:      time.Now().UTC(),
		ServerTime:      time.Now().UTC(),
		DeviceId:        "0287D9AA-4ADF-4B37-A60F-3E9E645C821E",
		Session:         "dfb",
		ParamStr:        "fgg",
		Ip:              1234567890,
		Sequence:        1,
		ParamInt:        1234,
		DeviceOs:        0,
		DeviceOsVersion: 1,
		Event:           10,
	}

	// single insert
	payload := []DataEventModel{event}
	err = db.Insert(ctx, payload)
	require.NoError(t, err)

	// batch insert
	payload1 := make([]DataEventModel, 0, 1000)
	for i := 0; i < 1000; i++ {
		payload1 = append(payload1, event)
	}
	err = db.Insert(ctx, payload1)
	require.NoError(t, err)
}
