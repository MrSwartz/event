package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

type Clickhouse struct {
	Host            string
	Port            string
	Username        string
	Password        string
	DBName          string
	DialTimeout     uint32
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime uint32
	Debug           bool
}

type Service struct {
	ExposeSwagger bool
	Port          int
	ReadTimeout   int
	WriteTimeout  int
	IdleTimeout   int
}

type Config struct {
	DataBase Clickhouse
	Buffer   Buffer
	Service  Service
	Mappers  Mappers
}

type Buffer struct {
	RetriesLeft int
	LoopTimeout int
	Size        int
}

type Mappers struct {
	Events    map[string]uint8
	DeviceOs  map[string]uint8
	OsVersion map[string]uint16
}

func ReadConfig() (*Config, error) {
	cnf, err := readConfig()
	if err != nil {
		return nil, err
	}

	mp, err := readMappers()
	if err != nil {
		return nil, err
	}

	cnf.Mappers = Mappers{
		DeviceOs:  mp.DeviceOs,
		OsVersion: mp.OsVersion,
		Events:    mp.Events,
	}
	return cnf, nil
}

func readConfig() (*Config, error) {
	configFile := "../cmd/config-" + os.Getenv("ENV") + ".toml"

	file, err := readFile(configFile)
	if err != nil {
		return nil, err
	}

	cnf := Config{}

	if _, err := toml.Decode(string(file), &cnf); err != nil {
		return nil, err
	}

	// эти параметры лучше ложить например в Vault, но не хочется усложнять, поэтому сделал так
	cnf.DataBase.DBName = os.Getenv("CLICKHOUSE_NAME")
	cnf.DataBase.Host = os.Getenv("CLICKHOUSE_HOST")
	cnf.DataBase.Password = os.Getenv("CLICKHOUSE_PASSWORD")
	cnf.DataBase.Port = os.Getenv("CLICKHOUSE_PORT")
	cnf.DataBase.Username = os.Getenv("CLICKHOUSE_USER")

	return &cnf, nil
}

func readMappers() (*Mappers, error) {

	devos := "../mappers/device_os.json"
	file, err := readFile(devos)
	if err != nil {
		return nil, err
	}

	mapDevOs := make(map[string]uint8)
	if err := json.Unmarshal(file, &mapDevOs); err != nil {
		return nil, err
	}

	events := "../mappers/events.json"
	file, err = readFile(events)
	if err != nil {
		return nil, err
	}

	mapEvents := make(map[string]uint8)
	if err := json.Unmarshal(file, &mapEvents); err != nil {
		return nil, err
	}

	osver := "../mappers/os_version.json"
	file, err = readFile(osver)
	if err != nil {
		return nil, err
	}

	mapOsVer := make(map[string]uint16)
	if err := json.Unmarshal(file, &mapOsVer); err != nil {
		return nil, err
	}

	return &Mappers{
		DeviceOs:  mapDevOs,
		OsVersion: mapOsVer,
		Events:    mapEvents,
	}, nil
}

func readFile(fname string) ([]byte, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return io.ReadAll(file)
}
