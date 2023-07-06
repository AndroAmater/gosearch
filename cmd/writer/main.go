package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func openPipe() *os.File {
	fifoFile := "/tmp/fifo"

	f, err := os.OpenFile(fifoFile, os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	return f
}

func main() {
	fmt.Println("Open reader...")

	f := openPipe()
	defer f.Close()

	for i := 0; i < 3; i++ {
		fmt.Println("Write:")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		_, err := f.WriteString(fmt.Sprint(input.Text(), i, "\n"))
		if err != nil {
			panic(err)
		}
	}

	_, err := f.WriteString(fmt.Sprint(io.EOF))
	if err != nil {
		panic(err)
	}
}
