package pipeline

import (
	"encoding/json"
	"fmt"
)

type ConsolePipeline struct {
}

func NewConsolePipeline() ConsolePipeline {
	return ConsolePipeline{}
}

func (c ConsolePipeline) ProcessData(v []map[string]interface{}, taskName string, processName string) {
	bytes, _ := json.Marshal(v)
	fmt.Println("Pipeline :", string(bytes))
}
