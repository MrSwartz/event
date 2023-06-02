package data

import (
	"context"
	"event/internal/config"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInsert(t *testing.T) {
	cnf := config.ClickHouse{
		Host:     "127.0.0.1",
		Port:     "9000",
		Password: "qwerty123",
		Username: "default",
		DBName:   "default",
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(2)*time.Second)
	defer cancel()

	db, err := NewClickHouseDB(ctx, cnf)
	require.NoError(t, err)

	err = db.db.Ping(ctx)
	require.NoError(t, err)

	defer db.CloseData()

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
