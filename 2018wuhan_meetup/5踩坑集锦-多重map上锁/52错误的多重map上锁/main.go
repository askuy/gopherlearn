package main

import (
	"fmt"
	"github.com/robfig/cron"
	"sync"
)

type store struct {
	cron *cron.Cron
	data map[int]singleData
}

type singleData struct {
	l    sync.RWMutex
	data map[int]string
}

var storeObj *store

func main() {
	storeObj = &store{
		data: make(map[int]singleData, 0),
		cron: cron.New(),
	}
	storeObj.initStore()
	storeObj.syncStore()
	storeObj.cron.AddFunc("* * * * * *", storeObj.syncStore)
	storeObj.cron.Start()

	var wg sync.WaitGroup
	for j := 1; j < 3; j++ {
		for i := 1; i < 6; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup, j, i int) {
				fmt.Println("j: ", j, " name: ", storeObj.Get(j, i))
				wg.Done()
			}(&wg, j, i)
		}
	}
	wg.Wait()
}

func (s *store) initStore() {
	s.data[1] = singleData{
		data: make(map[int]string, 0),
	}
	s.data[2] = singleData{
		data: make(map[int]string, 0),
	}
}

func (s *store) syncStore() {
	tmpStore1 := getRemoteData()
	tmpStore2 := getRemoteData()
	data1 := s.data[1]
	data2 := s.data[2]
	data1.l.Lock()
	data1.data = tmpStore1
	data1.l.Unlock()
	data2.l.Lock()
	data2.data = tmpStore2
	data2.l.Unlock()
}

func (s *store) Get(level int, id int) string {
	dataInfo, ok := s.data[level]
	if !ok {
		return ""
	}
	dataInfo.l.RLock()
	defer dataInfo.l.RUnlock()

	resp, ok := dataInfo.data[id]
	if !ok {
		return ""
	}
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
