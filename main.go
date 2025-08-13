package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	var fn string
	flag.StringVar(&fn, "f", "", "firmware file")
	flag.Parse()

	if fn == "" {
		fmt.Println("need -f file")
		os.Exit(1)
	}

	f, err := os.Open(fn)
	if err != nil {
		fmt.Println("open:", err)
		os.Exit(1)
	}
	defer f.Close()

	st, _ := f.Stat()
	fmt.Printf("File: %s (%d bytes)\n", fn, st.Size())

	buf := make([]byte, 16)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		fmt.Println("read:", err)
		os.Exit(1)
	}

	fmt.Printf("First %d bytes: %x\n", n, buf[:n])
}
