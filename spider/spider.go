package spider

import (
	"github.com/dllen/go-crawler/model"
	"github.com/dllen/go-crawler/pipline"
	"github.com/dllen/go-crawler/process"
	json_process "github.com/dllen/go-crawler/process/json"
	template_process "github.com/dllen/go-crawler/process/template"
)

type Spider struct {
	ID       string
	Name     string
	Depth    int
	EndCount int
	Requests []*model.Request
	Process  map[string][]process.Process
	Pipline  pipline.Pipline
}

func (s *Spider) GetPipline() pipline.Pipline {
	return s.Pipline
}

func (s *Spider) GetProcess(name string) []process.Process {
	return s.Process[name]
}

func (s *Spider) GetRequests() []*model.Request {
	return s.Requests
}

func (s *Spider) AddProcess(name string, p process.Process) {
	if s.Process == nil {
		s.Process = make(map[string][]process.Process)
	}
	processs, ok := s.Process[name]
	if !ok {
		ps := []process.Process{}
		s.Process[name] = append(ps, p)
	} else {
		processs = append(processs, p)
	}
}

func InitWithTask(task *model.Task) *Spider {
	s := &Spider{}
	s.ID = task.ID
	s.Name = task.Name
	s.Depth = task.Depth
	s.EndCount = task.EndCount
	s.Requests = task.Request
	s.Process = make(map[string][]process.Process)
	for i, p := range task.Process {
		switch p.Type {
		case "template":
			processs, ok := s.Process[p.Name]
			if !ok {
				processs = []process.Process{}
				s.Process[p.Name] = processs
			}
			s.Process[p.Name] = append(processs, template_process.NewTemplateProcess(&task.Process[i]))
		case "json":
			processs, ok := s.Process[p.Name]
			if !ok {
				processs = []process.Process{}
				s.Process[p.Name] = processs
			}
			s.Process[p.Name] = append(processs, json_process.NewJsonProcess(&task.Process[i]))
		}
	}
	switch task.Pipline {
	case "console":
		s.Pipline = pipline.NewConsolePipline()
	case "file":
		s.Pipline = pipline.NewFilePipline("./")
	default:
		s.Pipline = pipline.NewConsolePipline()
	}
	return s
}
