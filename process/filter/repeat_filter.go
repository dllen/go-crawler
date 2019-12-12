package filter

import (
	"crypto/md5"
	"fmt"
	"github.com/dllen/go-crawler/model"
	"sync"
)

var Filter map[string]int
var lock sync.RWMutex

func init() {
	Filter = make(map[string]int)
}

func RepeatFilter(url string, process *model.Process) bool {
	data := []byte(url)
	sign := fmt.Sprintf("%x", md5.Sum(data))
	if ok := get(sign); ok {
		return false
	}
	put(sign)
	return true
}

func get(str string) bool {
	lock.RLock()
	defer lock.RUnlock()
	_, ok := Filter[str]
	return ok
}

func put(str string) {
	lock.Lock()
	defer lock.Unlock()
	Filter[str] = 1
}
