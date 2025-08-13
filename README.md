# GOWY - firmware tool

Extract and analyze firmware images.

## usage
```bash
go run . -f firmware.bin
go run . -f firmware.elf
go run . -f firmware.hex -x
go run . -f fw.bin -arm
go run . -f fw.bin -o dump.bin
```

## test
```bash
cd tests
python3 gen.py
./gowy -f tests/test_arm.bin -arm
```


## build
```bash 
go build -o gowy
```
