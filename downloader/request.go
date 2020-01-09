package downloader

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ddliu/go-httpclient"
	"github.com/dllen/go-crawler/logger"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"sync"
)

var Clients map[string]*httpclient.HttpClient
var lock sync.RWMutex

func init() {
	Clients = make(map[string]*httpclient.HttpClient)
}

func makeCookieJar() http.CookieJar {
	cookieJarOptions := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, _ := cookiejar.New(&cookieJarOptions)

	return jar
}

func makeClient(jar http.CookieJar) *httpclient.HttpClient {
	return httpclient.Defaults(httpclient.Map {
		httpclient.OPT_USERAGENT: "",
		httpclient.OPT_COOKIEJAR: jar,
	})
}

func Get(taskId string, url string, headers map[string]string) (*httpclient.Response, error) {
	res, err := doRequest(taskId, "GET", url, nil, headers)
	if err != nil {
		logger.Info("Download fail doRequest,url:", url, "err:", err)
		return nil, err
	}
	logger.Info("GET", url, " =>", res.StatusCode)
	if res.StatusCode >= 400 {
		return nil, errors.New(fmt.Sprintf("download fail,url %s, StatusCode %d", url, res.StatusCode))
	}
	return res, nil
}

func PostJson(taskId string, url string, data interface{}, headers map[string]string) (*httpclient.Response, error) {
	dataJson, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	res, err := doRequest(taskId, "POST", url, dataJson, headers)
	if err != nil {
		return nil, err
	}
	logger.Info("POST", url, "=>", res.StatusCode)
	if res.StatusCode >= 400 {
		return nil, errors.New(fmt.Sprintf("download fail, StatusCode %d", res.StatusCode))
	}
	return res, nil
}

func doRequest(id string, method string, url string, data []byte, headers map[string]string) (resp *httpclient.Response, err error) {
	client := getClient(id)
	if client == nil {
		client = makeClient(makeCookieJar())
		setClient(id, client)
	}
	return client.Begin().Do(method, url, headers, bytes.NewBuffer(data))
}

func setClient(id string, client *httpclient.HttpClient) {
	lock.Lock()
	defer lock.Unlock()
	Clients[id] = client
}

func getClient(id string) *httpclient.HttpClient {
	lock.RLock()
	defer lock.RUnlock()
	client := Clients[id]
	return client
}
