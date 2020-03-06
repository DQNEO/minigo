# This Makefile can run on both of Linux and MacOS

BIN_FILES := $(patsubst %.c,%.bin,$(wildcard *.c))
AS_FILES := $(patsubst %.c,%.s,$(wildcard *.c))

all: $(BIN_FILES) $(AS_FILES)

%.bin: %.c
	gcc -O0 -o $@ $<

%.s: %.c
	gcc -O0 -S -fno-asynchronous-unwind-tables -o $@ $<

clean:
	rm -f *.s *.bin
