#!/usr/bin/env python3
import sys

def gen_bin(out):
    buf = bytes(range(256)) * 4
    with open(out, 'wb') as f:
        f.write(buf)
    print(f"wrote {out}")

if __name__ == "__main__":
    fn = sys.argv[1] if len(sys.argv) > 1 else "test.bin"
    gen_bin(fn)
