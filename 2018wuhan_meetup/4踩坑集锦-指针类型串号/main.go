package main

import (
	"fmt"
	"github.com/robfig/cron"
	"sync"
)

type store struct {
	sync.RWMutex
	cron *cron.Cron
	data map[int]*user
}

type user struct {
	uid   int
	name  string
	level int
}

var storeObj *store

func main() {
	storeObj = &store{
		cron: cron.New(),
		data: make(map[int]*user, 0),
	}
	storeObj.syncStore()
	storeObj.cron.AddFunc("* */10 * * * *", storeObj.syncStore)
	storeObj.cron.Start()

	// 读模式下问题不大
	testCorrectRead()

	// 指针容易串号,并且这种方式是存在数据竞争a
	// go run --race main.go
	testWrongReadAndWrite()
}

// 模拟多个http请求获取内存中数据，并只读内存数据
func testCorrectRead() {
	fmt.Println("==========内存正确只读模式 开始==========")
	var wg sync.WaitGroup
	for i := 1; i < 6; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			fmt.Println(storeObj.getUser(i))
			wg.Done()
		}(&wg, i)
	}
	wg.Wait()
	fmt.Println("==========内存正确只读模式 结束==========")
}

// 模拟多个http请求获取内存中数据，并只读内存数据
func testWrongReadAndWrite() {
	fmt.Println("==========内存错误的读写模式 开始==========")
	var wg sync.WaitGroup
	for i := 1; i < 6; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			fmt.Println(storeObj.getUser(i))
			wg.Done()
		}(&wg, i)
	}

	wg.Add(1)
	go func(wg *sync.WaitGroup, i int) {
		// 1. 有时候会存在实时计算，调整某个用户级别
		// 2. 有的开发会偷懒，拿着内存数据结构进行简单修改，吐给前端
		// 3. 这种修改方式，会导致内存中数据串号
		user := storeObj.getUser(i)
		user.level = 8
		fmt.Println("modify user===>", storeObj.getUser(i))
		wg.Done()
	}(&wg, 2)
	wg.Wait()
	fmt.Println("==========内存错误的读写模式 结束==========")
}

func (s *store) syncStore() {
	// 1模拟远端拉去大量用户数据放入到内存缓存中
	userList := getRemoterUserList()

	// 2将远端数据进行缓存
	tmpStore := make(map[int]*user, 0)
	for _, user := range userList {
		tmpStore[user.uid] = user
	}
	s.Lock()
	s.data = tmpStore
	s.Unlock()
}

func (s *store) getUser(uid int) *user {
	s.RLock()
	defer s.RUnlock()
	resp, ok := s.data[uid]
	if !ok {
		return nil
	}
	return resp
}

func getRemoterUserList() []*user {
	resp := make([]*user, 0)
	resp = append(resp, &user{
		1, "hello1", 1,
	}, &user{
		2, "hello2", 2,
	}, &user{
		3, "hello3", 3,
	}, &user{
		4, "hello4", 4,
	}, &user{
		5, "hello5", 5,
	})
	return resp
}
