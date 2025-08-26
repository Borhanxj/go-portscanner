package main

import (
	"context"
	"flag"
	"github.com/Borhanxj/go-portscanner/internal/scanner"
	"github.com/fatih/color"
	"time"
)

func main() {

	// CLI flags
	target := flag.String("target", "", "IP or hostname to scan")
	start := flag.Int("start", 1, "Start port number")
	end := flag.Int("end", 1024, "End port number")
	workers := flag.Int("workers", 500, "Number of concurrent workers")
	timeout := flag.Duration("timeout", 500*time.Millisecond, "Per-connection timeout")

	flag.Parse()

	// Input check
	if *target == "" {
		color.Red("Error: target is required")
		return
	}

	if *start < 1 || *end > 65535 || *start > *end {
		color.Red("Invalid port range")
		return
	}

	if *workers <= 0 {
		color.Red("Workers must be > 0")
		return
	}

	color.Cyan("Scanning target: %s", *target)
	color.Green("Starting from port: %d", *start)
	color.Green("Ending at port: %d", *end)
	color.Yellow("Workers: %d", *workers)
	color.Magenta("Timeout: %s", timeout.String())

	// Scan
	ports, err := scanner.ScanRange(
		context.Background(),
		*target,
		*start, *end,
		*workers,
		*timeout,
	)

	if err != nil {
		color.Red("Error encountered: %v", err)
		return
	}

	if len(ports) == 0 {
		color.Yellow("No open ports found in %d-%d", *start, *end)
	} else {
		color.Green("Open ports (%d):", len(ports))
		for _, p := range ports {
			color.Cyan("%s:%d open", *target, p)
		}
	}
}
