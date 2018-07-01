# Embedded Predicate Language

EPL is a tiny expression-only language compile and runtime implemented in Go. A Go program can use EPL to express arbitrarily complex predicate logic tersely and unambiguously or to accept such logic from a user.

```go
context := map[string]interface{}{
  "greeting": "Hello",
}
program, _ := epl.Compile(`greeting == "Hello"`)
result, _ := program.Exec(context)
fmt.Println(result) // true
```

The tests in `epl_test.go` do a reasonably good job of illustrating the usage of the language.

## History

This Go version is the latest and most fully realized incarnation of this project, a previous version being [Predicate Kit](https://github.com/bww/PredicateKit), implemented in Objective-C and intended as a more flexible replacement for [`NSPredicate`](https://developer.apple.com/documentation/foundation/nspredicate?changes=_5).
