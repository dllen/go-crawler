package crawler

import (
	"github.com/dllen/go-crawler/config"
	"github.com/dllen/go-crawler/core"
	"github.com/dllen/go-crawler/spider"
)

func init() {
	var err error
	if err = config.InitConfig(); err != nil {
		panic(err)
	}
}

type Boot struct {
	engine *core.Engine
}

func New() *Boot {
	s := &Boot{}
	s.engine = core.New()
	return s
}

func (s *Boot) AddSpider(spider *spider.Spider) *core.Engine {
	return s.engine.AddSpider(spider)
}

func (s *Boot) Run() {
	s.engine.Run()
}
