package service

import (
	"event/internal/utils"
	"event/pkg/eventservice/service/data"
	"time"

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
	clientTime, err := time.Parse("2006-01-02 15:04:05", e.ClientTime)
	if err != nil {
		logrus.Errorf("wrong time format clientTime: %v", e.ClientTime)
		return nil, false
	}

	if ok := utils.CheckUUID(e.DeviceId); !ok {
		logrus.Errorf("not valid UUID format: %v", e.DeviceId)
		return nil, false
	}

	ipcode := utils.IPstringV4ToInt(e.Ip)
	if ipcode == 0 {
		logrus.Errorf("not valid ip addr: %v", e.Ip)
		return nil, false
	}

	os, version, err := utils.SplitOsAndVersion(e.DeviceOs)
	if err != nil {
		logrus.Errorf("can't extract os and version from incoming data: %v", e.DeviceId)
		return nil, false
	}

	osCode, ok := osVersion[version]
	if !ok {
		logrus.Errorf("os not supported: %v", version)
		return nil, false
	}

	versionCode, ok := osType[os]
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
		ClientTime:      clientTime,
		ServerTime:      e.ServerTime,
		DeviceId:        e.DeviceId,
		Session:         e.Session,
		ParamStr:        e.ParamStr,
		Ip:              ipcode,
		Sequence:        e.Sequence,
		ParamInt:        e.ParamInt,
		DeviceOs:        data.DeviceOS(osCode),
		DeviceOsVersion: data.DeviceOSVersion(versionCode),
		Event:           data.EventType(eventCode),
	}, true
}
