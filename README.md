# 🐵 Monkey compiler 🐒

[![wercker status](https://app.wercker.com/status/d99bbfc328db0183527cfb7e259fbdbf/s/master "wercker status")](https://app.wercker.com/project/byKey/d99bbfc328db0183527cfb7e259fbdbf)


Monkey programming language compiler designed in [_Writing A Compiler In Go_ (Ball, T. 2018)](https://compilerbook.com). The book is  awesome and I believe every programmer who mainly uses dynamically typed languages such as Ruby or Python should read it.


## Prerequisites

Go 1.12 or later is required to build the compiler.


## Usage

Install the Monkey compiler using `go get` command:

```sh
$ go get -v -u github.com/skatsuta/monkey-compiler/...
```

The easiest way to get started is to run REPL:

```
$ $GOPATH/bin/monkey-compiler
This is the Monkey programming language!
Feel free to type in commands
>> puts("Hello, world!")
Hello, world!
>> 
```

The compiler also supports running a single Monkey script file (for example `script.monkey` file):

```sh
$ cat script.monkey
puts("Hello, world!")
$ $GOPATH/bin/monkey-compiler script.monkey
Hello, world!
```

## Getting started with Monkey

### Integer types and variable bindings

You can define variables using `let` keyword. Supported number types are integers.

```sh
>> let a = 1;
>> a
1
```

### Arithmetic expressions

You can do basic arithmetic operations for numbers, such as `+`, `-`, `*` and `/`. 

```sh
>> let a = 10;
>> let b = a * 2;
>> (a + b) / 2 - 3;
12
>> let c = 2;
>> b + c
22
```

### If expressions

You can use `if` and `else` keywords for conditional expressions. The last value in an executed block are returned from the expression.

```sh
>> let a = 10;
>> let b = a * 2;
>> let c = if (b > a) { 99 } else { 100 };
>> c
99
```

### Functions and closures

You can define functions using `fn` keyword. All functions are closures in Monkey and you have to use `let` along with `fn` to bind a closure to a variable. Closures close over an environment where they are defined, and are evaluated in *the* environment when called. The last value in an executed function body are returned as a return value.

```sh
>> let multiply = fn(x, y) { x * y };
>> multiply(50 / 2, 1 * 2)
50
>> fn(x) { x + 10 }(10)
20
>> let newAdder = fn(x) { fn(y) { x + y }; };
>> let addTwo = newAdder(2);
>> addTwo(3);
5
>> let sub = fn(a, b) { a - b };
>> let applyFunc = fn(a, b, func) { func(a, b) };
>> applyFunc(10, 2, sub);
8
```

### Strings

You can build strings using a pair of double quotes `""`. Strings are immutable values just like numbers. You can concatenate strings with `+` operator.

```sh
>> let makeGreeter = fn(greeting) { fn(name) { greeting + " " + name + "!" } };
>> let hello = makeGreeter("Hello");
>> hello("John");
Hello John!
```

### Arrays

You can build arrays using square brackets `[]`. Arrays can contain any type of values, such as integers, strings, even arrays and functions (closures). To get an element at an index from an array, use `array[index]` syntax.

```sh
>> let myArray = ["Thorsten", "Ball", 28, fn(x) { x * x }];
>> myArray[0]
Thorsten
>> myArray[4 - 2]
28
>> myArray[3](2);
4
```

### Hash tables

You can build hash tables using curly brackets `{}`. Hash literals are `{key1: value1, key2: value2, ...}`. You can use numbers, strings and booleans as keys, and any type of objects as values. To get a value of a key from a hash table, use `hash[key]` syntax.

```sh
>> let myHash = {"name": "Jimmy", "age": 72, true: "yes, a boolean", 99: "correct, an integer"};
>> myHash["name"]
Jimmy
>> myHash["age"]
72
>> myHash[true]
yes, a boolean
>> myHash[99]
correct, an integer
```

### Built-in functions

There are some built-in functions in Monkey.

#### `len`

`len` built-in function allows you to get the length of strings or arrays. Note that `len` returns the number of bytes instead of characters for strings.

```sh
>> len("hello");
5
>> len("∑");
3
>> let myArray = ["one", "two", "three"];
>> len(myArray)
3
```

#### `puts`

`puts` built-in function allows you to print out one or more objects to console (i.e. stdout).

```sh
>> puts("Hello, World")
Hello, World
```

#### `first`

`first` built-in function allows you to get the first element from an array. If the array is empty, `first` returns `nil`.

```sh
>> let myArray = ["one", "two", "three"];
>> first(myArray)
one
>> first([])
nil
```

#### `last`

`last` built-in function allows you to get the last element from an array. If the array is empty, `last` returns `nil`.

```sh
>> let myArray = ["one", "two", "three"];
>> last(myArray)
three
>> last([])
nil
```

#### `rest`

`rest` built-in function allows you to create a new array containing all elements of a given array except the first one. If the array is empty, `rest` returns `nil`.

```sh
>> let myArray = ["one", "two", "three"];
>> rest(myArray)
[two, three]
>> rest([])
nil
```

#### `push`

`push` built-in function allows you to add a new element to the end of an existing array. It allocates a new array instead of modifying the given one.

```sh
>> let myArray = ["one", "two", "three"];
>> push(myArray, "four")
[one, two, three, four]
```

#### `quote` / `unquote`

Special function, `quote`, returns an unevaluated code block (think it as an AST). Opposite function to `quote`, `unquote`, evaluates code inside `quote`.

```sh
>> quote(2 + 2)
Quote((2 + 2)) # Unevaluated code
>> quote(unquote(1 + 2))
Quote(3)
```

### Macros

You can define macros using `macro` keyword. Note that macro definitions must return `Quote` objects generated from `quote` function.

```sh
# Define `unless` macro which does the opposite to `if`
>> let unless = macro(condition, consequence, alternative) {
     quote(
       if (!(unquote(condition))) {
         unquote(consequence);
       } else {
         unquote(alternative);
       }
     );
   };
>> unless(10 > 5, puts("not greater"), puts("greater"));
greater
nil
```

### Example

Here is a Fibonacci function implemented in Monkey:

```
// fib.monkey

let fib = fn(x) {
   if (x == 0) {
     0;
   } else {
     if (x == 1) {
       1;
     } else {
       fib(x - 1) + fib(x - 2);
     }
   }
};
puts(fib(15));
```

Running the above script gives us:


```sh
$ $GOPATH/bin/monkey-compiler fib.monkey
610
```
