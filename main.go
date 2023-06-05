package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func randomThrow(r float64) (float64, float64) {
	// Return a random point in the square
	// [(-r, r), (r, r), (r, -r), (-r, -r)]
	x := rand.Float64()*2*r - r
	y := rand.Float64()*2*r - r
	return x, y
}

func inCircle(r float64, x float64, y float64) bool {
	// Return true if (x, y) is in a circle with radius r
	// centered at 0. Otherwise, return false.
	distanceFromCenter := math.Sqrt(x*x + y*y)
	return distanceFromCenter <= r
}

func calculatePi(nTotal int, r float64, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	nCircle := 0
	for i := 0; i < nTotal; i++ {
		x, y := randomThrow(r)
		if inCircle(r, x, y) {
			nCircle++
		}
	}
	results <- nCircle
}

func main() {
	rand.Seed(time.Now().UnixNano())
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	for i := 0; i < 10; i++ {
		nTotal := int(math.Pow10(i))
		nPerCPU := nTotal / numCPU

		wg := sync.WaitGroup{}
		results := make(chan int, numCPU)

		startTime := time.Now()

		for j := 0; j < numCPU; j++ {
			wg.Add(1)
			go calculatePi(nPerCPU, 0.5, results, &wg)
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		nCircle := 0
		for result := range results {
			nCircle += result
		}

		pi := 4.0 * float64(nCircle) / float64(nTotal)
		elapsedTime := time.Since(startTime)
		fmt.Printf("1x10^%d: %f (Elapsed Time: %s)\n", i, pi, elapsedTime)
	}
}
