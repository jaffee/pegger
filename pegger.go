package pegger

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
)

type Pegger struct {
	Concurrency int    `help:"Number of goroutines to spawn."`
	Iterations  uint64 `help:"Number of times each routine should loop."`
	Profiling   string `help:"Bind address for pprof."`
}

func NewPegger() *Pegger {
	return &Pegger{
		Concurrency: runtime.NumCPU(),
		Iterations:  1 << 40,
		Profiling:   "localhost:6060",
	}
}

func (m *Pegger) Run() error {
	wg := sync.WaitGroup{}
	vals := make([]int, m.Concurrency)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "stuff\n")
	})
	go func() {
		log.Println(http.ListenAndServe(m.Profiling, nil))
	}()

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
	log.Println(vals) // "use" vals so compiler can't optimize computation away.
	return nil
}
