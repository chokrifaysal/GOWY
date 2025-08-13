package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var fn string
	var sz int
	flag.StringVar(&fn, "f", "", "firmware file")
	flag.IntVar(&sz, "s", 256, "bytes to dump")
	flag.Parse()

	if fn == "" {
		fmt.Println("need -f file")
		os.Exit(1)
	}

	buf := load(fn)
	if len(buf) == 0 {
		os.Exit(1)
	}

	dump(buf, sz)
}

func load(fn string) []byte {
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println("open:", err)
		return nil
	}
	defer f.Close()

	st, _ := f.Stat()
	buf := make([]byte, st.Size())
	n, err := f.Read(buf)
	if err != nil {
		fmt.Println("read:", err)
		return nil
	}
	return buf[:n]
}

func dump(buf []byte, sz int) {
	if sz > len(buf) {
		sz = len(buf)
	}
	for i := 0; i < sz; i += 16 {
		end := i + 16
		if end > sz {
			end = sz
		}
		fmt.Printf("%08x: ", i)
		for j := i; j < end; j++ {
			fmt.Printf("%02x ", buf[j])
		}
		for j := end - i; j < 16; j++ {
			fmt.Printf("   ")
		}
		fmt.Print(" ")
		for j := i; j < end; j++ {
			c := buf[j]
			if c < 32 || c > 126 {
				c = '.'
			}
			fmt.Printf("%c", c)
		}
		fmt.Println()
	}
}
