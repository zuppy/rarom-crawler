package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

var activeCountyWorkers uint32 = 0

func countyWorker(currentWorkerNumber uint16, job <- chan string, jobOut chan <- serviceItem) {
	atomic.AddUint32(&activeCountyWorkers, 1);
	fmt.Printf("[county worker %d] start\n", currentWorkerNumber)

	client := http.Client{Timeout: httpClientTimeout}

	for {
		url, more := <-job

		// channel has been closed
		if !more {
			break;
		}

		content, err := getContent(url, httpRetryCount, httpRetryWait, client)
		if err != nil {
			fmt.Printf("[county worker %d] fatal error: %s\n", currentWorkerNumber, err)
			os.Exit(1)
		}

		err = processServiceItemList(content, jobOut)
		if err != nil {
			fmt.Printf("[county worker %d] fatal error: %s\n\n%s\n",
				currentWorkerNumber, err, content)
		}
	}

	fmt.Printf("[county worker %d] end\n", currentWorkerNumber)

	// marking job as finished, if this is the last worker
	atomic.AddUint32(&activeCountyWorkers, ^uint32(0));
	if atomic.LoadUint32(&activeCountyWorkers) == 0 {
		close(jobOut)
	}
}

func processServiceItemList(html []byte, jobDetail chan <- serviceItem) (err error) {
	// i know i could use regexps, but it's faster this way for my case (even if it's hackish)
	markerStart := []byte("alege atelier service")
	markerEnd   := []byte("</SELECT>")
	match, newCursorPosition := findTextBetweenMarkers(html, markerStart, markerEnd)
	if newCursorPosition < 1 {
		return errors.New("Could not read service list in county file");
	}

	// this only contains the list of services
	html = match[16:]
	maxLength := len(html)

	var currentPos = 0
	for {
		var newItem = serviceItem{}

		// county
		matchCounty, groupLength := findTextBetweenMarkers(html[currentPos:], []byte("'?jud="), []byte("&"))
		if groupLength < 1 {
			break;
		}
		newItem.county = string(matchCounty)
		currentPos = currentPos + groupLength

		// id
		matchId, groupLength := findTextBetweenMarkers(html[currentPos:], []byte("&id="), []byte("'"))
		if groupLength < 1 {
			break;
		}
		newItem.id = string(matchId)
		currentPos = currentPos + groupLength

		if currentPos >= maxLength {
			break
		}

		jobDetail <- newItem
		time.Sleep(httpNextDelay)
	}

	return nil
}