package eventservice

import (
	"event/pkg/eventservice/service"
	"time"
)

type EventModel struct {
	ClientTime string `json:"client_time"`
	DeviceId   string `json:"device_id"`
	DeviceOs   string `json:"device_os"`
	Session    string `json:"session"`
	Event      string `json:"event"`
	ParamStr   string `json:"param_str"`
	Sequence   uint32 `json:"sequence"`
	ParamInt   uint32 `json:"param_int"`
}

func (e EventModel) toServiseDataEventModel(nowTime time.Time, ip string) *service.ServiceEventModel {
	return &service.ServiceEventModel{
		ClientTime: e.ClientTime,
		ServerTime: nowTime,
		DeviceId:   e.DeviceId,
		Session:    e.Session,
		ParamStr:   e.ParamStr,
		Ip:         ip,
		Sequence:   e.Sequence,
		ParamInt:   e.ParamInt,
		DeviceOs:   e.DeviceOs,
		Event:      e.Event,
	}
}
