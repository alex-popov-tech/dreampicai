package main

import (
	"bufio"
	"bytes"
	"log"
	"log/slog"
	"os"
	"strconv"
)

// type NoHashingMap struct {
// 	size int
// }
// // how to convert string to bucket index with minimal collisions? Options:
// // 1. grab first letter and convert to unicode number, last letter to unicode number,
// func (it NoHashingMap) makeBucketIndex(key string) int64 {
// }
//
// func (it NoHashingMap) Set(key string, val int64) {
// }
//
// func (it NoHashingMap) Get(key string) int64 {
// }

// go run cmd/scanner/scanner.go  81.36s user 4.22s system 99% cpu 1:25.93 total
func main() {
	file, _ := os.Open("./testfile")
	defer file.Close()

	// try pre-cache approx cities count
	type City struct {
		Sum   int64
		Min   int64
		Max   int64
		Count int64
	}
	cities := make(map[string]*City, 1000)

	scanner := bufio.NewScanner(file)
	// read till \n char
	for scanner.Scan() {
		lineBytes := scanner.Bytes()
		// instead of cut, we can slice last 2 chars?
		cityBytes, temperatureBytes, found := bytes.Cut(lineBytes, []byte{59}) // []byte{59} == ";"
		if !found {
			log.Fatal("Invalid line format " + string(lineBytes))
		}
		// better parseint?
		// use unicode table to calculate digits from position of them in unicode?
		temp, err := strconv.ParseInt(string(temperatureBytes), 10, 8)
		if !found {
			log.Fatal("Cannot parse float", err)
		}
		c, ok := cities[string(cityBytes)]
		if !ok {
			cities[string(cityBytes)] = &City{Count: 1, Sum: temp, Min: temp, Max: temp}
			continue
		}
		c.Count += 1
		c.Sum += temp
		c.Min = min(c.Min, temp)
		c.Max = min(c.Max, temp)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for _, v := range cities {
		slog.Any("City", v)
	}
}
