package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func writer(job <- chan serviceItem, jobDone chan <- bool, outputFile string) {
	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Could not write output file ", outputFile, "; error: ", err)
		os.Exit(1)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.WriteString(
		"Judet;Nume;Telefon;Sediu social;Punct de lucru\n")
	if err != nil {
		fmt.Println("Could not write into buffer.")
		os.Exit(1)
	}

	for {
		item, more := <-job

		// channel has been closed
		if !more {
			break;
		}

		_, err := w.WriteString(
			stupidSimpleHackishSanitizer(item.county) + ";" +
			stupidSimpleHackishSanitizer(item.name) + ";" +
			stupidSimpleHackishSanitizer(item.phone) + ";" +
			stupidSimpleHackishSanitizer(item.addressLegal) + ";" +
			stupidSimpleHackishSanitizer(item.addressWork) + "\n")
		if err != nil {
			fmt.Println("Could not write into buffer.")
			os.Exit(1)
		}
	}

	err = w.Flush()
	if err != nil {
		fmt.Println("Could not flush buffer.")
		os.Exit(1)
	}

	jobDone <- true
}

func stupidSimpleHackishSanitizer(in string) (out string)  {
	// @TODO: do proper escape for csv, not this...
	out = strings.Replace(in, ";", ",", -1)
	out = strings.Replace(in, "\n", " ", -1)
	out = strings.Replace(in, "\r", "", -1)
	return
}
