package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-ping/ping"
)

type CheckResult struct {
	Up      bool
	Status  string
	Latency time.Duration
}

var wg sync.WaitGroup

func main() {

	fmt.Println("Starting outage test...")

	for {
		ch := make(chan CheckResult)

		wg.Add(2)
		go httpUp(ch)
		go pingUp(ch)

		go func() {
			wg.Wait()
			close(ch)
		}()

		for res := range ch {
			if !res.Up {
				fmt.Printf("[%v] Outage detected status: %s. Latency:%v\n", time.Now().Format("2006-01-02 15:04:05"), res.Status, res.Latency)
				return
			}

			fmt.Printf("[%v] Successful CheckResult: %s. Latency:%v\n", time.Now().Format("2006-01-02 15:04:05"), res.Status, res.Latency)
		}

		fmt.Println("Sleep 5 seconds...")
		time.Sleep(time.Second * 5)
	}
}

func httpUp(c chan CheckResult) {
	defer wg.Done()

	res, err := http.Get("https://icanhazip.com/")

	result := CheckResult{}
	if err != nil {
		result.Up = false
	}

	result.Up = true
	result.Status = res.Status + " icanhazip.com"

	c <- result
}

func pingUp(c chan CheckResult) {
	defer wg.Done()

	checkResult := CheckResult{}
	pinger, err := ping.NewPinger("8.8.8.8")
	if err != nil {
		checkResult.Up = false
		checkResult.Status = err.Error()
		return
	}

	pinger.Count = 3
	pinger.Timeout = time.Second * 10
	pinger.Run()
	stats := pinger.Statistics()

	checkResult.Up = true
	checkResult.Status = "Ping recv"
	checkResult.Latency = stats.AvgRtt

	c <- checkResult
}
