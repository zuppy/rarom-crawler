package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

var activeDetailWorkers uint32 = 0

func detailWorker(currentWorkerNumber uint16, job <- chan serviceItem, jobWriter chan <- serviceItem, jobDone chan <- bool) {
	atomic.AddUint32(&activeDetailWorkers, 1);
	fmt.Printf("[detail worker %d] start\n", currentWorkerNumber)

	client := http.Client{Timeout: httpClientTimeout}

	for {
		serviceData, more := <-job

		// channel has been closed
		if !more {
			break;
		}

		url := baseUrl + "?jud=" + serviceData.county + "&cls=&act=&cat=&id=" + serviceData.id

		content, err := getContent(url, httpRetryCount, httpRetryWait, client)
		if err != nil {
			fmt.Printf("[detail worker %d] fatal error: %s\n", currentWorkerNumber, err)
			os.Exit(1)
		}

		err = processDetailItemList(content, serviceData, jobWriter)
		if err != nil {
			fmt.Printf("[detail worker %d] fatal error: %s\n\n%s\n",
				currentWorkerNumber, err, content)
			os.Exit(1)
		}

		time.Sleep(httpNextDelay)
	}

	fmt.Printf("[detail worker %d] end\n", currentWorkerNumber)

	// marking job as finished, if this is the last worker
	atomic.AddUint32(&activeDetailWorkers, ^uint32(0));
	if atomic.LoadUint32(&activeDetailWorkers) == 0 {
		close(jobWriter)
	}
}

func processDetailItemList(html []byte, serviceData serviceItem, jobWriter chan <- serviceItem) (err error) {
	// i know i could use regexps, but it's faster this way for my case (even if it's hackish)
	markerStart := []byte("Nume societate")
	markerEnd   := []byte("height: 1px; ")
	match, newCursorPosition := findTextBetweenMarkers(html, markerStart, markerEnd)
	if newCursorPosition < 1 {
		return errors.New("Could not parse detail block");
	}

	// this only contains the details
	html = match

	// name
	matchName, groupLength := findTextBetweenMarkers(html, []byte("100%\">"), []byte("</TD>"))
	if groupLength < 1 {
		return errors.New("Could not get name");
	}
	serviceData.name = string(matchName)

	// address: legal
	matchAddressLegal, groupLength := findTextBetweenMarkers(html, []byte("social:</TD><TD>"), []byte("</TD>"))
	if groupLength < 1 {
		return errors.New("Could not get legal address");
	}
	serviceData.addressLegal = string(matchAddressLegal)

	// address
	matchAddress, groupLength := findTextBetweenMarkers(html, []byte("lucru:</TD><TD>"), []byte("</TD>"))
	if groupLength < 1 {
		return errors.New("Could not get address");
	}
	serviceData.addressWork = string(matchAddress)

	// phone
	matchPhone, groupLength := findTextBetweenMarkers(html, []byte("Telefon:</TD><TD>"), []byte("</TD>"))
	if groupLength < 1 {
		return errors.New("Could not get phone");
	}
	serviceData.phone = string(matchPhone)

	jobWriter <- serviceData

	return nil
}