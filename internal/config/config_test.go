package config

import (
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

func TestReadConfig(t *testing.T) {
	err := os.Setenv("APP_ENV", "test")
	require.NoError(t, err)

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

	cnf, err := readConfig("../../cmd/config-test.toml")
	require.NoError(t, err)
	require.NotNil(t, cnf)

	require.Equal(t, "test1", cnf.DataBase.DBName)
	require.Equal(t, "test2", cnf.DataBase.Host)
	require.Equal(t, "test3", cnf.DataBase.Password)
	require.Equal(t, "test4", cnf.DataBase.Port)
	require.Equal(t, "test5", cnf.DataBase.Username)
	require.Equal(t, uint32(10), cnf.DataBase.DialTimeout, cnf.DataBase)
	require.Equal(t, 5, cnf.DataBase.MaxOpenConns)
	require.Equal(t, 5, cnf.DataBase.MaxIdleConns)
	require.Equal(t, uint32(3600), cnf.DataBase.ConnMaxLifetime)
	require.Equal(t, true, cnf.DataBase.Debug)

	require.Equal(t, true, cnf.Service.ExposeSwagger)
	require.Equal(t, 10, cnf.Service.IdleTimeout)
	require.Equal(t, 10, cnf.Service.ReadTimeout)
	require.Equal(t, 120, cnf.Service.WriteTimeout)
	require.Equal(t, 8080, cnf.Service.Port)

	require.Equal(t, 3, cnf.Buffer.RetriesLeft)
	require.Equal(t, 10, cnf.Buffer.LoopTimeout)
	require.Equal(t, 60000, cnf.Buffer.Size)

}
