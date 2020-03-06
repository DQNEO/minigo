#include <stdio.h>

typedef struct MyStruct {
    int a;
    int b;
    int c;
    int d;
    int e;
    int f;
    int g;
} MyStruct;

MyStruct f() {
    MyStruct s;
    s.a = 1;
    s.b = 2;
    s.c = 3;
    s.d = 4;
    s.e = 4;
    s.f = 4;
    s.g = 4;
    return s;
}

int main() {
    MyStruct s2;
    s2 = f();
    printf("%d\n", s2.c);
}
