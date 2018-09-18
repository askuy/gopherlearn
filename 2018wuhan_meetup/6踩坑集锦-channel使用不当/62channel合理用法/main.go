package main

import (
	"fmt"
	"sync"
	"time"
)

type info struct {
	id   int
	time time.Time
}

type monitor struct {
	notifyChan chan info
}

var monitorObj *monitor

func main() {
	monitorObj = &monitor{
		notifyChan: make(chan info, 10),
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
	select {
	case m.notifyChan <- info{
		id:   id,
		time: time.Now(),
	}:
	default:
		// todo something
		fmt.Println("todo something")
	}
}

func (m *monitor) syncStore() {
	for {
		select {
		case info := <-m.notifyChan:
			if time.Now().Sub(info.time) > 1*time.Second {
				fmt.Println("drop this message")
				continue
			}

			//time.Sleep(time.Microsecond)
			//time.Sleep(time.Second)
			fmt.Println(info.id)
		}
	}
}
