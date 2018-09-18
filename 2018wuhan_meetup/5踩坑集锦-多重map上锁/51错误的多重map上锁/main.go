package main

import (
	"fmt"
	"github.com/robfig/cron"
	"sync"
)

type store struct {
	cron *cron.Cron
	lMap map[int]sync.RWMutex
	data map[int]map[int]string
}

var storeObj *store

func main() {
	storeObj = &store{
		lMap: make(map[int]sync.RWMutex, 0),
		data: make(map[int]map[int]string, 0),
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
	s.lMap[1] = sync.RWMutex{}
	s.data[1] = make(map[int]string, 0)
	s.lMap[2] = sync.RWMutex{}
	s.data[2] = make(map[int]string, 0)
}

func (s *store) syncStore() {
	tmpStore1 := getRemoteData()
	tmpStore2 := getRemoteData()
	lock1 := s.lMap[1]
	lock2 := s.lMap[2]
	lock1.Lock()
	s.data[1] = tmpStore1
	lock1.Unlock()
	lock2.Lock()
	s.data[2] = tmpStore2
	lock2.Unlock()
}

func (s *store) Get(level int, id int) string {
	l, ok := s.lMap[level]
	if !ok {
		return ""
	}
	l.RLock()
	defer l.RUnlock()

	resp, ok := s.data[level][id]
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
