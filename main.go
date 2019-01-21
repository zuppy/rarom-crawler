package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"time"
)

func main() {
	var Args TypeArgs
	arg.MustParse(&Args)

	jobCounty := make(chan string, 50);
	jobDetail := make(chan serviceItem, 50);
	jobWriter := make(chan serviceItem, 250);
	jobDone   := make(chan bool);

	for workerNumber := uint16(1); workerNumber <= countyWorkers; workerNumber ++ {
		go countyWorker(workerNumber, jobCounty, jobDetail)
	}
	for workerNumber := uint16(1); workerNumber <= detailWorkers; workerNumber ++ {
		go detailWorker(workerNumber, jobDetail, jobWriter, jobDone)
	}
	go writer(jobWriter, jobDone, Args.OutputFile)

	time.Sleep(time.Millisecond * 500)

	for _, county := range counties {
		jobCounty <- baseUrl + "?jud=" + county
	}
	close(jobCounty)

	<-jobDone

	fmt.Println("Done")
}
