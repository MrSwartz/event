package service

import (
	"testing"
	"time"

	"github.com/MrSwartz/event/pkg/eventservice/service/data"

	"github.com/stretchr/testify/require"
)

func TestDataModelCast(t *testing.T) {
	initMappers()
	tests := []struct {
		name          string
		payload       ServiceEventModel
		ExpectedModel *data.DataEventModel
		isValid       bool
	}{
		{
			name: "example from task",
			payload: ServiceEventModel{
				ServerTime: time.Date(2023, time.December, 1, 23, 59, 0, 0, time.UTC),
				ClientTime: "2020-12-01 23:59:00",
				DeviceId:   "0287D9AA-4ADF-4B37-A60F-3E9E645C821E",
				Session:    "ybuRi8mAUypxjbxQ",
				ParamStr:   "some text",
				Ip:         "8.8.8.8",
				DeviceOs:   "IOS 13.5.1",
				Event:      "app_start",
				Sequence:   1,
				ParamInt:   0,
			},
			ExpectedModel: &data.DataEventModel{
				ClientTime:      time.Date(2020, time.December, 1, 23, 59, 0, 0, time.UTC),
				ServerTime:      time.Date(2023, time.December, 1, 23, 59, 0, 0, time.UTC),
				DeviceId:        "0287D9AA-4ADF-4B37-A60F-3E9E645C821E",
				Session:         "ybuRi8mAUypxjbxQ",
				ParamStr:        "some text",
				Ip:              0x8080808,
				Sequence:        1,
				ParamInt:        0,
				DeviceOs:        osType["ios"],
				DeviceOsVersion: osVersion["13.5.1"],
				Event:           eventType["app_start"],
			},
			isValid: true,
		},
		{
			name: "invalid client time format",
			payload: ServiceEventModel{
				ServerTime: time.Date(2023, time.December, 1, 23, 59, 0, 0, time.UTC),
				ClientTime: "2020-12-01 23:59:00:00",
			},
			ExpectedModel: nil,
			isValid:       false,
		}, {
			name: "invalid uuid",
			payload: ServiceEventModel{
				ServerTime: time.Date(2023, time.December, 1, 23, 59, 0, 0, time.UTC),
				ClientTime: "2020-12-01 23:59:00",
				DeviceId:   "0287D9AA-4ADF-4B37-A60F-3E9ENOT_HEHE",
			},
			ExpectedModel: nil,
			isValid:       false,
		}, {
			name: "invalid session length",
			payload: ServiceEventModel{
				ServerTime: time.Date(2023, time.December, 1, 23, 59, 0, 0, time.UTC),
				ClientTime: "2020-12-01 23:59:00",
				DeviceId:   "0287D9AA-4ADF-4B37-A60F-3E9E645C821E",
				Session:    "ybuRi8mAUypxjbxQ+useless_payload",
			},
			ExpectedModel: nil,
			isValid:       false,
		}, {
			name: "invalid param_str length",
			payload: ServiceEventModel{
				ServerTime: time.Date(2023, time.December, 1, 23, 59, 0, 0, time.UTC),
				ClientTime: "2020-12-01 23:59:00",
				DeviceId:   "0287D9AA-4ADF-4B37-A60F-3E9E645C821E",
				Session:    "ybuRi8mAUypxjbxQ",
				ParamStr:   "some text",
			},
			ExpectedModel: nil,
			isValid:       false,
		}, {
			name: "unsupported os version template",
			payload: ServiceEventModel{
				ServerTime: time.Date(2023, time.December, 1, 23, 59, 0, 0, time.UTC),
				ClientTime: "2020-12-01 23:59:00",
				DeviceId:   "0287D9AA-4ADF-4B37-A60F-3E9E645C821E",
				Session:    "ybuRi8mAUypxjbxQ",
				ParamStr:   "some text",
				Ip:         "8.8.8.8",
				DeviceOs:   "IOS 13.5.1.0",
			},
			ExpectedModel: nil,
			isValid:       false,
		}, {
			name: "not supported event",
			payload: ServiceEventModel{
				ServerTime: time.Date(2023, time.December, 1, 23, 59, 0, 0, time.UTC),
				ClientTime: "2020-12-01 23:59:00",
				DeviceId:   "0287D9AA-4ADF-4B37-A60F-3E9E645C821E",
				Session:    "ybuRi8mAUypxjbxQ",
				ParamStr:   "some text",
				Ip:         "8.8.8.8",
				DeviceOs:   "IOS 13.5.1",
				Event:      "bad_app_start",
			},
			ExpectedModel: nil,
			isValid:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, ok := tt.payload.toDataModel()
			require.Equal(t, tt.isValid, ok)

			if tt.isValid {
				require.Equal(t, tt.ExpectedModel.ClientTime, obj.ClientTime)
				require.Equal(t, tt.ExpectedModel.ServerTime, obj.ServerTime)
				require.Equal(t, tt.ExpectedModel.DeviceId, obj.DeviceId)
				require.Equal(t, tt.ExpectedModel.Session, obj.Session)
				require.Equal(t, tt.ExpectedModel.ParamStr, obj.ParamStr)
				require.Equal(t, tt.ExpectedModel.Ip, obj.Ip)
				require.Equal(t, tt.ExpectedModel.Sequence, obj.Sequence)
				require.Equal(t, tt.ExpectedModel.ParamInt, obj.ParamInt)
				require.Equal(t, tt.ExpectedModel.DeviceOs, obj.DeviceOs)
				require.Equal(t, tt.ExpectedModel.DeviceOsVersion, obj.DeviceOsVersion)
				require.Equal(t, tt.ExpectedModel.Event, obj.Event)
			} else {
				require.Nil(t, obj)
			}
		})
	}
}
