package json

import "github.com/dllen/go-crawler/model"

type JsonProcess struct {
	jsonProcess *model.Process
}

func NewJsonProcess(jsonProcess *model.Process) *JsonProcess {
	return &JsonProcess{jsonProcess: jsonProcess}
}

func (j *JsonProcess) Process(context model.Context) (*model.Page, error) {
	return JsonRuleProcess(j.jsonProcess, context)
}
