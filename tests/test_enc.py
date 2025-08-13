#!/usr/bin/env python3
import os
import sys

def gen_enc(out):
    # make high entropy data
    buf = os.urandom(1024)
    with open(out, 'wb') as f:
        f.write(b'\x00' * 512)  # low entropy
        f.write(buf)            # high entropy
        f.write(b'\x00' * 512)  # low entropy
    print(f"wrote {out}")

if __name__ == "__main__":
    fn = sys.argv[1] if len(sys.argv) > 1 else "test_enc.bin"
    gen_enc(fn)
