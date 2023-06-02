package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadFile(t *testing.T) {
	data, err := readFile("./data/test.json")
	require.NoError(t, err)

	expectedData := []uint8([]byte{
		0x7b, 0xa, 0x20, 0x20, 0x20, 0x20, 0x22, 0x6b,
		0x65, 0x79, 0x31, 0x22, 0x3a, 0x20, 0x31, 0x2c,
		0xa, 0x20, 0x20, 0x20, 0x20, 0x22, 0x6b, 0x65,
		0x79, 0x32, 0x22, 0x3a, 0x20, 0x32, 0x2c, 0xa,
		0x20, 0x20, 0x20, 0x20, 0x22, 0x6b, 0x65, 0x79,
		0x33, 0x22, 0x3a, 0x20, 0x33, 0x2c, 0xa, 0x20,
		0x20, 0x20, 0x20, 0x22, 0x6b, 0x65, 0x79, 0x34,
		0x22, 0x3a, 0x20, 0x34, 0xa, 0x7d})
	require.Equal(t, expectedData, data)

	// test open not existing file

	data, err = readFile("./data/notexists.json")
	require.Nil(t, data)
	require.Equal(t, "open ./data/notexists.json: no such file or directory", err.Error())
}

func TestReadMappers(t *testing.T) {
	mappers, err := readMappers()
	require.NoError(t, err)
	require.NotNil(t, mappers)

	events := map[string]uint8{"appStart": 0x0, "onClose": 0x5, "onCreate": 0x3, "onDestroy": 0x4, "onPause": 0x1, "onRotate": 0x2, "panic": 0x6}
	require.Equal(t, fmt.Sprintf("%v", events), fmt.Sprintf("%v", mappers.Events))

	deviceOs := map[string]uint8{"android": 0x2, "ios": 0x1, "linux": 0x4, "macos": 0x5, "unix": 0x4, "unsupported": 0x0, "windows": 0x3}
	require.Equal(t, fmt.Sprintf("%v", deviceOs), fmt.Sprintf("%v", mappers.DeviceOs))

	osVersion := map[string]uint16{"10.0.1": 0x5209, "13.5.1": 0x2c57, "13.5.2": 0x2c58, "13.5.3": 0x2c59, "4.4.4": 0x98c, "5.0.1": 0x9c5}
	require.Equal(t, fmt.Sprintf("%v", osVersion), fmt.Sprintf("%v", mappers.OsVersion))
}

func TestReadConfig(t *testing.T) {
	err := os.Setenv("ENV", "test")
	require.NoError(t, err)
	/*
		export CLICKHOUSE_NAME=default
		export CLICKHOUSE_HOST=127.0.0.1
		export CLICKHOUSE_PASSWORD=qwerty123
		export CLICKHOUSE_PORT=8123
		export CLICKHOUSE_USER=default
	*/
	err = os.Setenv("CLICKHOUSE_NAME", "test1")
	require.NoError(t, err)

	err = os.Setenv("CLICKHOUSE_HOST", "test2")
	require.NoError(t, err)

	err = os.Setenv("CLICKHOUSE_PASSWORD", "test3")
	require.NoError(t, err)

	err = os.Setenv("CLICKHOUSE_PORT", "test4")
	require.NoError(t, err)

	err = os.Setenv("CLICKHOUSE_USER", "test5")
	require.NoError(t, err)

	cnf, err := ReadConfig()
	require.NoError(t, err)
	require.NotNil(t, cnf)

	require.Equal(t, "test1", cnf.DataBase.DBName)
	require.Equal(t, "test2", cnf.DataBase.Host)
	require.Equal(t, "test3", cnf.DataBase.Password)
	require.Equal(t, "test4", cnf.DataBase.Port)
	require.Equal(t, "test5", cnf.DataBase.Username)

	require.Equal(t, 10, cnf.Service.IdleTimeout)
	require.Equal(t, 10, cnf.Service.ReadTimeout)
	require.Equal(t, 120, cnf.Service.WriteTimeout)
	require.Equal(t, 8080, cnf.Service.Port)

	events := map[string]uint8{"appStart": 0x0, "onClose": 0x5, "onCreate": 0x3, "onDestroy": 0x4, "onPause": 0x1, "onRotate": 0x2, "panic": 0x6}
	require.Equal(t, fmt.Sprintf("%v", events), fmt.Sprintf("%v", cnf.Mappers.Events))

	deviceOs := map[string]uint8{"android": 0x2, "ios": 0x1, "linux": 0x4, "macos": 0x5, "unix": 0x4, "unsupported": 0x0, "windows": 0x3}
	require.Equal(t, fmt.Sprintf("%v", deviceOs), fmt.Sprintf("%v", cnf.Mappers.DeviceOs))

	osVersion := map[string]uint16{"10.0.1": 0x5209, "13.5.1": 0x2c57, "13.5.2": 0x2c58, "13.5.3": 0x2c59, "4.4.4": 0x98c, "5.0.1": 0x9c5}
	require.Equal(t, fmt.Sprintf("%v", osVersion), fmt.Sprintf("%v", cnf.Mappers.OsVersion))
}
