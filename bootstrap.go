package crawler

import (
	"github.com/dllen/go-crawler/config"
	"github.com/dllen/go-crawler/core"
	"github.com/dllen/go-crawler/register/etcd"
	"github.com/dllen/go-crawler/spider"
)

type Bootstrap struct {
	engine *core.Engine
}

func init() {
	var err error
	if err = config.InitConfig(); err != nil {
		panic(err)
	}
}

func New() *Bootstrap {
	s := &Bootstrap{}
	s.engine = core.New()
	return s
}

func (s *Bootstrap) AddSpider(spider *spider.Spider) *core.Engine {
	return s.engine.AddSpider(spider)
}

func (s *Bootstrap) Run() {
	s.engine.Run()
	if len(config.Conf.Etcd) > 0 {
		worker := etcd.NewWorker(config.Conf.Name, config.Conf.HttpAddr, config.Conf.Etcd)
		go worker.HeartBeat()
	}
}
