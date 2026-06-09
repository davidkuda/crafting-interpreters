# crafting-interpreters

This repository holds the code that I write while following the book [Crafting Interpreters](https://craftinginterpreters.com/) by Robert Nystrom.

In the book, you will create two interpreters.

The first interpreter focuses on concepts. The code samples are in Java. **For the first interpreter, I chose Go as my programming language.**

From the book:

> We’ll write our first interpreter, jlox, in Java. The focus is on concepts. We’ll write the simplest, cleanest code we can to correctly implement the semantics of the language. This will get us comfortable with the basic techniques and also hone our understanding of exactly how the language is supposed to behave.


The second interpreter is in C. From the book:

> So in the next part, we start all over again, but this time in C. C is the perfect language for understanding how an implementation really works, all the way down to the bytes in memory and the code flowing through the CPU.

# Lox expression grammar

```
expression -> literal
            | unary
            | binary
            | grouping
            ;

literal    -> NUMBER
            | STRING
            | "true"
            | "false"
            | "nil"
            ;

grouping   -> "(" expression ")" ;

unary      -> ( "-" | "!" ) expression ;

binary     -> expression operator expression ;

operator   -> "=="
            | "!="
            | "<"
            | "<="
            | ">"
            | ">="
            | "+"
            | "-"
            | "*"
            | "/"
            ;
```

The problem with this grammar is its ambiguity. We need precedence and associativity rules to make the grammar unambiguous.


### Precedence and Associativity

from lowest to highest precedence (same as in C):

|Name        | Operators         | Associates |
| ---------- | ----------------- | ---------- |
| Equality   | `==` `!=`         | Left       |
| Comparison | `>` `>=` `<` `<=` | Left       |
| Term       | `-` `+`           | Left       |
| Factor     | `/` `*`           | Left       |
| Unary      | `!` `-`           | Right      |


This yields a new grammar (without binary and without operators):

```
expression -> equality ;

equality   -> comparison
              ( ( "!=" | "==" ) comparison )* ;

comparison -> term
              ( ( ">" | ">=" | "<" | "<=" ) term )* ;

term       -> factor
              ( ( "-" | "+" ) factor )* ;

factor     -> unary
              ( ( "/" | "*" ) unary )* ;

unary      -> ( "-" | "!" ) unary
            | primary ;

primary    -> NUMBER
            | STRING
            | "true"
            | "false"
            | "nil"
            | "(" expression ")" ;
```

The `( )*` means 0 or many times.

On left-recursiveness:

The following two rules are equivalent:

```
factor     -> factor ( "/" | "*" ) unary
              | unary ;

factor     -> unary
              ( ( "/" | "*" ) unary )* ;
```

Notice how it's similar to unary. (unary is in fact right-associative.)

Problem is its left-recursiveness. This will require a different technique to parse than what we are going to use for golox. therefore, we use the second expression rule.


## Evolution in chapter 8 Statements and State

The middle of the grammar tree stays the same.

At the top, we prepend:

```
program     -> decleration* EOF ;
declaration -> varDecl | statement ;
varDecl     -> "var" IDENTIFIER ( "=" expression )? ";" ;
statement   -> exprStmt | printStmt ;
exprStmt    -> expression ";" ;
printStmt   -> "print" expression ";" ;
```

`declaration` for now only contains variable declarations or statements. Later, we'll add functions and classes.

At the bottom, we add one line to `primary`:

```
primary    -> NUMBER
            | STRING
            | "true"
            | "false"
            | "nil"
            | "(" expression ")
            | IDENTIFIER
            ;
```

### 8.4.1 Assignment

```
expression -> assignment ;

assignment -> IDENTIFIER "=" assignment
              | equality
              ;

```

### 8.5.2 Blocks

```
statement    -> exprStmt
              | printStmt
              | block
              ;

block        -> "{" declaration* "}" ;
```

### New grammar:

```
program     -> decleration* EOF ;
declaration -> varDecl | statement ;
varDecl     -> "var" IDENTIFIER ( "=" expression )? ";" ;
statement   -> exprStmt | printStmt | block ;
exprStmt    -> expression ";" ;
printStmt   -> "print" expression ";" ;

expression -> assignment ;
assignment -> IDENTIFIER "=" assignment | equality ;
equality    -> comparison ( ( "!=" | "==" ) comparison )* ;
comparison  -> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term        -> factor ( ( "-" | "+" ) factor )* ;
factor      -> unary ( ( "/" | "*" ) unary )* ;
unary       -> ( "-" | "!" ) unary | primary ;

primary     -> NUMBER
            | STRING
            | "true"
            | "false"
            | "nil"
            | "(" expression ")
            | IDENTIFIER
            ;
```


## Evolution in chapter 9 Control Flow

9.2 Conditional Execution
```
statement   -> exprStmt | ifStmt | printStmt | block ;
ifStmt      -> "if" "(" expression ")" statement
               ( "else" statement )? ;
```

9.3 Logical Operators

```
expression -> assignment ;
assignment -> IDENTIFIER "=" assignment | logic_or ;
logic_or   -> logic_and ( "or" logic_and )* ;
logic_and  -> equality ( "and" equality )* ;
```

