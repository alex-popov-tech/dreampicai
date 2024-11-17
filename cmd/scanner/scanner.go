package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

const semicolonByte = 59
	"github.com/pkg/profile"

var semicolonBytes = []byte{59}
var linebreakByte = "\n"[0]

// go run cmd/scanner/scanner.go  67.98s user 3.29s system 99% cpu 1:11.67 total
func main() {
	p := profile.Start(
		profile.CPUProfile,
		// profile.MemProfile,
		profile.ProfilePath("."),
		profile.NoShutdownHook,
	)
	defer p.Stop()
	start := time.Now()
	file, _ := os.Open("./testfile")
	afterOpen := time.Now()
	fmt.Println("Time spent opening file is seconds -", afterOpen.Sub(start).Seconds())
	defer file.Close()

	type City struct {
		Sum   uint16
		Count uint16
		Min   uint8
		Max   uint8
	}
	cities := make(map[string]*City, 1000)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineBytes := scanner.Bytes()
		cityBytes, temperatureBytes := separateBySemi(lineBytes)
		temp := parseNumberInStringAsBytes(temperatureBytes)
		c := cities[string(cityBytes)]
		if c == nil {
			cities[string(cityBytes)] = &City{
				Count: 1,
				Sum:   uint16(temp),
				Min:   uint8(temp),
				Max:   uint8(temp),
			}
		} else {
			c.Count += 1
			c.Sum += uint16(temp)
			c.Min = min(c.Min, uint8(temp))
			c.Max = max(c.Max, uint8(temp))
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

func parseNumberInStringAsBytes(input []byte) uint8 {
	switch string(input) {
	case "0":
		return 0
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	case "4":
		return 4
	case "5":
		return 5
	case "6":
		return 6
	case "7":
		return 7
	case "8":
		return 8
	case "9":
		return 9
	case "10":
		return 10
	case "11":
		return 11
	case "12":
		return 12
	case "13":
		return 13
	case "14":
		return 14
	case "15":
		return 15
	case "16":
		return 16
	case "17":
		return 17
	case "18":
		return 18
	case "19":
		return 19
	case "20":
		return 20
	case "21":
		return 21
	case "22":
		return 22
	case "23":
		return 23
	case "24":
		return 24
	case "25":
		return 25
	case "26":
		return 26
	case "27":
		return 27
	case "28":
		return 28
	case "29":
		return 29
	case "30":
		return 30
	case "31":
		return 31
	case "32":
		return 32
	case "33":
		return 33
	case "34":
		return 34
	case "35":
		return 35
	case "36":
		return 36
	case "37":
		return 37
	case "38":
		return 38
	case "39":
		return 39
	case "40":
		return 40
	case "41":
		return 41
	case "42":
		return 42
	case "43":
		return 43
	case "44":
		return 44
	case "45":
		return 45
	case "46":
		return 46
	case "47":
		return 47
	case "48":
		return 48
	case "49":
		return 49
	case "50":
		return 50
	case "51":
		return 51
	case "52":
		return 52
	case "53":
		return 53
	case "54":
		return 54
	case "55":
		return 55
	case "56":
		return 56
	case "57":
		return 57
	case "58":
		return 58
	case "59":
		return 59
	case "60":
		return 60
	case "61":
		return 61
	case "62":
		return 62
	case "63":
		return 63
	case "64":
		return 64
	case "65":
		return 65
	case "66":
		return 66
	case "67":
		return 67
	case "68":
		return 68
	case "69":
		return 69
	case "70":
		return 70
	default:
		return 0
	}
}
