# bs

[![Build Status](https://secure.travis-ci.org/ricallinson/bs.png?branch=master)](http://travis-ci.org/ricallinson/bs)

__WARNING: IN PROGRESS__

Bash Scripting is a simple programming language that compiles to Bash.

## Use

### Execute Code

    bs /path/to/file.bs

### Generate Script

    bs -script /path/to/file.bs

## Syntax

The syntax of bs is [go-based](https://en.wikipedia.org/wiki/Go_(programming_language)) (derived from C programming language). If you have learned C, Java, C++ or JavaScript, bs is quite easy for you.

### Assignment

```javascript
a = 1
b = "string"
c = [1, 2, "str", true, false]
```

### Expression

```javascript
a = 1 + 2
b = a * 7
c = concat("Con", "cat")
d = c + b
```

### Command

```javascript
output = ls()
ex = exists("file.txt")
```

### If condition

```javascript
a = 3;
if (a > 2) {
    println("Yes")
} else if (a == 2) {
    println("No")
}
```

### Loop

```javascript
n = 0
i = 0
j = 1
while (n < 60) {
    k = i + j
    i = j
    j = k
    n = n + 1
    println(k)
}
```

### Function

```go
v1 = "Global V1"
v2 = "Global V2"
func foo(p) {
    v1 = concat("Local ", p)
    v2 = "V3 Modified."
}
foo("Var")
```

### Recursion

```javascript
func fibonacci(num) {
    if (num == 0) {
        return 0
    } else if (num == 1) {
        return 1
    } else {
        return fibonacci(num - 2) + fibonacci(num - 1)
    }
}
println(fibonacci(8))
```

## Built-in functions

### `print(text, ...)`

Prints a text string to console without a newline.

### `println(text, ...)`

Prints a text string to console with a new line.

### `call(path, arg, ...)`

Runs command from path through shell.

### TODO: `bash(rawStatement)`

Put `rawStatement` into compiled code for Bash.

### `ls(path)`

Equals to `ls` and `dir /w`.

### `exists(path)`

Test existence of given path. 

## Inspiration

This project was inspired by [Batsh](https://github.com/BYVoid/Batsh).
