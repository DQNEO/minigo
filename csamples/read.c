/**
 * Execute:  gcc read.c && ./a.out
 * Assemble: gcc -S -fno-asynchronous-unwind-tables read.c -o read.s
 */

#include <stdio.h>
#include <fcntl.h>
#include <sys/types.h>
#include <sys/uio.h>
#include <unistd.h>

#define MYBUFSIZ 1024 * 1024

int main() {
    int fd;
    char buf[MYBUFSIZ];

    fd = open("/etc/hosts", O_RDONLY);
    read(fd, buf, MYBUFSIZ);
    printf("%s", buf);
}
