#!/usr/bin/env python3
import hashlib
import sys

def gen_hash(out):
    # file with known hashes
    data = b"GOWY hash test data" + b"\x00" * 77
    with open(out, 'wb') as f:
        f.write(data)

    md5h = hashlib.md5(data).hexdigest()
    sha1h = hashlib.sha1(data).hexdigest()
    sha256h = hashlib.sha256(data).hexdigest()

    print(f"wrote {out}")
    print(f"MD5:    {md5h}")
    print(f"SHA1:   {sha1h}")
    print(f"SHA256: {sha256h}")

if __name__ == "__main__":
    fn = sys.argv[1] if len(sys.argv) > 1 else "test_hash.bin"
    gen_hash(fn)
