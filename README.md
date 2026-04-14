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
