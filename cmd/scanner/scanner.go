package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	// "github.com/pkg/profile"
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
	Sum   uint
	Count uint
	Min   uint
	Max   uint
}

// go run cmd/scanner/scanner.go  51.41s user 3.09s system 98% cpu 55.426 total
func main() {
	// p := profile.Start(
	// 	// profile.CPUProfile,
	// 	profile.MemProfile,
	// 	profile.ProfilePath("."),
	// 	profile.NoShutdownHook,
	// )
	// defer p.Stop()
	start := time.Now()
	mmap, err := Map("./testfile")
	// mmap, err := Map("./testfile_small")
	if err != nil {
		log.Fatal(err)
	}
	afterOpen := time.Now()
	fmt.Println("Time spent opening file is seconds -", afterOpen.Sub(start).Seconds())

	bytesReader := bytes.NewReader(mmap)
	scanner := bufio.NewScanner(bytesReader)
	defer mmap.Unmap()

	cities := make(map[string]*City, 5000)

	for scanner.Scan() {
		line := scanner.Bytes()

		semiIndex := bytes.IndexByte(line, ';')

		city := line[:semiIndex]
		temperature := line[semiIndex+1:]

		temp := fastParseUint(string(temperature))
		c := cities[string(city)]

		if c == nil {
			cities[string(city)] = &City{
				Name:  city,
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

	for _, v := range cities {
		fmt.Printf("%s, sum=%d, min=%d, max=%d, count=%d\n", v.Name, v.Sum, v.Min, v.Max, v.Count)
	}
}

func fastParseUint(s string) uint64 {
	if len(s) == 0 || len(s) > 2 {
		return 0
	}

	if len(s) == 1 {
		ch := s[0]
		if ch < '0' || ch > '9' {
			return 0
		}
		return uint64(ch - '0')
	}

	// Two-digit parsing
	if s[0] < '0' || s[0] > '7' || s[1] < '0' || s[1] > '9' {
		return 0
	}

	return uint64((s[0]-'0')*10 + (s[1] - '0'))
}

