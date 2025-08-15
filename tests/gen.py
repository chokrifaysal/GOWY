#!/usr/bin/env python3
import os
import subprocess

os.chdir(os.path.dirname(__file__))

for t in ["test_hex", "test_elf", "test_bin", "test_arm", "test_enc", "test_base", "test_str", "test_crc"]:
    subprocess.run(["python3", t + ".py"])
