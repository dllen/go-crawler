package schedule

import (
	"github.com/dllen/go-crawler/common"
	"github.com/dllen/go-crawler/config"
	"github.com/dllen/go-crawler/logger"
	"github.com/dllen/go-crawler/model"
)

type ChanSchedule struct {
	waitQueue chan *model.Request
}

func NewChanSchedule(config *config.Config) Schedule {
	schedule := &ChanSchedule{}
	schedule.waitQueue = make(chan *model.Request, config.MaxWaitNum)
	return schedule
}

func (d *ChanSchedule) Push(req *model.Request) {
	parseReqs := common.ParseReq([]*model.Request{req}, nil)
	for _, req := range parseReqs {
		logger.Info("Push URL:", req.Url, req.ProcessName, len(d.waitQueue))
		d.waitQueue <- req
	}
}

func (d *ChanSchedule) PushMulti(reqs []*model.Request) {
	ParseReqs := common.ParseReq(reqs, nil)
	for _, req := range ParseReqs {
		logger.Info("Push URL:", req.Url, req.ProcessName, len(d.waitQueue))
		d.waitQueue <- req
	}
}

func (d *ChanSchedule) Pop() (*model.Request, bool) {
	req, ok := <-d.waitQueue
	logger.Info("Pop Url:", req.Url, req.ProcessName, len(d.waitQueue))
	return req, ok
}

func (d *ChanSchedule) Count() int {
	return len(d.waitQueue)
}

func (d *ChanSchedule) Close() {
	close(d.waitQueue)
}

func init() {
	RegisterSchedule("chan", NewChanSchedule)
}
