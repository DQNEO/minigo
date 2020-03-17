#include <time.h>
#include <stdio.h>
#include <unistd.h>

/**
struct timespec {
   time_t tv_sec;        // seconds: 64bit
   long   tv_nsec;       // nanoseconds: 64bit
};
*/
int main() {
    struct timespec spec;
    spec.tv_sec  = 2;
    spec.tv_nsec = 0;
    printf("size of time_t = %ld\n", sizeof(spec.tv_sec));
    printf("size of long = %ld\n", sizeof(spec.tv_nsec));
    write(1, "hello ", 6);
    nanosleep(&spec, NULL);
    write(1, "woke up\n ", 8);
    return 0;
}
