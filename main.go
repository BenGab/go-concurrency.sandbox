package main

import (
	"fmt"
	"math/rand"
	"time"
)

var cache = map[int]Book{}
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	for i := 0; i < 10; i++ {
		id := rnd.Intn(10) + 1
		go func(id int) {
			if b, ok := queryCache(id); ok {
				fmt.Println("From Cache")
				fmt.Println(b)
			}
		}(id)

		go func(id int) {
			if b, ok := queryDataBase(id); ok {
				fmt.Println("From DB")
				fmt.Println(b)
			}
		}(id)

		time.Sleep(150 * time.Millisecond)
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
