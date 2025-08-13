#!/usr/bin/env python3
import sys

def gen_str(out):
    with open(out, 'wb') as f:
        f.write(b'\x00' * 128)
        f.write(b'HelloWorld\x00')
        f.write(b'ERROR: invalid param\x00')
        f.write(b'fw_version=1.2.3\x00')
        f.write(b'\x00' * 64)
        f.write(b'secret_key_12345\x00')
        f.write(b'\x00' * 256)
    print(f"wrote {out}")

if __name__ == "__main__":
    fn = sys.argv[1] if len(sys.argv) > 1 else "test_str.bin"
    gen_str(fn)
