// Fibonacci
n = 0;
i = 0;
j = 1;
while (n < 60) {
    k = i + j;
    i = j;
    j = k;
    n = n + 1;
    println(k);
}
