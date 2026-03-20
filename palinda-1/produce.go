package main

import (
	"fmt"
	"sync"
	"time"
)

var cheesePortions = 5
var mu sync.Mutex

func chef(chefID int, orderChan <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for order := range orderChan {
		fmt.Printf("Chef %d received order: %s\n", chefID, order)
		mu.Lock()
		hasCheese := false
		if cheesePortions > 0 {
			cheesePortions--
			hasCheese = true
		}
		mu.Unlock()

		if hasCheese {
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("--> Chef %d finished baking %s\n", chefID, order)
		} else {
			fmt.Printf("--> Chef %d cannot make %s: OUT OF CHEESE!\n", chefID, order)
		}
	}
}

func main() {
	var wg sync.WaitGroup
	orderChan := make(chan string)

	fmt.Println("Pizzeria is open!")

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go chef(i, orderChan, &wg)
	}

	orders := []string{"Pepperoni", "Margherita", "Hawaiian", "Meat Lovers", "Veggie", "BBQ Chicken", "Mushroom"}

	for _, order := range orders {
		orderChan <- order
	}

	close(orderChan)

	wg.Wait()

	fmt.Println("Pizzeria is closed for the night.")
}
