package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	// "github.com/pkg/profile"
)

const semicolonByte = 59

var semicolonBytes = []byte{59}
var linebreakByte = "\n"[0]

// o run cmd/scanner/scanner.go  54.51s user 3.14s system 99% cpu 57.932 total
func main() {
	// p := profile.Start(
	// 	profile.CPUProfile,
	// 	// profile.MemProfile,
	// 	profile.ProfilePath("."),
	// 	profile.NoShutdownHook,
	// )
	// defer p.Stop()
	start := time.Now()
	file, _ := os.Open("./testfile")
	afterOpen := time.Now()
	fmt.Println("Time spent opening file is seconds -", afterOpen.Sub(start).Seconds())
	defer file.Close()

	type City struct {
		Sum   uint
		Count uint
		Min   uint
		Max   uint
	}
	cities := make(map[string]*City, 1000)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineBytes := scanner.Bytes()
		cityBytes, temperatureBytes := separateBySemi(lineBytes)

		temp, _ := strconv.ParseUint(string(temperatureBytes), 10, 8)
		c := cities[string(cityBytes)]
		if c == nil {
			cities[string(cityBytes)] = &City{
				Count: 1,
				Sum:   uint(temp),
				Min:   uint(temp),
				Max:   uint(temp),
			}
		} else {
			c.Count += 1
			c.Sum += uint(temp)
			c.Min = min(c.Min, uint(temp))
			c.Max = max(c.Max, uint(temp))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for k, v := range cities {
		fmt.Printf("%s, sum=%d, min=%d, max=%d, count=%d\n", k, v.Sum, v.Min, v.Max, v.Count)
	}
}

func separateBySemi(row []byte) (cityName []byte, temperature []byte) {
	for i := len(row) - 1; i > 0; i-- {
		if row[i] == semicolonByte {
			return row[:i], row[i+1:]
		}
	}
	panic("separateBySemi failed on " + string(row))
}
