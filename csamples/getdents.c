// original: http://man7.org/linux/man-pages/man2/getdents64.2.html#top_of_page
#define _GNU_SOURCE
#include <dirent.h>     /* Defines DT_* constants */
#include <fcntl.h>
#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <sys/stat.h>
#include <sys/syscall.h>

#define handle_error(msg) \
        do { perror(msg); exit(EXIT_FAILURE); } while (0)

struct linux_dirent64 {
   ino64_t        d_ino;    /* 8 bytes: 64-bit inode number */
   off64_t        d_off;    /* 8 bytes: 64-bit offset to next structure */
   unsigned short d_reclen; /* 2 bytes: Size of this dirent  */
   unsigned char  d_type;   /* 1 byte: File type */
   char           d_name[]; /* Filename (null-terminated) */
};

#define BUF_SIZE 4096

int
main(int argc, char *argv[])
{
/*
    printf("unsigned short=%ld\n", sizeof(unsigned short));
    printf("int=%ld\n", sizeof(int));
    printf("long=%ld\n", sizeof(long));
    printf("off_t=%ld\n", sizeof(off_t));
    printf("size_t=%ld\n", sizeof(size_t));
*/

    int fd, nread;
    char buf[BUF_SIZE];
    struct linux_dirent64 *d;
    int bpos;
    char d_type;

    fd = open(argc > 1 ? argv[1] : ".", O_RDONLY | O_DIRECTORY);
    if (fd == -1)
        handle_error("open");

    for ( ; ; ) {
        nread = syscall(SYS_getdents64, fd, buf, BUF_SIZE);
        if (nread == -1)
            handle_error("getdents");

        if (nread == 0)
            break;

        printf("--------------- nread=%d ---------------\n", nread);
        printf("inode#    file type  d_reclen  d_off   d_name\n");
        for (bpos = 0; bpos < nread;) {
            d = (struct linux_dirent64 *) (buf + bpos);
            printf("%8ld  ", d->d_ino);
            d_type = d->d_type;
            printf("%-10s ", (d_type == DT_REG) ?  "regular" :
                             (d_type == DT_DIR) ?  "directory" :
                             (d_type == DT_FIFO) ? "FIFO" :
                             (d_type == DT_SOCK) ? "socket" :
                             (d_type == DT_LNK) ?  "symlink" :
                             (d_type == DT_BLK) ?  "block dev" :
                             (d_type == DT_CHR) ?  "char dev" : "???");
            printf("%4d %10lld  %s\n", d->d_reclen,
                    (long long) d->d_off, d->d_name);
            bpos += d->d_reclen;
        }
    }

    exit(EXIT_SUCCESS);
}
