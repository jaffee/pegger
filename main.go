package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	conc := flag.Int("conc", 8, "number of goroutines to spawn")
	flag.Parse()

	wg := sync.WaitGroup{}
	vals := make([]int, conc)

	for i := 0; i < conc; i++ {
		vals[i] = rand.Int()
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for j := 0; j < 1000*1000*1000*1000; j++ {
				vals[idx] = vals[idx] ^ rand.Int()
			}
		}(i)
	}

	wg.Wait()
	fmt.Println(vals)
}
