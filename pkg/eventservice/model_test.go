package eventservice

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshalling(t *testing.T) {
	tests := []struct {
		name           string
		payload        []byte
		expectedResult []EventModel
		expectedError  error
	}{
		{
			name: "payload from example, but not array",
			payload: []byte(
				`{
					"client_time": "2020-12-01 23:59:00",
					"device_id": "0287D9AA-4ADF-4B37-A60F-3E9E645C821E",
					"device_os": "ios 13.5.1",
					"session":"ybuRi8mAUypxjbxQ",
					"sequence": 1,
					"event": "app_start",
					"param_int": 100,
					"param_str": "some text"
				}`,
			),
			expectedResult: nil,
			expectedError:  errors.New("json: cannot unmarshal object into Go value of type []eventservice.EventModel"),
		}, {
			name: "payload from example, but array",
			payload: []byte(
				`[
					{
						"client_time": "2020-12-01 23:59:00",
						"device_id": "0287D9AA-4ADF-4B37-A60F-3E9E645C821E",
						"device_os": "ios 13.5.1",
						"session":"ybuRi8mAUypxjbxQ",
						"sequence": 1,
						"event": "app_start",
						"param_int": 100,
						"param_str": "some text"
					}
				]`,
			),
			expectedResult: []EventModel{
				{
					ClientTime: "2020-12-01 23:59:00",
					DeviceId:   "0287D9AA-4ADF-4B37-A60F-3E9E645C821E",
					DeviceOs:   "ios 13.5.1",
					Session:    "ybuRi8mAUypxjbxQ",
					Sequence:   1,
					Event:      "app_start",
					ParamInt:   100,
					ParamStr:   "some text",
				},
			},
			expectedError: nil,
		}, {
			name: "payload from example, array of 3 elems",
			payload: []byte(
				`[
					{
						"client_time": "2020-12-01 23:59:00",
						"device_id": "0287D9AA-4ADF-4B37-A60F-3E9E645C821E",
						"device_os": "ios 13.5.1",
						"session":"ybuRi8mAUypxjbxQ",
						"sequence": 1,
						"event": "app_start",
						"param_int": 100,
						"param_str": "some text"
					},{
						"client_time": "2020-12-01 23:59:00",
						"device_id": "0287D9AA-4ADF-4B37-A60F-3E9E645C821E",
						"device_os": "ios 14.5.1",
						"session":"ybuRi8mAUypxjbxQ",
						"sequence": 1,
						"event": "app_start",
						"param_int": 100,
						"param_str": "some text"
					},{
						"client_time": "2020-12-01 23:59:00",
						"device_id": "0287D9AA-4ADF-4B37-A60F-3E9E645C821E",
						"device_os": "ios 15.5.1",
						"session":"ybuRi8mAUypxjbxQ",
						"sequence": 1,
						"event": "app_start",
						"param_int": 100,
						"param_str": "some text"
					}
				]`,
			),
			expectedResult: []EventModel{
				{
					ClientTime: "2020-12-01 23:59:00",
					DeviceId:   "0287D9AA-4ADF-4B37-A60F-3E9E645C821E",
					DeviceOs:   "ios 13.5.1",
					Session:    "ybuRi8mAUypxjbxQ",
					Sequence:   1,
					Event:      "app_start",
					ParamInt:   100,
					ParamStr:   "some text",
				}, {
					ClientTime: "2020-12-01 23:59:00",
					DeviceId:   "0287D9AA-4ADF-4B37-A60F-3E9E645C821E",
					DeviceOs:   "ios 14.5.1",
					Session:    "ybuRi8mAUypxjbxQ",
					Sequence:   1,
					Event:      "app_start",
					ParamInt:   100,
					ParamStr:   "some text",
				}, {
					ClientTime: "2020-12-01 23:59:00",
					DeviceId:   "0287D9AA-4ADF-4B37-A60F-3E9E645C821E",
					DeviceOs:   "ios 15.5.1",
					Session:    "ybuRi8mAUypxjbxQ",
					Sequence:   1,
					Event:      "app_start",
					ParamInt:   100,
					ParamStr:   "some text",
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			models := []EventModel{}
			err := json.Unmarshal(tt.payload, &models)
			if tt.expectedError != nil {
				require.Equal(t, tt.expectedError.Error(), err.Error())
				require.Empty(t, models)
			} else {
				require.Equal(t, len(tt.expectedResult), len(models))

				for i, obj := range tt.expectedResult {
					require.Equal(t, obj.ClientTime, models[i].ClientTime)
					require.Equal(t, obj.DeviceId, models[i].DeviceId)
					require.Equal(t, obj.DeviceOs, models[i].DeviceOs)
					require.Equal(t, obj.Session, models[i].Session)
					require.Equal(t, obj.Sequence, models[i].Sequence)
					require.Equal(t, obj.Event, models[i].Event)
					require.Equal(t, obj.ParamInt, models[i].ParamInt)
					require.Equal(t, obj.ParamStr, models[i].ParamStr)
				}
			}
		})
	}
}
