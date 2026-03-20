package main

import (
	"fmt"
	"sync"
)

var mutex sync.Mutex
var cheese = 1000

func makePizza(wg *sync.WaitGroup) {

	defer wg.Done()

	mutex.Lock()
	if cheese > 0 {
		cheese--
	}
	defer mutex.Unlock()
}

func main() {

	var wg sync.WaitGroup

	fmt.Println("Starting kitchen")

	for i := +0; i < 1000; i++ {
		wg.Add(1)

		go makePizza(&wg)
	}

	wg.Wait()

	fmt.Printf("Chesse remainin %d\n", cheese)
}
