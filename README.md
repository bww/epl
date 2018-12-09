# Embedded Predicate Language

EPL is a tiny, general-purpose, expression-only language compiler and runtime implemented in Go. A Go program can use EPL to express arbitrarily complex predicate logic tersely and unambiguously or to accept such logic from a user.

```go
context := map[string]interface{}{
  "greeting": "Hello",
}

program, _ := epl.Compile(`greeting == "Hello"`)
result, _ := program.Exec(context)
fmt.Println(result) // true
```

## History

This Go version is the latest and most fully realized incarnation of this project. A previous version was [Predicate Kit](https://github.com/bww/PredicateKit), implemented in Objective-C and intended as a more flexible replacement for [`NSPredicate`](https://developer.apple.com/documentation/foundation/nspredicate?changes=_5).

# Syntax
EPL will be reasonably familiar to anyone with experience using langauges with C-style syntax.

## Identifiers
Identifiers in EPL have the same rules as Go. A valid identifier is a letter followed by zero or more letters or digits. An underscore is considered to be a letter.
```
a
ThisIsALongIdentifier
_a9
```

## UUID Identifiers
In addition to normal identifiers EPL has support for the use of UUIDs as identifiers. A UUID identifier is the sequence `u:` (or `U:`, if you prefer) followed by a valid UUID. A UUID identifier can be used anywhere a normal identifier can be used.
```
u:9515976f-cdb4-4e56-bd07-b1ae6efc00da
U:7388AA2B-44C3-4146-8F17-C78F89B5F7D8
```

## Literals
String, number, and boolean literals are supported.

## Strings
Strings literals have essentially the same rules as Go. A string begins with `"`, is terminated by `"`, and contains zero or more characters or escape sequences. Unlike regular Go strings, an EPL string may contain newlines.
```
"Hello!"
""
"Hello
       world!"
```

The following escape sequences are allowed in strings. All escape sequences are introduced by a backslash `\` character.

| Escape | Value |
|--------|-------|
| `\uXXXX` | A Unicode character with the codepoint `XXXX`. |
| `\xXX` | A Unicode character with the codepoint `XX`. |
| `\\` | A literal `\` |
| `\"` | A literal `"` |
| `\a` | Audible bell | 
| `\b` | Backspace |
| `\f` | Form feed |
| `\n` | Newline |
| `\r` | Carriage return |
| `\t` | Tab |
| `\v` | Vertical tab |

## Numbers
Numeric literals have essentially the same rules as Go. Decimal, octal, and hexidecimal integers and decimal floating point numeric literals are supported. Floating points may use exponential notation.

### Integers
```
42
0600
0xBadFace
170141183460469231731687303715884105727
```

### Floating points
```
0.
72.40
2.71828
1.e+0
6.67428e-11
1E6
```

## Booleans
Boolean literals are `true` and `false`.
```
true
false
```

## Nil
The special literal Nil uses the special identifier `nil`.
```
nil
```

## `&&`, `||`, `!` Logical Operators
The standard logical and, or, and not operators are supported and have the same meaning as in Go.
```
1 < 2 && 2 < 3
1 > 2 || 2 < 3
!(1 > 2)
```

## `<`, `<=`, `==`, `>=`, `>` Relational Operators
The standard relational operators are supported. Values are comparable if their types are comparable in Go. Unlike Go, EPL will automatically convert numeric types so that they can be compared.
```
 1 < 2
 1 <= 3
 1 == 1
 5 >= 5
 5 > 4
```

## `.` Dereference Operator
The `.` operator dereferences a property. This operator can be use more liberally in Ego than Go. You can use this operator to:

* Obtain the value of an exported struct field
* Obtain a value from a map that has string keys if the key is also a valid identifier. That is: `a_map.string_key` is equivalent to `a_map["string_key"]`.
* Obtain the result of a method invocation if that method does not take any arguments. That is: `an_interface.Foo` is equivalent to `an_interface.Foo()`.

## `[]` Subscript Operator
The subscript operator obtains the value at an index when the operand is an array or slice and obtains the value of a key when the operand is a map.
```
 a_slice[5]
 a_map["the_key"]
```

## `()` Functions and Methods
Functions and methods are invoked as they are in Go. When a function invocation follows a dereference it is treated as a method invocation.
```
len(v)
val.Len()
```

The manner in which underlying Go functions are mapped to EPL is a bit more nuanced, however. There are various rules governing how different return values are handled, aimed at producing the expected result.

# Standard Library
A few builtin functions are provided in the standard library. They are aimed at providing functionality that cannot be reasonably addressed using the language syntax alone.

| Function | Detail |
|----------|--------|
| `len(v)` | Determine the length of an array, slice, or map, `v`. Providing any other type as an argument will produce an error. ` |
| `match(e, v)` | Match the regular expression `e` in the text `v`. If a match is found, `true` is returned, otherwise `false`. |
| `printf(...v)` | Print to standard output. This method has the same semantics as `fmt.Printf`. |

## Environment
The current environment is exposed via the variable `env`. Environment variables are accessed by their name as properties of `env`.
```
env.TERM
env.POSTGRES_HOME
```



