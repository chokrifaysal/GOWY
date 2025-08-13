package main

import (
	"bufio"
	"debug/elf"
	"encoding/binary"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	var fn, out string
	var sz, min int
	var x, arm, ent, base, str bool
	flag.StringVar(&fn, "f", "", "firmware file")
	flag.StringVar(&out, "o", "", "output file")
	flag.IntVar(&sz, "s", 256, "bytes to dump")
	flag.IntVar(&min, "m", 4, "min string len")
	flag.BoolVar(&x, "x", false, "extract all")
	flag.BoolVar(&arm, "arm", false, "parse ARM vectors")
	flag.BoolVar(&ent, "e", false, "entropy scan")
	flag.BoolVar(&base, "b", false, "find base addr")
	flag.BoolVar(&str, "str", false, "dump strings")
	flag.Parse()

	if fn == "" {
		fmt.Println("need -f file")
		os.Exit(1)
	}

	if str {
		dumpStr(fn, min)
		return
	}
	if base {
		findBase(fn)
		return
	}
	if ent {
		scanEnt(fn)
		return
	}
	if arm {
		parseARM(fn)
		return
	}
	if x {
		if isELF(fn) {
			extrELF(fn)
		} else if strings.HasSuffix(fn, ".hex") {
			extrHEX(fn)
		} else {
			fmt.Println("raw extract not supported")
		}
		return
	}
	if isELF(fn) {
		parseELF(fn)
	} else if strings.HasSuffix(fn, ".hex") {
		parseHEX(fn, sz)
	} else {
		buf := load(fn)
		if out != "" {
			os.WriteFile(out, buf, 0644)
			fmt.Println("wrote", out)
		} else {
			dump(buf, sz)
		}
	}
}

func dumpStr(fn string, min int) {
	buf := load(fn)
	for off := 0; off < len(buf); {
		end := off
		for end < len(buf) && buf[end] >= 32 && buf[end] < 127 {
			end++
		}
		if end-off >= min {
			fmt.Printf("%08x: %s\n", off, string(buf[off:end]))
		}
		off = end + 1
	}
}

func findBase(fn string) {
	buf := load(fn)
	if len(buf) < 0x1000 {
		fmt.Println("too small for base scan")
		return
	}
	for off := 0; off < len(buf)-0x40; off += 0x100 {
		sp := binary.LittleEndian.Uint32(buf[off:])
		pc := binary.LittleEndian.Uint32(buf[off+4:])
		if sp>>20 == 0x200 && (pc&0xFFFFFFFE) > 0x08000000 && (pc&0xFFFFFFFE) < 0x08100000 {
			fmt.Printf("possible base: 0x%08x\n", 0x08000000+uint32(off))
			fmt.Printf("SP: 0x%08x PC: 0x%08x\n", sp, pc)
		}
	}
}

func scanEnt(fn string) {
	buf := load(fn)
	if len(buf) == 0 {
		return
	}
	win := 256
	for i := 0; i < len(buf)-win; i += win {
		e := ent(buf[i : i+win])
		if e > 7.5 {
			fmt.Printf("0x%08x-0x%08x entropy %.2f [PACKED]\n", i, i+win, e)
		}
	}
}

func ent(b []byte) float64 {
	cnt := [256]int{}
	for _, v := range b {
		cnt[v]++
	}
	h := 0.0
	for _, c := range cnt {
		if c == 0 {
			continue
		}
		p := float64(c) / float64(len(b))
		h -= p * math.Log2(p)
	}
	return h
}

func parseARM(fn string) {
	buf := load(fn)
	if len(buf) < 0x40 {
		fmt.Println("too small for ARM vectors")
		return
	}
	sp := binary.LittleEndian.Uint32(buf[0:4])
	pc := binary.LittleEndian.Uint32(buf[4:8])
	nmi := binary.LittleEndian.Uint32(buf[8:12])
	hard := binary.LittleEndian.Uint32(buf[12:16])
	fmt.Printf("ARM Cortex-M vectors:\n")
	fmt.Printf("SP:  0x%08x\n", sp)
	fmt.Printf("PC:  0x%08x\n", pc)
	fmt.Printf("NMI: 0x%08x\n", nmi)
	fmt.Printf("HARDFAULT: 0x%08x\n", hard)
	if pc&1 != 0 {
		fmt.Println("Thumb mode enabled")
	}
}

func isELF(fn string) bool {
	f, err := os.Open(fn)
	if err != nil {
		return false
	}
	defer f.Close()
	buf := make([]byte, 4)
	n, err := f.Read(buf)
	return err == nil && n == 4 && string(buf) == "\x7fELF"
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
	mem := loadHEX(f)
	fmt.Printf("HEX: %s (%d bytes)\n", fn, len(mem))
	dump(mem, sz)
}

func extrELF(fn string) {
	f, err := elf.Open(fn)
	if err != nil {
		fmt.Println("elf open:", err)
		return
	}
	defer f.Close()
	os.MkdirAll("extr", 0755)
	for _, sec := range f.Sections {
		if sec.Type == elf.SHT_PROGBITS && sec.Size > 0 {
			buf, _ := sec.Data()
			out := fmt.Sprintf("extr/%s.bin", sec.Name)
			os.WriteFile(out, buf, 0644)
			fmt.Println("wrote", out)
		}
	}
}

func extrHEX(fn string) {
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println("hex open:", err)
		return
	}
	defer f.Close()
	mem := loadHEX(f)
	os.WriteFile("extr/firmware.bin", mem, 0644)
	fmt.Println("wrote extr/firmware.bin")
}

func loadHEX(f *os.File) []byte {
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
	return mem
}
