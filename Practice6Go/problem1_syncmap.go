package main

import (
	"fmt"
	"sync"
)

func main() {
	var m sync.Map
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			m.Store("key", val)
		}(i)
	}

	wg.Wait()

	value, _ := m.Load("key")
	fmt.Printf("Value: %d\n", value)
}