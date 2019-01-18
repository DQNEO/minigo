/**
 * Execute:  gcc hello-world.c && ./a.out
 * Assemble: gcc -S -fno-asynchronous-unwind-tables hello-world.c -o hello-world.s
 */

#include <stdio.h>

int main() {
    puts("hello world");
}
