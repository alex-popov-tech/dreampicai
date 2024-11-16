package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// go run cmd/scanner/scanner.go  225.80s user 6.31s system 100% cpu 3:50.96 total
func main() {
	file, _ := os.Open("./testfile")
	defer file.Close()

	// try pre-cache approx cities count
	sums := make(map[string]int64, 1000)
	minTemperatures := make(map[string]int64, 1000)
	maxTemperatures := make(map[string]int64, 1000)
	counts := make(map[string]int64, 1000)

	scanner := bufio.NewScanner(file)
	// read till \n char
	for scanner.Scan() {
		// conv to string, possible replace with bytes?
		line := scanner.Text()
		// instead of cut, we can slice last 2 chars?
		city, temp, found := strings.Cut(line, ";")
		if !found {
			log.Fatal("Invalid line format " + line)
		}
		// better parseint?
		// use unicode table to calculate digits from position of them in unicode?
		intTemp, err := strconv.ParseInt(temp, 10, 8)
		if !found {
			log.Fatal("Cannot parse float", err)
		}
		// 4 times we hash again and again
		counts[city]++
		sums[city] += intTemp
		minTemperatures[city] = min(minTemperatures[city], intTemp)
		maxTemperatures[city] = max(minTemperatures[city], intTemp)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for k, v := range sums {
		fmt.Printf("City - %s, min: %d, max: %d, count: %d, sum: %d\n", k, minTemperatures[k], maxTemperatures[k], counts[k], v)
	}

}
