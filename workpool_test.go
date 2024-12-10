package workpool_test

import (
	"fmt"
	"github.com/sam-rba/workpool"
)

// Vector dot product A·B.
func Example() {
	// Input vectors.
	// Imagine these are very long—much longer than the number of CPU threads.
	a := []int{0, 1, 2, 3, 4}
	b := []int{5, 6, 7, 8, 9}

	// Sum all the a[i]*b[i] terms in another goroutine.
	terms, total := make(chan int), make(chan int)
	go sum(terms, total)

	// Multiply each pair of elements concurrently.
	pool := workpool.New(workpool.DefaultSize)
	defer pool.Close()
	for i := range a {
		pool.Spawn(func() {
			terms <- a[i] * b[i]
		})
	}
	pool.Wait()          // wait for tasks to finish.
	close(terms)         // signal completion to the sum goroutine.
	fmt.Println(<-total) // read the result.
	// Output: 80
}

func sum(terms <-chan int, total chan<- int) {
	defer close(total)
	sum := 0
	for x := range terms {
		sum += x
	}
	total <- sum
}
