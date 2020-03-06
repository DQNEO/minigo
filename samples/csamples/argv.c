#include <stdio.h>

static int global_argc;
static char **global_argv;

int main(int argc, char **argv) {
    if (argc > 1) {
        printf("argc=%d\n", argc);
        printf("argv[1]=%s\n", argv[1]);
    }
    global_argc = argc;
    global_argv = argv;
}
