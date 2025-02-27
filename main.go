package main

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/tursodatabase/go-libsql"
)

func main() {
	args := os.Args[1:]

	cfg, err := configure(args)
	if err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}
	defer cfg.db.Close()

	cfg.logger.Infof("FETCHING: Total Results")
	totalResults, err := cfg.getTotalResults()
	if err != nil {
		cfg.logger.Fatalf("Failed to retrieve total results: %v", err)
	}
	cfg.logger.Infof("STATS: Total Agents: %d\n", totalResults)

	allRequests, err := cfg.getRequests(totalResults)
	if err != nil {
		cfg.logger.Fatalf("Failed to retrieve search requests: %v", err)
	}
	cfg.addRequests(allRequests)

	for {
		remainingReqs, isComplete := cfg.getRemainingRequests()
		cfg.logger.Infof(
			"STATS: (Total Requests: %d, Remaining Requests: %d)", cfg.getRequestCount(), len(remainingReqs))
		if isComplete {
			cfg.logger.Infof("========== COMPLETE ==========")
			return
		}

		count := 0
		for _, reqKey := range remainingReqs {
			request, err := cfg.getRequest(reqKey)
			if err != nil {
				cfg.logger.Fatalf("Failed to retrieve (request: %d): %v", reqKey, err)
			}
			if request.processed {
				continue
			}

			cfg.wg.Add(1)
			go cfg.processRequest(request)

			count++
			if count%cfg.throttleRequestLimit == 0 {
				remaining, _ := cfg.getRemainingRequests()
				cfg.logger.Infof(
					"THROTTLING: Remaining Requests: %d. Request limit of %d reached. Pausing for %v.",
					len(remaining),
					cfg.throttleRequestLimit,
					defaultThrottleTimeout,
				)
				time.Sleep(defaultThrottleTimeout)
			}

			if cfg.platform == "dev" && count >= cfg.throttleRequestLimit {
				break
			}

		}
		cfg.wg.Wait()
		if cfg.platform == "dev" {
			cfg.logger.Infof("DONE")
			fmt.Println()
			fmt.Println("Results: ")
			fmt.Println("* Threads:", cfg.maxThreadCount)
			fmt.Println("* Requests:", cfg.throttleRequestLimit)
			fmt.Println("* Agent Scraped:", cfg.agentCount.Load())
			fmt.Println()
			return
		}
	}
}
