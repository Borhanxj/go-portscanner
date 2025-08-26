package scanner

import (
	"context"
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

// ScanRange scans [start,end] on target with 'workers' goroutines and per-conn dial timeout.
// Returns a slice of open ports in ascending order.
func ScanRangeSingleWorker(
	ctx context.Context,
	target string,
	start int,
	end int,
	timeout time.Duration,
) ([]int, error) {

	open := make([]int, 0, end-start+1)

	for port := start; port < end; port++ {
		select {
		case <-ctx.Done():
			return open, ctx.Err()
		default:
		}

		// This is the standard format Go expects for a network address
		addr := net.JoinHostPort(target, fmt.Sprintf("%d", port))

		// Tries to open a TCP connection to the address we just built.
		conn, err := net.DialTimeout("tcp", addr, timeout)

		if err == nil {
			// Close the socket
			conn.Close()
			open = append(open, port)
		}
	}

	return open, nil
}

// dialOnce tries one TCP connect with timeout; returns true if port is open.
func dialOnce(target string, port int, timeout time.Duration) bool {
	addr := net.JoinHostPort(target, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err == nil {
		_ = conn.Close()
		return true
	}
	return false
}

// worker reads ports from jobs, tests them, and sends open ports to results.
func Worker(
	ctx context.Context,
	target string,
	timeout time.Duration,
	// 	jobs <-chan int → only receives port numbers.
	// results chan<- int → only sends open ports.
	jobs <-chan int,
	results chan<- int,
	wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case port, ok := <-jobs: // Jobs contains ports to scan
			if !ok { // If jobs channel is closed, finish the scan
				return
			}
			if dialOnce(target, port, timeout) { // If port is open, Add it to the results
				results <- port
			}
		}
	}
}

// Scan range with worker pool and go routines
func ScanRange(
	ctx context.Context,
	target string,
	start int,
	end int,
	workers int,
	timeout time.Duration,
) ([]int, error) {

	jobs := make(chan int, workers)
	results := make(chan int, end-start+1)
	var wg sync.WaitGroup

	// Start workers
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go Worker(ctx, target, timeout, jobs, results, &wg)
	}

	// Close results when all finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Feed jobs
	go func() {
		for port := start; port < end; port++ {
			select {
			case <-ctx.Done():
				close(jobs)
				return
			case jobs <- port:
			}
		}
		close(jobs)
	}()

	// Collect results
	open := make([]int, 0, end-start+1)
	for p := range results {
		open = append(open, p)
	}
	sort.Ints(open)
	return open, ctx.Err()
}
