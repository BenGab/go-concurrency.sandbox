package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var cache = map[int]Book{}
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	wg := &sync.WaitGroup{}
	cacheCh := make(chan Book)
	dbCh := make(chan Book)
	for i := 0; i < 10; i++ {
		id := rnd.Intn(10) + 1
		wg.Add(2)

		go func(id int, wg *sync.WaitGroup, ch chan<- Book) {
			if b, ok := queryCache(id); ok {
				ch <- b
			}
			wg.Done()
		}(id, wg, cacheCh)

		go func(id int, wg *sync.WaitGroup, ch chan<- Book) {
			if b, ok := queryDataBase(id); ok {
				ch <- b
			}
			wg.Done()
		}(id, wg, dbCh)

		go func(cacheCh, dbCh <-chan Book) {
			select {
			case b := <-cacheCh:
				fmt.Println("FROM CAHCE")
				fmt.Println(b)
				<-dbCh
			case b := <-dbCh:
				fmt.Println("FROM DB")
				fmt.Println(b)
			}
		}(cacheCh, dbCh)
		wg.Wait()
	}
}

func queryCache(id int) (Book, bool) {
	b, ok := cache[id]
	return b, ok
}

func queryDataBase(id int) (Book, bool) {
	for _, b := range books {
		if b.ID == id {
			cache[id] = b
			return b, true
		}
	}
	return Book{}, false
}
