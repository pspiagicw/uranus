# `uranus`

`uranus` is a language built by following the book 'Writing a Interpreter in Go.'.
More information can be found inside book, but general introduction is provided here.

## Example Code

Statements end with `;`.

### Maths
```
(5 + 10 * 2 + 15 / 3) * 2 + -10
```

### Booleans

```
(1 > 2) == false
!!false
```

### If statements
```
if (1 < 2) { 10 } else { 20 }
```

### Variable Assignment
```
let a = 5; let b = a; let c = b + a + 5; c;
```

### Functions
```
let add = fn(x,y) { x + y;}; add(1,4)
```

Including closures.

```
let newAdder = fn(x) {
    fn(y) { x + y };
};

let addTwo = newAdder(2);
```

### Strings

```

```

## Language Features

Currently the language supports

- Variables
- Booleans
- Integers
- Functions (Including higher-order function and closures.)
- Strings
- Assignment statements
- Return Statements

## Interpreter Features

It has a custom built lexer. Along with a Pratt Parser. The evaluator is a tree-walking evaluator.

## Code Structure

The code is structured with 6 major parts.

- AST
- Lexer
- Token
- Parser
- Object
- Evaluator

### Token

This package defines what kind of symbols are used in `uranus`. It also defines the type of INTEGERS and BOOLEANS.

### Lexer

This package separates the input into a stream of tokens according to the rules defined in the `token` package.


### AST
These are the rules of the language, this package defines the structure of each expression found in `uranus`.

### Parser

This package parses the tokens in into a IR, according to the rules defined by the `ast` package.

### Object

The language's internal object system. This does not mean the language is object-oriented. It simply means internally a structure for encapsulation each value is used.

### Evaluator

This package contains the all-important `eval` function. Which takes any Node and parses returning a value.

## Testing

All the tests are written alongside the code itself.
To test the code use

```sh
make test
```
