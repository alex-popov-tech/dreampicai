package main

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"syscall"
	// "github.com/pkg/profile"
)

const mb = 1024 * 1024

type City struct {
	Name  string
	Sum   int
	Count int
	Min   int
	Max   int
}

type Cities map[string]*City

func NewCities(capacity int) *Cities {
	cities := make(Cities, capacity)
	return &cities
}

func (it *Cities) SetCity(name string, sum, count, minimum, maximum int) {
	city := (*it)[name]
	if city == nil {
		city = &City{
			Name:  name,
			Sum:   sum,
			Count: count,
			Min:   minimum,
			Max:   maximum,
		}
		(*it)[name] = city
		return
	}

	city.Sum += sum
	city.Count += count
	city.Min = min(minimum, city.Min)
	city.Max = min(maximum, city.Max)
}

// go run cmd/scanner/scanner.go  109.59s user 5.29s system 857% cpu 13.401 total
func main() {
	// p := profile.Start(
	// 	profile.CPUProfile,
	// 	// profile.MemProfile,
	// 	profile.ProfilePath("."),
	// 	profile.NoShutdownHook,
	// )
	// defer p.Stop()

	f, err := os.OpenFile("./testfile", os.O_RDONLY, 0)
	// f, err := os.OpenFile("./testfile_small", os.O_RDONLY, 0)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		panic(err.Error())
	}
	fileSize := stat.Size()

	pageSize := int64(mb * 8)

	tmpRows := make([]byte, 1024)

	group := sync.WaitGroup{}

	cities := NewCities(1000)
	chunkCitiesChan := make(chan *Cities, 100)
	go func() {
		for chunk := range chunkCitiesChan {
			for k, v := range *chunk {
				cities.SetCity(k, v.Sum, v.Count, v.Min, v.Max)
			}
		}
	}()

	for index, offset := 0, int64(0); offset < fileSize; index, offset = index+1, offset+pageSize {
		mmap := Mmap(f, offset, pageSize)

		// cut first line or whatever rest of line from previous chunk
		newlineIndexFromStart := bytes.IndexByte(mmap, '\n')
		tmpRows = append([]byte(nil), append(tmpRows, mmap[:newlineIndexFromStart+1]...)...)
		// tmpRows = append(tmpRows, mmap[:newlineIndexFromStart+1]...)
		// processChunk(tmpRows, cities)

		// cut whatever rest incomplete line and save for next iteration
		newlineIndexFromEnd := bytes.LastIndexByte(mmap, '\n')
		// ended on newline, just clean slice
		if newlineIndexFromEnd == len(mmap)-1 {
			tmpRows = tmpRows[:0]
		}
		tmpRows = append([]byte(nil), mmap[newlineIndexFromEnd+1:]...)

		// process main chunk
		group.Add(1)
		go func() {
			cities := NewCities(1000)
			defer group.Done()
			chunk := mmap[newlineIndexFromStart+1 : newlineIndexFromEnd+1]
			processChunk(chunk, cities)
			chunkCitiesChan <- cities
			mmap.Unmap()
		}()
	}

	group.Wait()
	// all processing is done, close chan and process remaining chunks
	close(chunkCitiesChan)
	processChunk(tmpRows, cities)

	for _, v := range *cities {
		fmt.Println(v.Name, "Sum:", v.Sum, "Count:", v.Count, "Min:", v.Min, "Max:", v.Max)
	}
}

type MMap []byte

func Mmap(file *os.File, offset, pageSize int64) MMap {
	if stats, _ := file.Stat(); stats.Size() < offset+pageSize {
		pageSize = stats.Size() - offset
	}
	data, err := syscall.Mmap(
		int(file.Fd()),
		offset,
		int(pageSize),
		syscall.PROT_READ,
		syscall.MAP_SHARED,
	)
	if err != nil {
		panic("Cannot create mmap: " + err.Error())
	}

	return MMap(data)
}

func (m MMap) Unmap() error {
	return syscall.Munmap(m)
}

func processChunk(chunk []byte, cities *Cities) {
	var pointer int
	for pointer < len(chunk) {

		start, semi, newline := nextline(chunk, pointer)
		// appendFile(chunk[start:newline+1], "./log")
		city := chunk[start:semi]
		temperature := chunk[semi+1 : newline]
		pointer = newline + 1

		temp := fastParseUint(string(temperature))
		cities.SetCity(string(city), temp, 1, temp, temp)
	}
}

func nextline(data []byte, pointer int) (startIndex, semiIndex, newlineIndex int) {
	startIndex = pointer
	for i := startIndex; i < len(data); i++ {
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
	if len(s) == 1 {
		ch := s[0]
		if ch < '0' || ch > '9' {
			return 0
		}
		return int(ch - '0')
	}

	return int((s[0]-'0')*10 + (s[1] - '0'))
}
