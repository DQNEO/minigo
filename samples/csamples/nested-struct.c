/**
 * C sample of a nested struct
  */

#include <stdio.h>

typedef struct {
    int id;
    int weight;
} Ring;

typedef struct {
    int id;
    Ring ring;
} Hobbit;

typedef struct {
    Hobbit hobbit;
} Hole;

Hole ghole;

void f_global() {
    ghole.hobbit.ring.weight = 7;
    printf("%d\n", ghole.hobbit.ring.weight);
}

void f_local() {
    Hole lhole;
    lhole.hobbit.ring.weight = 7;
    printf("%d\n", lhole.hobbit.ring.weight);
}

int main() {
    f_global();
    f_local();
}
