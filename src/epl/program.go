// 
// Copyright (c) 2015 Brian William Wolter, All rights reserved.
// EPL - A little Embeddable Predicate Language
// 
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
// 
//   * Redistributions of source code must retain the above copyright notice, this
//     list of conditions and the following disclaimer.
// 
//   * Redistributions in binary form must reproduce the above copyright notice,
//     this list of conditions and the following disclaimer in the documentation
//     and/or other materials provided with the distribution.
//     
//   * Neither the names of Brian William Wolter, Wolter Group New York, nor the
//     names of its contributors may be used to endorse or promote products derived
//     from this software without specific prior written permission.
//     
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT,
// INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
// LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE
// OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED
// OF THE POSSIBILITY OF SUCH DAMAGE.
// 

package epl

import (
  "os"
  "io"
  "fmt"
  "reflect"
)

/**
 * A runtime context
 */
type Context interface {
  Variable(name string)(interface{}, error)
}

/**
 * A variable provider
 */
type VariableProvider func(name string)(interface{}, error)

/**
 * Runtime context
 */
type runtime struct {
  stdout    io.Writer
}

/**
 * Executable context
 */
type context struct {
  stack   []interface{}
}

/**
 * Create a new context
 */
func newContext(f interface{}) *context {
  return &context{[]interface{}{f}}
}

/**
 * Push a frame
 */
func (c *context) push(f interface{}) *context {
  if l := len(c.stack); cap(c.stack) > l {
    c.stack = c.stack[:l+1]
    c.stack[l] = f
  }else{
    c.stack = append(c.stack, f)
  }
  return c
}

/**
 * Pop a frame
 */
func (c *context) pop(f interface{}) *context {
  if l := len(c.stack); l > 0 {
    c.stack = c.stack[:l-1]
  }
  return c
}

/**
 * Obtain a value
 */
func (c *context) get(n string) (interface{}, error) {
  if l := len(c.stack); l > 0 {
    return derefProp(c.stack[l-1], n)
  }else{
    return nil, fmt.Errorf("No context")
  }
}

/**
 * Executable
 */
type executable interface {
  exec(*runtime, *context)([]interface{}, error)
}

/**
 * Executable
 */
type tree interface {
  add(c executable) *node
}

/**
 * An AST node
 */
type node struct {
  subnodes  []executable
}

/**
 * Execute
 */
func (n *node) exec(runtime *runtime, context *context) ([]interface{}, error) {
  return nil, fmt.Errorf("No implementation")
}

/**
 * A program
 */
type Program struct {
  root executable
}

/**
 * Execute
 */
func (p *Program) Exec(context interface{}) ([]interface{}, error) {
  return p.root.exec(&runtime{os.Stdout}, newContext(context))
}

/**
 * A program
 */
type emptyNode struct {
  node
}

/**
 * An expression node
 */
type exprNode struct {
  node
}

/**
 * A logical OR node
 */
type logicalOrNode struct {
  node
  left, right executable
}

/**
 * Execute
 */
func (n *logicalOrNode) exec(runtime *runtime, context *context) ([]interface{}, error) {
  
  lvi, err := execReturnSingle(runtime, context, n.left)
  if err != nil {
    return nil, err
  }
  lv, err := asBool(lvi)
  if err != nil {
    return nil, err
  }
  
  if lv {
    return []interface{}{true}, nil
  }
  
  rvi, err := execReturnSingle(runtime, context, n.right)
  if err != nil {
    return nil, err
  }
  rv, err := asBool(rvi)
  if err != nil {
    return nil, err
  }
  
  return []interface{}{rv}, nil
}

/**
 * A logical AND node
 */
type logicalAndNode struct {
  node
  left, right executable
}

/**
 * Execute
 */
func (n *logicalAndNode) exec(runtime *runtime, context *context) ([]interface{}, error) {
  
  lvi, err := execReturnSingle(runtime, context, n.left)
  if err != nil {
    return nil, err
  }
  lv, err := asBool(lvi)
  if err != nil {
    return nil, err
  }
  
  if !lv {
    return []interface{}{false}, nil
  }
  
  rvi, err := execReturnSingle(runtime, context, n.right)
  if err != nil {
    return nil, err
  }
  rv, err := asBool(rvi)
  if err != nil {
    return nil, err
  }
  
  return []interface{}{rv}, nil
}

/**
 * An arithmetic expression node
 */
type arithmeticNode struct {
  node
  op          token
  left, right executable
}

/**
 * Execute
 */
func (n *arithmeticNode) exec(runtime *runtime, context *context) ([]interface{}, error) {
  
  lvi, err := execReturnSingle(runtime, context, n.left)
  if err != nil {
    return nil, err
  }
  lv, err := asNumber(lvi)
  if err != nil {
    return nil, err
  }
  
  rvi, err := execReturnSingle(runtime, context, n.right)
  if err != nil {
    return nil, err
  }
  rv, err := asNumber(rvi)
  if err != nil {
    return nil, err
  }
  
  switch n.op.which {
    case tokenAdd:
      return []interface{}{ lv + rv }, nil
    case tokenSub:
      return []interface{}{ lv - rv }, nil
    case tokenMul:
      return []interface{}{ lv * rv }, nil
    case tokenDiv:
      return []interface{}{ lv / rv }, nil
    default:
      return nil, fmt.Errorf("Invalid operator: %v", n.op)
  }
  
}

