package data

import (
	"context"
	"fmt"
)

var _ Events = (*ChClient)(nil)

const (
	/*
		я знаю про миграции, в любом случае при изменении схемы данных придётся чтото менять в коде
		поэтому я решил не усложнять запуск и захардкодил запрос на создание таблицы если её нету

		возможные вопросы:

		вопрос: почему IP адрес не представлен в виде встроенного типа в CH?
		ответ: я понятия не имею как именно дальше будут использоваться данные,
		поэтому решил кастить в инт + любая бд лучше работает с примитивами, чем с обёртками над данными

		вопрос: почему ORDER BY ip?
		ответ: я незнаю какие именно будут составляться запросы, поэтому выбрал произвольно
	*/
	createTemplate = `
	CREATE TABLE IF NOT EXISTS %s.%s
	(
		client_time 		DateTime('UTC'),
		server_time 		DateTime('UTC'),
		device_id 			UUID,
		session 			FixedString(17),
		param_str 			FixedString(32),
		ip 					UInt32,
		sequence 			UInt32,
		param_int 			UInt32,
		device_os 			UInt8,
		device_os_version 	UInt16,
		event 				UInt8
	)
	ENGINE = MergeTree
	ORDER BY ip;
	`

	insertQuery = `
	INSERT INTO %s.%s (
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

	//nolint:unused
	dropTemplate = `DROP TABLE IF EXISTS %s.%s`
)

type Events interface {
	Insert(ctx context.Context, events []DataEventModel) error
	Ping(ctx context.Context) error
	CloseData() error
}

func (c *ChClient) Insert(ctx context.Context, rows []DataEventModel) error {
	batch, err := c.db.PrepareBatch(ctx, fmt.Sprintf(insertQuery, c.dbName, c.tableName))
	if err != nil {
		return err
	}

	for _, v := range rows {
		err := batch.Append(
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

		if err != nil {
			return err
		}
	}

	return batch.Send()
}

func (i *ChClient) CloseData() error {
	return i.db.Close()
}

func (i *ChClient) Ping(ctx context.Context) error {
	return i.db.Ping(ctx)
}
