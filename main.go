package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-ping/ping"
)

type TestResults struct {
	HttpStatus  bool
	PingStatus  bool
	PingLatency int64
}

func main() {

	fmt.Println("Starting outage test...")

	for {
		runResults := TestResults{
			HttpStatus: httpUp(),
		}

		pinger, err := ping.NewPinger("8.8.8.8")
		if err != nil {
			panic(err)
		}
		pinger.Count = 3
		pinger.Timeout = time.Second * 10
		pinger.Run()
		stats := pinger.Statistics()

		if !runResults.HttpStatus {
			fmt.Printf("[%v] Outage detected http Status: %t\n", time.Now(), runResults.HttpStatus)
		}
		if stats.AvgRtt > time.Second {
			fmt.Printf("[%v] Long ping RTT detected:%v", time.Now(), stats.AvgRtt)
		}

		time.Sleep(time.Second * 30)
	}

}

func httpUp() bool {
	res, err := http.Get("https://icanhazip.com/")

	if err != nil {
		return false
	}

	return res.StatusCode == 200
}
