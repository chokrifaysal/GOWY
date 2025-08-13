#!/usr/bin/env python3
import sys

def gen_hex(out):
    with open(out, 'w') as f:
        # :LLAAAATTDD...CC
        # test data at 0x1000
        f.write(":10000000000102030405060708090A0B0C0D0E0F10\n")
        f.write(":10001000101112131415161718191A1B1C1D1E1F20\n")
        f.write(":10002000202122232425262728292A2B2C2D2E2F30\n")
        f.write(":00000001FF\n")  # EOF

if __name__ == "__main__":
    fn = sys.argv[1] if len(sys.argv) > 1 else "test.hex"
    gen_hex(fn)
    print(f"wrote {fn}")
