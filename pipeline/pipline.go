package pipeline

type Pipeline interface {
	ProcessData(v []map[string]interface{}, taskName string, processName string)
}
