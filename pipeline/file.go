package pipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/dllen/go-crawler/logger"
)

type FilePipeline struct {
	root  string
	files map[string]*os.File
}

func NewFilePipeline(root string) FilePipeline {
	return FilePipeline{root: root, files: make(map[string]*os.File)}
}

func (c FilePipeline) ProcessData(v []map[string]interface{}, taskName string, processName string) {
	file, ok := c.files[processName]
	if !ok {
		var f *os.File
		var err error
		path := fmt.Sprintf("%s%s-%s.txt", c.root, taskName, processName)
		if f, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666); err != nil {
			logger.Error("FilePipeline Open File fail, path =", path, err)
			return
		}
		f.WriteString(fmt.Sprintf("========= Task : %s =============\n", taskName))
		f.WriteString(fmt.Sprintf("======= Task Begin : %s =============\n", time.Now()))
		c.files[processName] = f
		file = f
	}
	for _, value := range v {
		data, err := json.Marshal(value)
		if err != nil {
			logger.Error("FilePipeline json.Marshal fail, v = ", v)
			return
		}
		file.WriteString(string(data) + "\n")
	}
	logger.Info("File Pipeline write. Count:", len(v))
	return
}

func (c FilePipeline) Close() {
	for _, f := range c.files {
		f.Close()
	}
}
