package data

import "time"

type DeviceOS uint8

// 	unsupported DeviceOS = 0
// 	ios         DeviceOS = 1
// 	android     DeviceOS = 2
// 	linux       DeviceOS = 3

type DeviceOSVersion uint16

//  version 13.5.1 -> 11351
// 	i1351 DeviceOSVersion = 11351 // ios(1) 13.5.1
// 	i1352 DeviceOSVersion = 11352 // ios(1) 13.5.2
// 	i1353 DeviceOSVersion = 11353 // ios(1) 13.5.3

type EventType uint8

// 	appStart  EventType = 0
// 	onPause   EventType = 1
// 	onRotate  EventType = 2

// так clickhouse работает быстрее + буфер занимает меньше памяти
type DataEventModel struct {
	ClientTime      time.Time
	ServerTime      time.Time
	DeviceId        string
	Session         string
	ParamStr        string
	Ip              uint32 // возможно не лучшее решение, но сделал так
	Sequence        uint32
	ParamInt        uint32
	DeviceOs        DeviceOS
	DeviceOsVersion DeviceOSVersion
	Event           EventType
}