/**
 * An relational expression node
 */
type relationalNode struct {
  node
  op          token
  left, right executable
}

/**
 * Execute
 */
func (n *relationalNode) exec(runtime *runtime, context *context) ([]interface{}, error) {
  
  lvi, err := execReturnSingle(runtime, context, n.left)
  if err != nil {
    return nil, err
  }
  rvi, err := execReturnSingle(runtime, context, n.right)
  if err != nil {
    return nil, err
  }
  
  switch n.op.which {
    case tokenEqual:
      return []interface{}{ lvi == rvi }, nil
    case tokenNotEqual:
      return []interface{}{ lvi != rvi }, nil
  }
  
  lv, err := asNumber(lvi)
  if err != nil {
    return nil, err
  }
  rv, err := asNumber(rvi)
  if err != nil {
    return nil, err
  }
  
  switch n.op.which {
    case tokenLess:
      return []interface{}{ lv < rv }, nil
    case tokenGreater:
      return []interface{}{ lv > rv }, nil
    case tokenLessEqual:
      return []interface{}{ lv <= rv }, nil
    case tokenGreaterEqual:
      return []interface{}{ lv >= rv }, nil
    default:
      return nil, fmt.Errorf("Invalid operator: %v", n.op)
  }
  
}

/**
 * A dereference expression node
 */
type derefNode struct {
  node
  left  executable
  ident string
}

/**
 * Execute
 */
func (n *derefNode) exec(runtime *runtime, context *context) ([]interface{}, error) {
  v, err := execReturnSingle(runtime, context, n.left)
  if err != nil {
    return nil, err
  }
  return derefProp(v, n.ident)
}

/**
 * An identifier expression node
 */
type identNode struct {
  node
  ident string
}

/**
 * Execute
 */
func (n *identNode) exec(runtime *runtime, context *context) ([]interface{}, error) {
  v, err := context.get(n.ident)
  if err != nil {
    return nil, err
  }else{
   return []interface{}{v}, nil
 }
}

/**
 * A literal expression node
 */
type literalNode struct {
  node
  value interface{}
}

/**
 * Execute
 */
func (n *literalNode) exec(runtime *runtime, context *context) ([]interface{}, error) {
  return []interface{}{n.value}, nil
}

/**
 * Execute expecting a single return value
 */
func execReturnSingle(runtime *runtime, context *context, e executable) (interface{}, error) {
  rva, err := e.exec(runtime, context)
  if err != nil {
    return nil, err
  }
  if len(rva) != 1 {
    return nil, fmt.Errorf("Expected single return value, got %d", len(rva))
  }else{
    return rva[0], nil
  }
}

/**
 * Obtain an interface value as a bool
 */
func asBool(value interface{}) (bool, error) {
  switch v := value.(type) {
    case bool:
      return v, nil
    case uint8:
      return v != 0, nil
    case uint16:
      return v != 0, nil
    case uint32:
      return v != 0, nil
    case uint64:
      return v != 0, nil
    case int8:
      return v != 0, nil
    case int16:
      return v != 0, nil
    case int32:
      return v != 0, nil
    case int64:
      return v != 0, nil
    default:
      return false, fmt.Errorf("Cannot cast %v (%T) to bool", value, value)
  }
}

/**
 * Obtain an interface value as a number
 */
func asNumber(value interface{}) (float64, error) {
  switch v := value.(type) {
    case int:
      return float64(v), nil
    case uint8:
      return float64(v), nil
    case uint16:
      return float64(v), nil
    case uint32:
      return float64(v), nil
    case uint64:
      return float64(v), nil
    case int8:
      return float64(v), nil
    case int16:
      return float64(v), nil
    case int32:
      return float64(v), nil
    case int64:
      return float64(v), nil
    case float32:
      return float64(v), nil
    case float64:
      return v, nil
    default:
      return 0, fmt.Errorf("Cannot cast %v (%T) to numeric", value, value)
  }
}

/**
 * Dereference
 */
func derefProp(context interface{}, ident string) ([]interface{}, error) {
  switch v := context.(type) {
    
    case Context:
      r, err := v.Variable(ident)
      if err != nil {
        return nil, err
      }else{
        return []interface{}{r}, nil
      }
      
    case VariableProvider:
      r, err := v(ident)
      if err != nil {
        return nil, err
      }else{
        return []interface{}{r}, nil
      }
      
    case map[string]interface{}:
      return []interface{}{v[ident]}, nil
      
    default:
      return derefField(context, ident)
      
  }
}

/**
 * Execute
 */
func derefField(context interface{}, property string) ([]interface{}, error) {
  c, _ := derefValue(reflect.ValueOf(context))
  switch c.Kind() {
    case reflect.Struct:
      return []interface{}{c.FieldByName(property)}, nil
    default:
      return nil, fmt.Errorf("Cannot dereference context: %v (%T)", context, context)
  }
}

/**
 * Dereference a value
 */
func derefValue(value reflect.Value) (reflect.Value, int) {
  v := value
  c := 0
  for ; v.Kind() == reflect.Ptr; {
    v = v.Elem()
    c++
  }
  return v, c
}
