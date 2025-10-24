# passGen

`passGen` is a minimal command-line password generator written in Go.  
It is designed to be simple, fast, and dependency-free.

## Features

- Full random password generation
- Custom character ranges
- Adjustable password length
- Clean string output
- Command-line interface (CLI)

## Installation

Clone the repository and build from source:

```bash
	git clone https://github.com/lnkssr/passGen.git
	cd passGen
	make install
```

## Build with docker (coming soon)

``` bash
	docker build -t passgen-builder .
	docker create --name passgen-tmp passgen-builder
	mkdir -p bin
	docker cp passgen-tmp:/passGen ./bin/passGen
	docker rm passgen-tmp
```
