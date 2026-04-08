package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

// getLinesChannel reads lines from the given file and sends them to a channel.
func getLinesChannel(f io.ReadCloser) <-chan string {

	// Create a channel to send lines read from the file. The channel is buffered with a capacity of 1 to allow for asynchronous reading.
	out := make(chan string, 1)

	// Start a goroutine to read lines from the file and send them to the channel.
	go func() {

		defer f.Close()
		defer close(out)

		// str is used to accumulate data until a newline is found.
		str := ""

		// Read data from the file in chunks and process it to extract lines.
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				break
			}

			// Truncate the data slice to the number of bytes read.
			data = data[:n]

			// Check if there is a newline character in the data. If found, extract the line and send it to the channel.
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				str += string(data[:i])
				data = data[i+1:]
				out <- str
				str = ""
			}

			// Append the remaining data to str for the next iteration.
			str += string(data)
		}

		// If there is any remaining data in str after the loop, send it to the channel.
		if len(str) != 0 {
			out <- str
		}

	}()

	return out
}

func main() {

	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("error", "error", err)
	}

	lines := getLinesChannel(f)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}

}
