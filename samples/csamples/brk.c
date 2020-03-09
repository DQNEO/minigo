// Example of using brk syscall and according malloc(3)
// Watch real syscalls as follows:
//   gcc -O0 malloc.c && strace ./a.out
#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/syscall.h>

int main() {
    void *p1;
    void *p2;

    int i;
    open("/==== START ===", O_RDONLY);

    // see https://github.com/bminor/glibc/blob/7975f9a48a83b95174503bda6c48124f08cb4f62/sysdeps/unix/sysv/linux/x86_64/brk.c#L31
    long adr1 = syscall(SYS_brk, 0); // get the current
    long adr2 = syscall(SYS_brk, adr1 + 4096);

    // Do similar things by malloc(3)
    p1 = (void *)malloc(0);
    p2 = (void *)malloc(4096);
    open("/==== END t===", O_RDONLY);

    printf("adr1=%p\n", (void *)adr1);
    printf("adr2=%p\n", (void *)adr2);
    printf("p1  =%p\n", p1);
    printf("p2  =%p\n", p2);
}
