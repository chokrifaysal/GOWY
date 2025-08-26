#!/usr/bin/env python3
import sys

def gen_xor(out):
    # create file with XOR encoded data
    plain = b"GOWY XOR test data with readable text and some numbers 12345"
    key = 0x42
    xor_data = bytes([b ^ key for b in plain])

    # pad with some random-looking data
    padding = b"\x89\x12\x34\x56\x78\x9a\xbc\xde" * 8
    data = padding + xor_data + padding

    with open(out, 'wb') as f:
        f.write(data)

    print(f"wrote {out} (XOR key: 0x{key:02x})")
    print(f"encoded: {plain.decode()}")

if __name__ == "__main__":
    fn = sys.argv[1] if len(sys.argv) > 1 else "test_xor.bin"
    gen_xor(fn)
