#!/usr/bin/env python3
import struct
import sys

def gen_base(out):
    # fake ARM fw with vector table at offset 0x400
    with open(out, 'wb') as f:
        f.write(b'\x00' * 0x400)                # padding
        f.write(struct.pack('<I', 0x20002000))  # SP
        f.write(struct.pack('<I', 0x08000441))  # PC
        f.write(b'\x00' * 2048)                 # rest
    print(f"wrote {out}")

if __name__ == "__main__":
    fn = sys.argv[1] if len(sys.argv) > 1 else "test_base.bin"
    gen_base(fn)
