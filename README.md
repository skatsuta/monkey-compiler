# 🐵 Monkey compiler 🐒

[![wercker status](https://app.wercker.com/status/d99bbfc328db0183527cfb7e259fbdbf/s/master "wercker status")](https://app.wercker.com/project/byKey/d99bbfc328db0183527cfb7e259fbdbf)
[![Go Report Card](https://goreportcard.com/badge/github.com/skatsuta/monkey-compiler)](https://goreportcard.com/report/github.com/skatsuta/monkey-compiler)


Monkey programming language compiler designed in [_Writing A Compiler In Go_ (Ball, T. 2018)](https://compilerbook.com). The book is  awesome and I believe every programmer who mainly uses dynamically typed languages such as Ruby or Python should read it.

This implementation has added several features to the one implemented in the above book:

* Added support for running a single Monkey script file
* Added support for single-line comments (`#`)
* Added support for floating-point numbers and their arithmetic (`+`, `-`, `*`, `/`) and comparison (`<`, `>`, `==`, `!=`) operations 
* Added support for "greater than or equal to" (`>=`) and "less than or equal to" (`<=`) comparison operators
* Added support for logical AND (`&&`) and OR (`||`) operators
* Added support for variable assignment statements without `let` keyword
* Added support for variable reassignment statements
* Added support for setting values into existing arrays and hash maps
* Added support for `nil` literal


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

### Number types and variable bindings

You can define and reassign to variables using `=` operator. Variables are dynamically typed and can be assigned to objects of any type in Monkey. You can use `let` keyword when defining variables, but it's completely optional and there is no difference between with and without `let` keyword. 

Two number types are supported in this implementation: integers and floating-point numbers.

```sh
>> let a = 1;  # Assignment with `let` keyword
>> a
1
>> b = 2.5;  # Assignment without `let` keyword
>> b
2.5
>> b = "a";  # Reassignment to b
>> b
a
```

### Arithmetic and comparison expressions

You can do basic arithmetic and comparison operations for numbers, such as `+`, `-`, `*`, `/`, `<`, `>`, `<=`, `>=`, `==`, `!=`, `&&` and `||`.

```sh
>> let a = 10;
>> let b = a * 2;
>> (a + b) / 2 - 3;
12
>> let c = 2.25;
>> let d = -5.5;
>> b + c * d
7.625
>> a < b
true
>> c == d
false
```

### If expressions

You can use `if` and `else` keywords for conditional expressions. The last value in an executed block is returned from the expression.

```sh
>> let a = 10;
>> let b = a * 2;
>> let c = if (b > a) { 99 } else { 100 };
>> c
99
>> let d = if (b > a && c < b) { 199 } else { 200 };
>> d
200
```

### Functions and closures

You can define functions using `fn` keyword. All functions are closures in Monkey and you have to use `let` along with `fn` to bind a closure to a variable. Closures close over an environment where they are defined, and are evaluated in *the* environment when called. The last value in an executed function body is returned as a return value.

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

You can build arrays using square brackets `[]`. Array literal is `[value1, value2, ...]`. Arrays can contain values of any type, such as integers, strings, even arrays and functions (closures). To get an element at an index from an array, use `array[index]` syntax. To set a value at an index in an array to another value, use `array[index] = value` syntax.

```sh
>> let myArray = ["Thorsten", "Ball", 28, fn(x) { x * x }];
>> myArray[0]
Thorsten
>> myArray[4 - 2]
28
>> myArray[3](2);
4
>> myArray[2] = myArray[2] + 1
>> myArray[2]
29
```

### Hash maps

You can build hash maps using curly brackets `{}`. Hash literal is `{key1: value1, key2: value2, ...}`. You can use numbers, strings and booleans as keys, and objects of any type as values. To get a value under a key from a hash map, use `hash[key]` syntax. To set a value under a key in a hash map to another value, use `hash[key] = value` syntax.

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
>> myHash[0] = "right, zero"
>> myHash[0]
right, zero
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

### Comments

You can write single-line comments by starting with `#`. Comments begin with a hash mark (`#`) and continue to the end of the line. Thery are ignored by the compiler.

```sh
>> # This line is just a comment.
>> let a = 1;  # This is an integer.
1
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

##### fibonacci.monkey

```
let fib = fn(x) {
   if (x <= 1) {
     return x;
   }
   fib(x - 1) + fib(x - 2);
};

let N = 15;
puts(fib(N));
```

Running the above script gives us:


```sh
$ $GOPATH/bin/monkey-compiler fibonacci.monkey
610
```

Other example Monkey scripts are also placed in `examples` directory.
