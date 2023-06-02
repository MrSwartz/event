package data

import (
	"context"
)

const insertQuery = `INSERT INTO events (
	client_time, 
	server_time, 
	device_id, 
	session,
	param_str, 
	ip, 
	sequence, 
	param_int, 
	device_os, 
	device_os_version, 
	event
	)`

type Events interface {
	Insert(ctx context.Context, events []DataEventModel) error
	Ping(ctx context.Context) error
	CloseData() error
}

func (c *ChClient) Insert(ctx context.Context, rows []DataEventModel) error {
	batch, err := c.db.PrepareBatch(ctx, insertQuery)
	if err != nil {
		return err
	}

	for _, v := range rows {
		batch.Append(
			v.ClientTime.UTC(),
			v.ServerTime.UTC(),
			v.DeviceId,
			v.Session,
			v.ParamStr,
			v.Ip,
			v.Sequence,
			v.ParamInt,
			v.DeviceOs,
			v.DeviceOsVersion,
			v.Event,
		)
	}

	return batch.Send()
}

func (i *ChClient) CloseData() error {
	return i.db.Close()
}

func (i *ChClient) Ping(ctx context.Context) error {
	return i.db.Ping(ctx)
}
