package main

import (
	"debug/elf"
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

	if elf.IsELF(fn) {
		parseELF(fn)
	} else {
		buf := load(fn)
		dump(buf, sz)
	}
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

func parseELF(fn string) {
	f, err := elf.Open(fn)
	if err != nil {
		fmt.Println("elf open:", err)
		return
	}
	defer f.Close()

	fmt.Printf("ELF: %s\n", fn)
	fmt.Printf("Class: %s\n", f.FileHeader.Class)
	fmt.Printf("Machine: %s\n", f.FileHeader.Machine)

	for _, sec := range f.Sections {
		if sec.Type == elf.SHT_PROGBITS && sec.Size > 0 {
			fmt.Printf("\nSection: %s (0x%x bytes)\n", sec.Name, sec.Size)
			buf, _ := sec.Data()
			dump(buf, 256)
		}
	}
}
