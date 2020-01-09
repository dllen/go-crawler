package core

import (
	"github.com/ddliu/go-httpclient"
	"io/ioutil"
	"sync"
	"sync/atomic"

	"github.com/dllen/go-crawler/common"
	"github.com/dllen/go-crawler/config"
	"github.com/dllen/go-crawler/downloader"
	"github.com/dllen/go-crawler/logger"
	"github.com/dllen/go-crawler/model"
	"github.com/dllen/go-crawler/process"
	"github.com/dllen/go-crawler/schedule"
	"github.com/dllen/go-crawler/spider"
)

const Default_WorkNum = 1

type SpiderRuntime struct {
	sync.Mutex
	workNum     int
	schedule    schedule.Schedule
	spider      *spider.Spider
	stopFlag    bool
	recoverChan chan int
	TaskMeta    *TaskMeta
}

type TaskMeta struct {
	DownloadFailCount int32 `json:"download_fail_count"`
	DownloadCount     int32 `json:"download_fail_count"`
	URLNum            int32 `json:"url_num"`
	WaitURLNum        int   `json:"wait_url_num"`
	CrawlerResultNum  int32 `json:"crawler_result_num"`
}

func NewSpiderRuntime() *SpiderRuntime {
	workNum := config.Conf.WorkNum
	if workNum == 0 {
		workNum = Default_WorkNum
	}
	s := &SpiderRuntime{}
	s.workNum = workNum
	s.schedule = schedule.GetSchedule(config.Conf)
	s.recoverChan = make(chan int)
	meta := &TaskMeta{}
	meta.WaitURLNum = 0
	meta.URLNum = int32(0)
	meta.DownloadCount = int32(0)
	meta.DownloadFailCount = int32(0)
	meta.CrawlerResultNum = int32(0)
	s.TaskMeta = meta

	return s
}

func (s *SpiderRuntime) SetSpider(spider *spider.Spider) {
	s.spider = spider
}

func (s *SpiderRuntime) GetSpider() *spider.Spider {
	return s.spider
}

func (s *SpiderRuntime) Run() {
	if s.stopFlag {
		s.recoverChan <- 1
		return
	}
	for i := 0; i < s.workNum; i++ {
		go s.worker()
	}
	s.schedule.PushMulti(s.spider.GetRequests())
}

func (s *SpiderRuntime) Stop() {
	s.stopFlag = true
}

func (s *SpiderRuntime) worker() {
	context := model.Context{}
	for {
		if s.stopFlag {
			_, ok := <-s.recoverChan
			s.stopFlag = false
			if !ok {
				goto exit
			}
		}
		req, ok := s.schedule.Pop()
		if !ok {
			goto exit
		}
		if req == nil {
			logger.Info("schedule is emply")
			continue
		}
		atomic.AddInt32(&s.TaskMeta.DownloadCount, 1)
		response, err := s.download(req)
		if err != nil {
			logger.Error(err.Error())
			atomic.AddInt32(&s.TaskMeta.DownloadFailCount, 1)
			continue
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		context.Clear()
		context.Body, err = common.ToUtf8(body)
		if err != nil {
			context.Body = body
		}
		context.Request = response.Request
		context.Headers = response.Header
		ps, ok := s.spider.Process[req.ProcessName]
		if !ok {
			response.Body.Close()
			logger.Info("Process is not find ! please call SetProcess|SetTask")
			break
		}
		for _, p := range ps {
			page, err := processWrapper(p, context)
			if err != nil {
				logger.Info("Process fail | ", err.Error())
				continue
			}
			if page == nil {
				logger.Info("Process page is nil")
				continue
			}
			s.TaskMeta.WaitURLNum = s.schedule.Count()
			if page.Urls != nil && len(page.Urls) > 0 {
				atomic.AddInt32(&s.TaskMeta.URLNum, int32(len(page.Urls)))
				go func() {
					s.schedule.PushMulti(page.Urls)
				}()
			}
			if page.ResultCount > 0 {
				atomic.AddInt32(&s.TaskMeta.CrawlerResultNum, int32(page.ResultCount))
				s.spider.Pipline.ProcessData(page.Result, s.spider.Name, req.ProcessName)
			}
		}

		response.Body.Close()
	}

exit:
	logger.Info(s.spider.Name, "worker close")
}
func processWrapper(p process.Process, context model.Context) (*model.Page, error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()

	page, err := p.Process(context)
	return page, err
}

func (s *SpiderRuntime) download(req *model.Request) (*httpclient.Response, error) {
	switch req.Method {
	case "get":
		return downloader.Get(req.ProcessName, req.Url,req.Header)
	case "post":
		return downloader.PostJson(req.ProcessName, req.Url, req.Data, req.Header)
	}

	return nil, nil
}

func (s *SpiderRuntime) Exit() {
	s.schedule.Close()
	close(s.recoverChan)
}
