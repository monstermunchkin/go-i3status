package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var pattern = regexp.MustCompile(`^(MemFree|Buffers|Cached): + (\d+) kB`)

func gradient(lo, hi, memory uint) func(val uint) string {
	var (
		mid         float32 = (float32(hi - lo)) / 2
		scale       float32 = 255 / (mid - float32(lo))
		step                = memory / (hi - lo + 1)
		colors              = make([]string, hi-lo+1)
		breakpoints         = make([]uint, hi-lo)
	)

	color := func(val uint) string {
		switch {
		case val <= lo:
			// red
			return "#FF0000"
		case val >= hi:
			// green
			return "#00FF00"
		case float32(val) < mid:
			// red-ish
			return fmt.Sprintf("#FF%02X00", uint(float32(val-lo)*scale))
		default:
			// green-ish
			return fmt.Sprintf("#%02XFF00", uint(255-(float32(val)-mid)*scale))
		}
	}

	for i, k := step, 0; i < memory; i, k = i+step, k+1 {
		breakpoints[k] = i
	}

	for i, k := lo, 0; i <= hi; i, k = i+1, k+1 {
		colors[k] = color(i)
	}

	out := func(val uint) string {
		if val <= breakpoints[0] {
			return colors[0]
		} else if val > breakpoints[len(breakpoints)-1] {
			return colors[len(colors)-1]
		}

		for i := 0; i < len(breakpoints); i++ {
			if val > breakpoints[i] && val <= breakpoints[i+1] {
				return colors[i+1]
			}
		}

		panic("gradient")
	}

	return out
}

func getFreeMemory() uint {
	var freeMemory uint = 0

	file, err := os.Open("/proc/meminfo")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if pattern.MatchString(line) {
			memory := (pattern.FindAllStringSubmatch(line, -1)[0][2])
			tmp, err := strconv.ParseUint(memory, 10, 32)
			if err != nil {
				log.Fatal(err)
			}
			freeMemory += uint(tmp)
		}
	}

	// Return free memory in MB.
	return freeMemory / 1024
}

func main() {
	type Element struct {
		FullText            string `json:"full_text"`
		ShortText           string `json:"short_text,omitempty"`
		Color               string `json:"color,omitempty"`
		MinWidth            string `json:"min_width,omitempty"`
		Align               string `json:"align,omitempty"`
		Name                string `json:"name"`
		Instance            string `json:"instance,omitempty"`
		Urgent              string `json:"urgent,omitempty"`
		Separator           string `json:"separator,omitempty"`
		SeparatorBlockWidth string `json:"separator_block_width,omitempty"`
	}

	scanner := bufio.NewScanner(os.Stdin)
	var (
		buffer   bytes.Buffer
		status   []Element
		colorize = gradient(0, 3, 1024)
	)

	// Skip the first line which contains the version header.
	if scanner.Scan() {
		if _, err := buffer.WriteString(scanner.Text()); err != nil {
			log.Fatal(err)
		}
	}

	// The second line contains the start of the infinite array.
	if scanner.Scan() {
		if _, err := buffer.WriteString(scanner.Text()); err != nil {
			log.Fatal(err)
		}
	}

	for scanner.Scan() {
		line, prefix := scanner.Text(), ""
		// Remember the comma at the beginning of the line.
		if strings.HasPrefix(line, ",") {
			line, prefix = line[1:], ","
		}

		if _, err := buffer.WriteString(prefix); err != nil {
			log.Fatal(err)
		}

		freeMemory := getFreeMemory()

		memoryElement := Element{
			FullText: fmt.Sprintf("Mem: %d MB", freeMemory),
			Name:     "mem_info",
			Color:    colorize(freeMemory),
		}

		dec := json.NewDecoder(strings.NewReader(line))
		if err := dec.Decode(&status); err != nil {
			log.Fatal(err)
		}

		// Prepend the memory element.
		status = append([]Element{memoryElement}, status...)

		out, err := json.Marshal(status)
		if err != nil {
			log.Fatal(err)
		}

		buffer.Write(out)
		buffer.WriteTo(os.Stdout)

		// Reset output.
		status = []Element{}
	}
}
