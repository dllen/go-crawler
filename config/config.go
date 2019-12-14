package config

import (
	"encoding/json"

	"github.com/BurntSushi/toml"
	"github.com/dllen/go-crawler/logger"
)

var Conf *Config
var metaData toml.MetaData

type Config struct {
	Name       string `toml:"name"`
	Version    string `toml:"version"`
	WorkNum    int    `toml:"work_num"`
	MaxWaitNum int    `toml:"max_wait_num"`

	RedisAddr    string   `toml:"redis_addr"`
	ScheduleMode string   `toml:"schedule"`
	Etcd         []string `toml:"etcd"`

	Mysql string `toml:"mysql"`
}

func InitConfig() error {
	data, err := toml.DecodeFile("crawler.toml", &Conf)
	if err != nil {
		return err
	}
	bytes, _ := json.Marshal(Conf)
	logger.Info("Init config ", string(bytes))
	metaData = data
	return nil
}
