# GOWY - firmware tool

Extract and analyze firmware images.

## test
```bash
cd tests
python3 gen.py
../gowy -f test.hex
../gowy -f test.elf -x
```

## usage
```bash
go run . -f firmware.bin
go run . -f firmware.bin -s 512
go run . -f firmware.elf
go run . -f firmware.hex -x
go run . -f fw.bin -o dump.bin
```
## build
```bash
go build -o gowy
```
