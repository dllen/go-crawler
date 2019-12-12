package config

import "github.com/BurntSushi/toml"

var Conf *Config
var metaData toml.MetaData

type Config struct {
	Name       string `toml:"name"`
	Version    string `toml:"version"`
	WorkNum    int    `toml:"work_num"`
	MaxWaitNum int    `toml:"max_wait_num"`

	HttpAddr     string   `toml:"http_addr"`
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
	metaData = data
	return nil
}
