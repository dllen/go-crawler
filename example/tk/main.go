package main

import (
	crawler "github.com/dllen/go-crawler"
	"github.com/dllen/go-crawler/model"
	"github.com/dllen/go-crawler/spider"
)

func main() {
	task := &model.Task{
		ID:   "tuiku",
		Name: "tuiku",
		Request: []*model.Request{
			{
				Method:      "get",
				Url:         "http://www.tuicool.com/ah/0/{1-100,1}?lang=1",
				ProcessName: "tuikulist",
			},
			{
				Method:      "get",
				Url:         "http://www.tuicool.com/ah/101000000/{1-100,1}?lang=1",
				ProcessName: "tuikulist",
			},
			{
				Method:      "get",
				Url:         "http://www.tuicool.com/ah/101040000/{1-100,1}?lang=1",
				ProcessName: "tuikulist",
			},
			{
				Method:      "get",
				Url:         "http://www.tuicool.com/ah/101050000/{1-100,1}?lang=1",
				ProcessName: "tuikulist",
			},
			{
				Method:      "get",
				Url:         "http://www.tuicool.com/ah/20/{1-100,1}?lang=1",
				ProcessName: "tuikulist",
			},
			{
				Method:      "get",
				Url:         "http://www.tuicool.com/ah/108000000/{1-100,1}?lang=1",
				ProcessName: "tuikulist",
			},
			{
				Method:      "get",
				Url:         "http://www.tuicool.com/ah/114000000/{1-100,1}?lang=1",
				ProcessName: "tuikulist",
			},
		},
		Process: []model.Process{
			{
				Name: "tuikulist",
				Type: "template",
				TemplateRule: model.TemplateRule{
					Rule: map[string]string{
						"node":   "array|.list_article_item",
						"img":    "attr.src|.article_thumb_image img",
						"title":  "text|.title a",
						"author": "text|.tip span:nth-child(1)",
						"time":   "text|.tip span:nth-child(3)",
					},
				},
			},
		},
		Pipline: "console",
	}
	app := crawler.New()
	app.AddSpider(spider.InitWithTask(task))
	app.Run()
}
