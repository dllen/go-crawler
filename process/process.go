package process

import "github.com/dllen/go-crawler/model"

type Process interface {
	Process(context model.Context) (*model.Page, error)
}
