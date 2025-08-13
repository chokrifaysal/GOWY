package main

import (
	"bufio"
	"debug/elf"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	} else if strings.HasSuffix(fn, ".hex") {
		parseHEX(fn, sz)
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

func parseHEX(fn string, sz int) {
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println("hex open:", err)
		return
	}
	defer f.Close()

	var mem []byte
	max := 0

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		l := sc.Text()
		if len(l) < 11 || l[0] != ':' {
			continue
		}

		ln, _ := strconv.ParseUint(l[1:3], 16, 8)
		if ln == 0 {
			continue
		}

		addr, _ := strconv.ParseUint(l[3:7], 16, 16)
		typ, _ := strconv.ParseUint(l[7:9], 16, 8)
		if typ != 0 {
			continue
		}

		end := 9 + int(ln*2)
		if end > len(l) {
			continue
		}

		if int(addr)+int(ln) > max {
			nmem := make([]byte, int(addr)+int(ln))
			copy(nmem, mem)
			mem = nmem
			max = int(addr) + int(ln)
		}

		for i := 0; i < int(ln); i++ {
			b, _ := strconv.ParseUint(l[9+i*2:11+i*2], 16, 8)
			mem[addr+uint64(i)] = byte(b)
		}
	}

	fmt.Printf("HEX: %s (%d bytes)\n", fn, len(mem))
	dump(mem, sz)
}
