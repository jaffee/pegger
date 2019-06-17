package main

import (
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"

	"github.com/jaffee/commandeer"
)

type Main struct {
	Concurrency int    `help:"Number of goroutines to spawn."`
	Iterations  uint64 `help:"Number of times each routine should loop."`
	Profiling   string `help:"Bind address for pprof."`
}

func NewMain() *Main {
	return &Main{
		Concurrency: runtime.NumCPU(),
		Iterations:  1 << 40,
		Profiling:   "localhost:6060",
	}
}

func (m *Main) Run() error {
	go func() {
		log.Println(http.ListenAndServe(m.Profiling, nil))
	}()

	wg := sync.WaitGroup{}
	vals := make([]int, m.Concurrency)

	for i := 0; i < m.Concurrency; i++ {
		vals[i] = rand.Int()
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			rnd := rand.New(rand.NewSource(rand.Int63()))
			for j := uint64(0); j < m.Iterations; j++ {
				vals[idx] = vals[idx] ^ rnd.Int()
			}
		}(i)
	}

	wg.Wait()
	log.Println(vals)
	return nil
}

func main() {
	if err := commandeer.Run(NewMain()); err != nil {
		log.Fatal(err)
	}
}
