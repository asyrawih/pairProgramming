package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// placeholder url. Change it to your base url.
const BASE_URL = "https://jsonplaceholder.typicode.com/posts/"

// number of parallelism
// try using only two worker
const NUM_PARALLEL = 2000

// Stream inputs to input channel
func streamInputs(done <-chan struct{}, inputs []string) <-chan string {
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
func sendUser(user string) (string, error) {
	url := BASE_URL + user
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

// Wrapper for sendUser return value, used as result channel type
type result struct {
	bodyStr string
	err     error
}

func AsyncHTTP(users []string) ([]string, error) {
	done := make(chan struct{})
	defer close(done)

	inputCh := streamInputs(done, users)

	var wg sync.WaitGroup
	// bulk add goroutine counter at the start
	wg.Add(NUM_PARALLEL)

	resultCh := make(chan result)

	for i := 0; i < NUM_PARALLEL; i++ {
		// spawn N worker goroutines, each is consuming a shared input channel.
		go func() {
			for input := range inputCh {
				bodyStr, err := sendUser(input)
				resultCh <- result{bodyStr, err}
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
