package service

import (
	"time"

	"github.com/MrSwartz/event/internal/utils"
	"github.com/MrSwartz/event/pkg/eventservice/service/data"

	"github.com/sirupsen/logrus"
)

type ServiceEventModel struct {
	ServerTime time.Time
	ClientTime string
	DeviceId   string
	Session    string
	ParamStr   string
	Ip         string
	DeviceOs   string
	Event      string
	Sequence   uint32
	ParamInt   uint32
}

func (e ServiceEventModel) toDataModel() (*data.DataEventModel, bool) {
	/*
		идея такой строгой валидации заключается в идее более эффективно хранть
		данные в бд и валидировать данные, чтобы хранить только валидные значения

		про пакет "github.com/go-playground/validator/v10" знаю
		не использую его так как лишний проход по данным в рефлексии
		замедлит проверку, поэтому я совместил её с конвертером и сделал её через if

		вторым параметром функция возваращает bool, который нужен чтобы
		была возможность пропустить событие так как драйвер clickhouse
		слишком консервативный и падает в панике при попытке передать
		невалидные значения

		так же события с отсутствующими полями может внести некоторые
		"аномалии" в данных при попытке аггрегировать их, что может
		быть достаточно критичным
	*/
	clientTime, err := time.Parse("2006-01-02 15:04:05", e.ClientTime)
	if err != nil {
		logrus.Errorf("wrong time format clientTime: %v", e.ClientTime)
		return nil, false
	}

	if ok := utils.CheckUUID(e.DeviceId); !ok {
		logrus.Errorf("not valid UUID format: %v", e.DeviceId)
		return nil, false
	}

	if len(e.Session) > 16 {
		logrus.Errorf("session length must be not longer than 16 symbols")
		return nil, false
	}

	if len(e.ParamStr) > 31 {
		logrus.Errorf("param_str length must be not longer than 31 symbols")
		return nil, false
	}

	ipcode := utils.IPstringV4ToInt(e.Ip)
	if ipcode == 0 {
		logrus.Errorf("not valid ip addr: %v", e.Ip)
		return nil, false
	}

	os, version, err := utils.SplitOsAndVersion(e.DeviceOs)
	if err != nil {
		logrus.Errorf("can't extract os and version from incoming data: %v", e.DeviceOs)
		return nil, false
	}

	osCode, ok := osType[os]
	if !ok {
		logrus.Printf("my: %v, list: %v", version, osType)
		logrus.Errorf("os not supported: %v", version)
		return nil, false
	}

	versionCode, ok := osVersion[version]
	if !ok {
		logrus.Errorf("os version not supported: %v", versionCode)
		return nil, false
	}

	eventCode, ok := eventType[e.Event]
	if !ok {
		logrus.Errorf("event not supported: %v", e.Event)
		return nil, false
	}

	return &data.DataEventModel{
		ClientTime: clientTime,
		ServerTime: e.ServerTime,
		DeviceId:   e.DeviceId,
		Session:    e.Session,
		ParamStr:   e.ParamStr,
		Ip:         ipcode,
		Sequence:   e.Sequence,
		ParamInt:   e.ParamInt,
		//nolint:unconvert
		DeviceOs: data.DeviceOS(osCode),
		//nolint:unconvert
		DeviceOsVersion: data.DeviceOSVersion(versionCode),
		//nolint:unconvert
		Event: data.EventType(eventCode),
	}, true
}
