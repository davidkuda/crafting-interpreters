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
