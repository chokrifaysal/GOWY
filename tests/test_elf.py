#!/usr/bin/env python3
import struct
import sys

def gen_elf(out):
    # minimal 32-bit ELF
    hdr = struct.pack("<16B", 0x7f, 0x45, 0x4c, 0x46, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0)  # e_ident
    hdr += struct.pack("<HHIIIIIHHHHHH", 2, 3, 1, 0x08048000, 0x34, 0, 0, 0x34, 0, 0, 0x1000, 0x1000)
    
    # .text section
    sh = struct.pack("<IIIIIIIIII", 0x08048000, 0x1000, 1, 6, 0, 0, 0, 0, 0x10, 0)
    
    data = b"\x90\x90\x90\x90" * 64  # NOP sled
    
    with open(out, 'wb') as f:
        f.write(hdr)
        f.write(data)
        f.write(sh)
    
    print(f"wrote {out}")

if __name__ == "__main__":
    fn = sys.argv[1] if len(sys.argv) > 1 else "test.elf"
    gen_elf(fn)
