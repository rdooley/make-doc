make-doc
============

Better makefile doc printing.

### Usage

```
make-doc --help
usage: make-doc [<flags>] <makefiles>...

Flags:
  --help       Show context-sensitive help (also try --help-long and
               --help-man).
  --variables  optional flag to parse variables
  --version    Show application version.

Args:
  <makefiles>  paths to makefiles
```

### Example

```
make-doc Makefile
Usage: make [target] [VARIABLE=value]

build:
	clean         	Clean stuff.
	build         	Build the actual thing.
	install       	Install make-doc to /usr/local/bin.

help:
	help          	show this help..
	help-variables	show makefile customizable variables..

make-doc Makefile --variables
build:
	PROJECT_NAME	Project name. Default : make-doc
```
