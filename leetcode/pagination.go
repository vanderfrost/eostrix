package leetcode

import "sync"

type PageData struct {
	Problems []*Problem
	Page     int
	PageSize int
}

var (
	pageStore = make(map[string]*PageData)
	pageMutex sync.Mutex
)

func StorePage(userID string, data *PageData) {
	pageMutex.Lock()
	defer pageMutex.Unlock()
	pageStore[userID] = data
}

func GetPage(userID string) (*PageData, bool) {
	pageMutex.Lock()
	defer pageMutex.Unlock()
	data, ok := pageStore[userID]
	return data, ok
}

func DeletePage(userID string) {
	pageMutex.Lock()
	defer pageMutex.Unlock()
	delete(pageStore, userID)
}
