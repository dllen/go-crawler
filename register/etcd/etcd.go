package etcd

import (
	"context"
	"encoding/json"
	"log"
	"runtime"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/dllen/go-crawler/core"
)

type Worker struct {
	Name    string
	IP      string
	KeysAPI client.KeysAPI
}

type WorkerInfo struct {
	Name       string                 `json:"name"`
	IP         string                 `json:"ip"`
	CPU        int                    `json:"cpu"`
	MetaData   map[string]string      `json:"metadata"`
	SpiderData map[string]*SpiderData `json:"spider_data"`
}

type SpiderData struct {
	DownloadFailCount int32 `json:"download_fail_count"`
	DownloadCount     int32 `json:"download_count"`
	URLNum            int32 `json:"url_num"`
	WaitURLNum        int   `json:"wait_url_num"`
	CrawlerResultNum  int32 `json:"crawler_result_num"`
}

func NewWorker(name, IP string, endpoints []string) *Worker {
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}

	w := &Worker{
		Name:    name,
		IP:      IP,
		KeysAPI: client.NewKeysAPI(etcdClient),
	}
	return w
}

func (w *Worker) HeartBeat() {
	api := w.KeysAPI
	for {
		info := &WorkerInfo{
			Name:       w.Name,
			IP:         w.IP,
			CPU:        runtime.NumCPU(),
			SpiderData: getSpiderData(),
		}
		key := "spiders/" + w.Name
		value, _ := json.Marshal(info)
		_, err := api.Set(context.Background(), key, string(value), &client.SetOptions{
			TTL: time.Second * 15,
		})
		if err != nil {
			log.Println("Error update workerInfo:", err)
		}
		time.Sleep(time.Second * 5)
	}
}

func getSpiderData() map[string]*SpiderData {
	spiderDataMap := make(map[string]*SpiderData)
	metas := core.GetEngine().GetTaskMetas()
	for name, meta := range metas {
		data := &SpiderData{}
		data.CrawlerResultNum = meta.CrawlerResultNum
		data.DownloadFailCount = meta.DownloadFailCount
		data.DownloadCount = meta.DownloadCount
		data.WaitURLNum = meta.WaitURLNum
		data.URLNum = meta.URLNum
		spiderDataMap[name] = data
	}
	return spiderDataMap
}
