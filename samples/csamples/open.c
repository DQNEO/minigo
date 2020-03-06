/**
 * Execute:  gcc open.c && ./a.out
 * Assemble: gcc -S -fno-asynchronous-unwind-tables open.c -o open.s
 */

#include <stdio.h>
#include <fcntl.h>
#include <unistd.h>

int f1() {
    int fd;
    fd = open("/etc/hosts", O_RDONLY);
    close(fd);
    return fd;
}

int f2() {
    int fd;
    fd = open("/etc/foobar", O_RDONLY);
    return fd;
}

int main() {
    int fd;
    fd = f1();
    printf("open /etc/hosts returns %d\n", fd);

    fd = f2();
    printf("open /etc/foobar returns %d\n", fd);
}
