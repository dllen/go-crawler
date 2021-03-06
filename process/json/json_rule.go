package json

import (
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/dllen/go-crawler/common"
	"github.com/dllen/go-crawler/logger"
	"github.com/dllen/go-crawler/model"
)

func JsonRuleProcess(process *model.Process, context model.Context) (*model.Page, error) {
	return Process(process, context)
}

func Process(process *model.Process, context model.Context) (*model.Page, error) {
	jsonRule := process.JsonRule.Rule
	page := &model.Page{}

	sJson, err := simplejson.NewJson(context.Body)
	if err != nil {
		logger.Error("NewDocumentFromReader fail,", err)
		return nil, err
	}

	resultType := "map"
	rootSel := []string{}

	v, ok := jsonRule["node"]

	if ok {
		contentInfo := strings.Split(v, "|")
		resultType = contentInfo[0]
		selStr := contentInfo[1]
		rootSel = strings.Split(selStr, ".")
	}

	if resultType == "array" {
		for _, name := range rootSel {
			sJson = sJson.Get(name)
		}
		rootNode, err := sJson.Array()
		if err != nil {
			logger.Error("Json fail,", err)
			return nil, err
		}
		if len(rootNode) >= 0 {
			for _, node := range rootNode {
				nodeMap, ok := node.(map[string]interface{})
				if !ok {
					continue
				}
				data := map[string]interface{}{}
				for key, value := range jsonRule {
					if key == "node" {
						continue
					}
					data[key] = nodeMap[value]
				}
				if len(process.AddQueue) > 0 {
					page.AddUrls(common.ParseReq(process.AddQueue, data))
				}
				page.AddResult(data)
			}
		}
	}

	if resultType == "map" {

		result := map[string]interface{}{}

		for _, name := range rootSel {
			sJson = sJson.Get(name)
		}

		if err != nil {
			logger.Error("Json fail,", err)
			return nil, err
		}

		for key, value := range jsonRule {
			valueSel := []string{}
			valueSel = strings.Split(value, ".")
			valueNode := *sJson
			for _, name := range valueSel {
				valueNode = *valueNode.Get(name)
			}
			result[key] = valueNode.Interface()
		}

		if len(process.AddQueue) > 0 {
			page.AddUrls(common.ParseReq(process.AddQueue, result))
		}
		page.AddResult(result)
	}

	if resultType == "nil" {

		result := map[string]interface{}{}

		for _, name := range rootSel {
			sJson = sJson.Get(name)
		}
		rootNode, err := sJson.Map()

		if err != nil {
			logger.Error("Json fail,", err)
			return nil, err
		}

		for key, value := range jsonRule {
			result[key] = rootNode[value]
		}
		page.Urls = []*model.Request{}
		if len(process.AddQueue) > 0 {
			page.AddUrls(common.ParseReq(process.AddQueue, result))
		}
	}
	return page, nil
}
