package main

import (
	"fmt"
	// "github.com/pkg/profile"
	"log"
	"os"
	"syscall"
	"time"
)

type MMap []byte

func Map(path string) (MMap, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := fi.Size()

	data, err := syscall.Mmap(
		int(f.Fd()),
		0,
		int(size),
		syscall.PROT_READ,
		syscall.MAP_SHARED,
	)
	if err != nil {
		return nil, err
	}

	return MMap(data), nil
}

func (m MMap) Unmap() error {
	return syscall.Munmap(m)
}

type City struct {
	Name  []byte
	Sum   int
	Count int
	Min   int
	Max   int
}

// go run cmd/scanner/scanner.go  38.26s user 2.71s system 99% cpu 41.183 total
func main() {
	// p := profile.Start(
	// 	// profile.CPUProfile,
	// 	profile.MemProfile,
	// 	profile.ProfilePath("."),
	// 	profile.NoShutdownHook,
	// )
	// defer p.Stop()
	startTime := time.Now()
	mmap, err := Map("./testfile")
	// mmap, err := Map("./testfile_small")
	defer mmap.Unmap()
	if err != nil {
		log.Fatal(err)
	}
	afterOpen := time.Now()
	fmt.Println("Time spent opening file is seconds -", afterOpen.Sub(startTime).Seconds())

	cities := make(map[string]*City, 5000)

	var pointer int
	for pointer < len(mmap) {

		start, semi, newline := nextline(mmap, pointer)
		city := mmap[start:semi]
		temperature := mmap[semi+1 : newline]
		pointer = newline + 1

		temp := fastParseUint(string(temperature))
		c := cities[string(city)]

		if c == nil {
			cities[string(city)] = &City{
				Name:  city,
				Count: 1,
				Sum:   temp,
				Min:   temp,
				Max:   temp,
			}
		} else {
			c.Count += 1
			c.Sum += temp
			c.Min = min(c.Min, temp)
			c.Max = max(c.Max, temp)
		}
	}

	for _, v := range cities {
		fmt.Printf("%s, sum=%d, min=%d, max=%d, count=%d\n", v.Name, v.Sum, v.Min, v.Max, v.Count)
	}
}

func nextline(data []byte, pointer int) (startIndex, semiIndex, newlineIndex int) {
	startIndex = pointer
	for i := startIndex; ; i++ {
		if data[i] == ';' {
			semiIndex = i
			continue
		}
		if data[i] == '\n' {
			newlineIndex = i
			break
		}
	}

	return startIndex, semiIndex, newlineIndex
}

func fastParseUint(s string) int {
	if len(s) == 0 || len(s) > 2 {
		return 0
	}

	if len(s) == 1 {
		ch := s[0]
		if ch < '0' || ch > '9' {
			return 0
		}
		return int(ch - '0')
	}

	// Two-digit parsing
	if s[0] < '0' || s[0] > '7' || s[1] < '0' || s[1] > '9' {
		return 0
	}

	return int((s[0]-'0')*10 + (s[1] - '0'))
}
