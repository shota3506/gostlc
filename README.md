# gostlc - Simply Typed Lambda Calculus Interpreter

Go implemetation of Simple Typed Lambda Calculus (STLC) with base types - with zero external dependencies.

This is a toy interpreter for self learning purposes.

## Features

### Core Functionality
- STLC support: Lambda abstractions, applications, and variables
- Type system: Static type checking with Int and Bool base types, plus function types
- Conditional expressions: if-then-else constructs with type checking
- Literals: Integer and boolean literal support

### Language Features

#### Supported Types
- `Int` - Integer type
- `Bool` - Boolean type
- `T1 -> T2` - Function type from T1 to T2

#### Supported Syntax

The EBNF grammar for the supported subset of STLC is as follows:

```
expr ::= var
       | "\" var ":" type "." expr         (* abstraction *)
       | expr expr                         (* application *)
       | "(" expr ")"                      (* grouping *)
       | "true" | "false"                  (* boolean literals *)
       | "if" expr "then" expr "else" expr (* conditional *)
       | digit+                            (* integer literals *)

type ::= "Bool"                            (* boolean type *)
       | "Int"                             (* integer type *)
       | type "->" type                    (* function type *)
       | "(" type ")"                      (* grouping *)

var  ::= letter (letter | digit)*          (* variable names *)
```

## Installation

```bash
go install github.com/shota3506/gostlc/cmd/gostlc@latest
```

## Usage

### Interactive REPL

Start the REPL by running `gostlc` without arguments:

```bash
$ gostlc
STLC REPL
Type :quit or :q to exit, :help for help

gostlc> (\x:Int. x) 42
=> 42
gostlc> if true then 1 else 0
=> 1
```

### Execute from File

```bash
gostlc sample.stlc
```

### Execute from Command Line

```bash
gostlc -c "(\x:Int. x) 42"
```

### Execute from stdin

```bash
echo "(\x:Bool. x) true" | gostlc -
```

## Examples

### Identity Function
```stlc
(\x:Int. x) 42
# Result: 42
```

### Function Composition
```stlc
(\f:Int->Int. \g:Int->Int. \x:Int. f (g x))
```

### Church Numerals
```stlc
# Church numeral for 2
\f:Int->Int. \x:Int. f (f x)

# Church numeral for 3
\f:Int->Int. \x:Int. f (f (f x))
```

### Conditional Logic
```stlc
if (\x:Bool. x) true then 100 else 200
# Result: 100
```

### Higher-Order Functions
```stlc
(\f:Int->Int. f 10) (\x:Int. x)
# Result: 10
```

## TODOs

- Let bindings: `let x = e1 in e2` for local variable binding
- Recursive functions: Fixed-point operator or recursive let bindings
- Product types: Pairs/tuples with projection operations
- Sum types: Either/variant types with pattern matching
- Unit type: `()` for side-effect operations
- Type inference: Hindley-Milner style type inference to reduce type annotations
- Arithmetic operators: add, sub, mul, div, mod
- Comparison operators: eq, ne, lt, le, gt, ge
- Boolean operators: and, or, not
- String type and operations: String literals and concatenation
- Debugger: AST inspection and step-by-step evaluation


## License

This project is licensed under the MIT License.
