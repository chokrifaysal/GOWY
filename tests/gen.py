#!/usr/bin/env python3
import os
import subprocess

os.chdir(os.path.dirname(__file__))

subprocess.run(["python3", "test_hex.py"])
subprocess.run(["python3", "test_elf.py"])
subprocess.run(["python3", "test_bin.py"])
subprocess.run(["python3", "test_arm.py"])
subprocess.run(["python3", "test_enc.py"])
subprocess.run(["python3", "test_base.py"])
