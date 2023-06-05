package config

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type Clickhouse struct {
	Host            string
	Port            string
	Username        string
	Password        string
	DBName          string
	TableName       string
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
}

type Buffer struct {
	RetriesLeft int
	LoopTimeout int
	Size        int
}

func ReadConfig() (*Config, error) {
	configFile, err := readConfigName()
	if err != nil {
		return nil, err
	}

	return readConfig(configFile)
}

func readConfig(configFile string) (*Config, error) {
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

func readFile(fname string) ([]byte, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return io.ReadAll(file)
}

func readConfigName() (string, error) {
	args := os.Args
	if len(args) == 2 {
		if !strings.HasPrefix(args[1], "--cfg=") {
			return "", errors.New("invalid command")
		}

		configName, ok := strings.CutPrefix(args[1], "--cfg=")
		if !ok {
			return "", errors.New("can't extract config name")
		}

		return configName, nil
	}

	return "", errors.New("try to use ./binary --cfg=config-env.toml")
}
