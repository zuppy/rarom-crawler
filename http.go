package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func getContent(url string, retryCount uint16, retryWait time.Duration, client http.Client) ([]byte, error) {
	var currentRetryCount uint16;

	for currentRetryCount = 1; currentRetryCount <= retryCount; currentRetryCount ++ {
		if currentRetryCount > 1 {
			fmt.Printf("http => [ %s ] Retry wait %d ms\n", url, retryWait/1000000) // @TODO: proper ms to s conversion
			time.Sleep(retryWait)
		}
		fmt.Printf("http => [ %s ] try %d\n", url, currentRetryCount);

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Printf("http => [ %s ] error %s\n", url, err)
			continue
		}
		request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.3; Trident/7.0; rv:11.0) like Gecko")

		response, err := client.Do(request)
		if err != nil {
			fmt.Printf("http => [ %s ] error %s\n", url, err)
			continue
		}

		if response.StatusCode != 200 {
			fmt.Printf("http => [ %s ] invalid status code %d\n", url, response.StatusCode)
			continue
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			continue
		}

		return body, nil
	}
	return nil, errors.New("Failed retrieving url: " + url)
}
