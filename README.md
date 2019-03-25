# Monkey compiler

Monkey programming language compiler designed in [_Writing An Compiler In Go_](https://https://compilerbook.com/).


## Usage

Install the Monkey compiler using `go get`:

```sh
$ go get -v -u github.com/skatsuta/monkey-compiler/...
```

Then run REPL:

```sh
$ $GOPATH/bin/monkey-compiler
This is the Monkey programming language!
Feel free to type in commands
>> 
```

Or run a Monkey script file (for example `script.monkey` file):

```sh
$ $GOPATH/bin/monkey-compiler script.monkey
```

## Getting started with Monkey

### Variable bindings and arithmetic expressions

You can define variables using `let` keyword. Supported number types are integers and floating-point numbers. You can do usual arithmetic operations against numbers, such as `+`, `-`, `*` and `/`.

```sh
>> let a = 10;
>> let b = a * 2;
>> (a + b) / 2 - 3;
12
>> let c = 2.5;
>> b + c
22.5
```

### If expressions

You can use `if` and `else` keywords for conditional expressions. The final value in an executed block are returned from the expression.

```sh
>> let a = 10;
>> let b = a * 2;
>> let c = if (b > a) { 99 } else { 100 };
>> c
99
```

### Functions and closures

You can define functions using `fn` keyword. All functions are closures in Monkey, so you must use `let` along with `fn` to bind a closure to a variable. Closures enclose defined environments and are executed along with the environments (e.g. variables) when called. The final value in an executed function body are returned as a return value.

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

You can build strings using a pair of double quotes `""`. Strings are values just like numbers. You can concatenate strings with `+` operator.

```sh
>> let makeGreeter = fn(greeting) { fn(name) { greeting + " " + name + "!" } };
>> let hello = makeGreeter("Hello");
>> hello("John");
Hello John!
```

### Arrays

You can build arrays using square brackets `[]`. Arrays can contain any type of values, such as even arrays and functions (closures). To get an element at an index from an array, use `array[index]` syntax.

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

There are many built-in functions in Monkey, for example `len()`, `first()` and `last()`. Special function, `quote`, returns an unevaluated code block (think it as an AST). Opposite function to `quote`, `unquote`, evaluates code inside `quote`.

```sh
>> len("hello");
5
>> len("∑");
3
>> let myArray = ["one", "two", "three"];
>> len(myArray)
3
>> first(myArray)
one
>> rest(myArray)
[two, three]
>> last(myArray)
three
>> push(myArray, "four")
[one, two, three, four]
>> puts("Hello World")
Hello World
nil
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
