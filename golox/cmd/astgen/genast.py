from dataclasses import dataclass


@dataclass
class ASTNodeTypeField:
    name: str
    go_type: str

@dataclass
class ASTNodeType:
    name: str
    fields: list[ASTNodeTypeField]


INPUT = [
    "Binary   : left Expr, operator Token, right Expr",
    "Grouping : expression Expr",
    "Literal  : value any",
    "Unary    : operator Token, right Expr",
]


def main():
    print("package main")
    print()
    print("type Expr interface{}")
    print()

    ts = create_ast_node_types()
    for i, t in enumerate(ts, 1):
        print(f"type {t.name} struct {{")

        longest_name = 0
        for field in t.fields:
            if len(field.name) > longest_name:
                longest_name = len(field.name)

        for field in t.fields:
            print(f"\t{field.name:<{longest_name}} {field.go_type}")

        print("}")

        if i < len(ts):
            print()


def create_ast_node_types():
    ast_types: list[ASTNodeType] = []

    for t in INPUT:
        name = t.split(":")[0].strip()
        fields = t.split(":")[1].strip().split(", ")

        ast_fields: list[ASTNodeTypeField] = []
        for field in fields:
            field_name, field_type = field.split(" ")
            ast_fields.append(
                ASTNodeTypeField(
                    name=field_name,
                    go_type=field_type,
                )
            )

        ast_types.append(
            ASTNodeType(
                name=name,
                fields=ast_fields,
            )
        )

    return ast_types


if __name__ == "__main__":
    main()
