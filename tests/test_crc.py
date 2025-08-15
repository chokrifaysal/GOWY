#!/usr/bin/env python3
import binascii
import sys

def gen_crc(out):
    # file with known CRC32
    data = b"GOWY firmware test data" + b"\x00" * 100
    with open(out, 'wb') as f:
        f.write(data)
    chk = binascii.crc32(data) & 0xffffffff
    print(f"wrote {out} (CRC32: 0x{chk:08x})")

if __name__ == "__main__":
    fn = sys.argv[1] if len(sys.argv) > 1 else "test_crc.bin"
    gen_crc(fn)
