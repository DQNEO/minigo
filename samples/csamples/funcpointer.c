#include <stdio.h>

typedef int (*MyFunc)(int a,int b);


int sum(int a,int b)
{
    return a + b;
}

int multiple(int a,int b)
{
    return a * b;
}

void main()
{
    int a,b;
    int result;
    MyFunc fnc;

    a = 3;
    b = 4;

    fnc = &sum;
    result = (fnc)(a,b);
    printf("sum = %d\n",result);

    fnc = &multiple;
    result = (fnc)(a,b);
    printf("multiple = %d\n",result);
}
