package main

import (
	"fmt"
	"github.com/robfig/cron"
	"sync"
)

type store struct {
	sync.RWMutex
	cron *cron.Cron
	data map[int]string
}

var storeObj *store

func main() {
	storeObj = &store{
		data: make(map[int]string, 0),
		cron: cron.New(),
	}
	storeObj.syncStore()
	storeObj.cron.AddFunc("*/10 * * * * *", storeObj.syncStore)
	storeObj.cron.Start()

	var wg sync.WaitGroup
	for i := 1; i < 6; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			fmt.Println(storeObj.Get(i))
			wg.Done()
		}(&wg, i)
	}
	wg.Wait()
}

func (s *store) syncStore() {
	dataTmp := getRemoteData()
	s.Lock()
	s.data = dataTmp
	s.Unlock()
}

func (s *store) Get(id int) string {
	s.RLock()
	defer s.RUnlock()
	resp, _ := s.data[id]
	return resp
}

func getRemoteData() map[int]string {
	resp := make(map[int]string)
	resp[1] = "test1"
	resp[2] = "test2"
	resp[3] = "test3"
	resp[4] = "test4"
	resp[5] = "test5"
	return resp
}
