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

func gradient(shades uint, max_value uint) func(val uint) string {
	var (
		colors              = make([]string, shades)
		breakpoints         = make([]uint, shades-1)
		step                = max_value / shades
		scale       float32 = 255 / (float32(max_value) / 2)
	)

	color := func(val uint) string {
		switch {
		case val < step:
			// red
			return "#FF0000"
		case val >= step*(shades-1):
			// green
			return "#00FF00"
		case val < max_value/2:
			// red-ish
			return fmt.Sprintf("#FF%02X00", uint(float32(val)*scale))
		default:
			// green-ish
			return fmt.Sprintf("#%02XFF00", uint(255-(float32(val)-float32(max_value/2-step))*scale))
		}
	}

	for i := uint(0); i < shades-1; i++ {
		breakpoints[i] = step * (i + 1)
	}

	for i := uint(0); i < shades; i++ {
		colors[i] = color(step * i)
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

		log.Println("failed to determine color")
		return ""
	}

	return out
}

func getFreeMemory() (freeMemory uint) {
	file, err := os.Open("/proc/meminfo")

	if err != nil {
		log.Println(err)
		freeMemory = 0
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if pattern.MatchString(line) {
			memory := (pattern.FindAllStringSubmatch(line, -1)[0][2])
			tmp, err := strconv.ParseUint(memory, 10, 32)

			if err != nil {
				log.Println(err)
				continue
			}

			freeMemory += uint(tmp)
		}
	}

	// Return free memory in MB.
	freeMemory /= 1024
	return
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

	var (
		buffer   bytes.Buffer
		status   []Element
		colorize = gradient(4, 1024)
	)

	scanner := bufio.NewScanner(os.Stdin)

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
			log.Println(err)
			continue
		}

		freeMemory := getFreeMemory()

		memoryElement := Element{
			FullText: fmt.Sprintf("Mem: %d MB", freeMemory),
			Name:     "mem_info",
			Color:    colorize(freeMemory),
		}

		dec := json.NewDecoder(strings.NewReader(line))

		if err := dec.Decode(&status); err != nil {
			log.Println(err)
			continue
		}

		// Prepend the memory element.
		status = append([]Element{memoryElement}, status...)

		out, err := json.Marshal(status)

		if err != nil {
			log.Println(err)
			continue
		}

		buffer.Write(out)
		buffer.WriteTo(os.Stdout)

		// Reset output.
		status = []Element{}
	}
}
