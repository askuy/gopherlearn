package main

import (
	"fmt"
	"sync"
	"time"
)

type monitor struct {
	notifyChan chan int
}

var monitorObj *monitor

func main() {
	monitorObj = &monitor{
		notifyChan: make(chan int, 10),
	}
	go monitorObj.syncStore()

	t1 := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			monitorObj.push(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println(time.Since(t1))
}

func (m *monitor) push(id int) {
	m.notifyChan <- id
}

func (m *monitor) syncStore() {
	for {
		select {
		case id := <-m.notifyChan:
			//time.Sleep(time.Nanosecond)
			time.Sleep(time.Microsecond)
			fmt.Println(id)
		}
	}
}
