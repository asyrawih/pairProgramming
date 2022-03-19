package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Wrapper for sendUser return value, used as result channel type
type Result struct {
	bodyStr string
	err     error
}

// placeholder url. Change it to your base url.
const BASE_URL = "https://jsonplaceholder.typicode.com/posts/"

// number of parallelism
// try using only two worker
const NUM_PARALLEL = 100 

// Stream inputs to input channel 
// Return Pointer To Some Memory Location
func streamInputs(done <-chan struct{}, inputs []string) <-chan string {

  // inputs contains some user id on type Array<string>
	inputCh := make(chan string)

	go func() {
		defer close(inputCh)
		for _, input := range inputs {
			select {
			case inputCh <- input:
			case <-done:
				// in case done is closed prematurely (because error midway),
				// finish the loop (closing input channel)
				break
			}
		}
	}()
	return inputCh
}

// Normal function for HTTP call, no knowledge of goroutine/channels
func sendUser(userId string) (string, error) {
	url := BASE_URL + userId 
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	bodyStr := string(body)
	return bodyStr, nil
}


func AsyncHTTP(users []string) ([]string, error) {
	done := make(chan struct{})
	defer close(done)

	inputChannel := streamInputs(done, users)

	var wg sync.WaitGroup
	// bulk add goroutine counter at the start
	wg.Add(NUM_PARALLEL)

	resultCh := make(chan Result)

	for i := 0; i < NUM_PARALLEL; i++ {
		// spawn N worker goroutines, each is consuming a shared input channel.
		go func() {
			for input := range inputChannel {
        // API CALL
				bodyStr, err := sendUser(input)
				resultCh <- Result{bodyStr, err}
			}
			wg.Done()
		}()
	}

	// Wait all worker goroutines to finish. Happens if there's no error (no early return)
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	results := []string{}

	for result := range resultCh {

		if result.err != nil {
			// return early. done channel is closed, thus input channel is also closed.
			// all worker goroutines stop working (because input channel is closed)
			return nil, result.err

		}
		results = append(results, result.bodyStr)
	}

	return results, nil
}


func main() {
	// populate users param
	users := []string{}
	for i := 1; i <= 100; i++ {
		users = append(users, strconv.Itoa(i))
	}

	fmt.Println("Please Press Enter To Continue")
	fmt.Scanln()
	// Print User Appended Before goroutine Launnch
	start := time.Now()

	results, err := AsyncHTTP(users)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, result := range results {
		fmt.Println(result)
	}

	fmt.Println("finished in ", time.Since(start))
}
