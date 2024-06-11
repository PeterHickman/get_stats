package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type callStats struct {
	size    int
	status  int
	elapsed float64
}

type allStats struct {
	count     int
	sizeTotal int
	sizeMin   int
	sizeMax   int
	timeTotal float64
	timeMin   float64
	timeMax   float64
	codes     map[int]int
}

func (a *allStats) Update(size int, elapsed float64, code int) {
	if a.count == 0 {
        // A hack so that we don't have to initialise this seperatly
	    a.codes = make(map[int]int)
		a.sizeMin = size
		a.sizeMax = size
		a.timeMin = elapsed
		a.timeMax = elapsed
	} else {
		if a.sizeMin > size {
			a.sizeMin = size
		}
		if a.sizeMax < size {
			a.sizeMax = size
		}
		if a.timeMin > elapsed {
			a.timeMin = elapsed
		}
		if a.timeMax < elapsed {
			a.timeMax = elapsed
		}
	}
	a.count++
	a.sizeTotal += size
	a.timeTotal += elapsed

	a.codes[code]++
}

func (a *allStats) Print() {
	fmt.Printf("Minimum response size was %d\n", a.sizeMin)
	fmt.Printf("Average response size was %d\n", a.sizeTotal/a.count)
	fmt.Printf("Maximum response size was %d\n", a.sizeMax)

	fmt.Printf("Minimum response time was %f\n", a.timeMin)
	fmt.Printf("Average response time was %f\n", a.timeTotal/float64(a.count))
	fmt.Printf("Maximum response time was %f\n", a.timeMax)

	for k, v := range a.codes {
		fmt.Printf("Status %d occurred %d times\n", k, v)
	}
}

func callUrl(url string) (callStats, error) {
	start := time.Now()
	resp, err := http.Get(url)
	t := time.Now()

	// Duration is in nanoseconds
	elapsed := float64(t.Sub(start)) / 1000000000

	if err != nil {
		return callStats{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return callStats{}, err
	}

	return callStats{len(body), resp.StatusCode, elapsed}, nil
}

func main() {
	urlPrefix := flag.String("prefix", "", "Prefix to be added to the input strings")
	urlSuffix := flag.String("suffix", "", "Suffix to be added to the input strings")
	dropdead := flag.Bool("dropdead", false, "Bail of the first 5xx status code")
	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatalln("No file of urls supplied")
	}

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	allStatii := new(allStats)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		url := *urlPrefix + scanner.Text() + *urlSuffix

		s, err := callUrl(url)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("Size = %d, Status = %d, Elapsed = %.4f %s\n", s.size, s.status, s.elapsed, url)

		allStatii.Update(s.size, s.elapsed, s.status)

		if *dropdead && s.status >= 500 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	allStatii.Print()
}
