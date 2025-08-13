#!/usr/bin/env python3
import struct
import sys

def gen_arm(out):
    # ARM Cortex-M vector table
    vectors = [
        0x20001000,  # SP
        0x08000041,  # PC (thumb mode)
        0x08000045,  # NMI
        0x08000049,  # HardFault
    ]
    
    with open(out, 'wb') as f:
        for v in vectors:
            f.write(struct.pack('<I', v))
        f.write(b'\x00' * 256)  # padding
    
    print(f"wrote {out}")

if __name__ == "__main__":
    fn = sys.argv[1] if len(sys.argv) > 1 else "test_arm.bin"
    gen_arm(fn)
